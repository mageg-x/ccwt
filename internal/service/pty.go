package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ccwt/ccwt/internal/config"
	"github.com/creack/pty"
)

const (
	maxScrollback = 5 * 1024 * 1024 // 5MB 回滚缓冲区
)

// PtySession 单个 PTY 会话
type PtySession struct {
	ID       string
	UserName string
	Cmd      *exec.Cmd
	Pty      *os.File
	Buf      *RingBuffer
	CreateAt time.Time
	mu       sync.Mutex
	// 订阅者：WebSocket 连接订阅输出
	subs   []chan []byte
	subsMu sync.Mutex
	done   chan struct{}
}

// Subscribe 订阅终端输出
func (s *PtySession) Subscribe() chan []byte {
	ch := make(chan []byte, 256)
	s.subsMu.Lock()
	s.subs = append(s.subs, ch)
	s.subsMu.Unlock()
	return ch
}

// Unsubscribe 取消订阅
func (s *PtySession) Unsubscribe(ch chan []byte) {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()
	for i, c := range s.subs {
		if c == ch {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)
			close(ch)
			return
		}
	}
}

// broadcast 将数据广播给所有订阅者
func (s *PtySession) broadcast(data []byte) {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()
	// 复制数据，防止并发问题
	cp := make([]byte, len(data))
	copy(cp, data)
	for _, ch := range s.subs {
		select {
		case ch <- cp:
		default:
			// 订阅者来不及消费，丢弃旧数据
		}
	}
}

// Done 返回会话结束的 channel
func (s *PtySession) Done() <-chan struct{} {
	return s.done
}

// RingBuffer 环形缓冲区
type RingBuffer struct {
	data []byte
	size int
	pos  int
	full bool
	mu   sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{data: make([]byte, size), size: size}
}

func (r *RingBuffer) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, b := range p {
		r.data[r.pos] = b
		r.pos = (r.pos + 1) % r.size
		if r.pos == 0 {
			r.full = true
		}
	}
	return len(p), nil
}

func (r *RingBuffer) Bytes() []byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.full {
		return append([]byte{}, r.data[:r.pos]...)
	}
	out := make([]byte, r.size)
	copy(out, r.data[r.pos:])
	copy(out[r.size-r.pos:], r.data[:r.pos])
	return out
}

// PtyManager 管理所有 PTY 会话
type PtyManager struct {
	sessions map[string]*PtySession
	mu       sync.RWMutex
}

var Pty = &PtyManager{
	sessions: make(map[string]*PtySession),
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func mergeEnv(base []string, overrides map[string]string) []string {
	idx := make(map[string]int, len(base))
	out := make([]string, 0, len(base)+len(overrides))
	for _, kv := range base {
		i := strings.IndexByte(kv, '=')
		if i <= 0 {
			continue
		}
		k := kv[:i]
		v := kv[i+1:]
		if j, ok := idx[k]; ok {
			out[j] = k + "=" + v
			continue
		}
		idx[k] = len(out)
		out = append(out, k+"="+v)
	}
	for k, v := range overrides {
		if j, ok := idx[k]; ok {
			out[j] = k + "=" + v
		} else {
			idx[k] = len(out)
			out = append(out, k+"="+v)
		}
	}
	return out
}

func writeBashInit(username string) (string, error) {
	userHome := config.UserDir(username)
	workspace := config.UserWorkspace(username)
	initPath := filepath.Join(userHome, ".ccwt_bashrc")

	content := fmt.Sprintf(`# CCWT managed bash init
export HOME=%s
export CLAUDE_CONFIG_DIR=%s
export HISTFILE=%s
export HISTSIZE=10000
export HISTFILESIZE=20000
shopt -s histappend cmdhist checkwinsize
export PROMPT_COMMAND="history -a; history -n${PROMPT_COMMAND:+; $PROMPT_COMMAND}"
__CCWT_WORKSPACE=%s
cd() {
  builtin cd "$@" || return
  case "$PWD/" in
    "$__CCWT_WORKSPACE"/|"$__CCWT_WORKSPACE"/*) ;;
    *)
      printf "\\033[33m[CCWT] 工作目录受限，已返回 workspace\\033[0m\\n"
      builtin cd "$__CCWT_WORKSPACE" || return
      ;;
  esac
}
case "$PWD/" in
  "$__CCWT_WORKSPACE"/|"$__CCWT_WORKSPACE"/*) ;;
  *) builtin cd "$__CCWT_WORKSPACE" ;;
esac
`, shellQuote(userHome), shellQuote(config.UserClaudeDir(username)), shellQuote(filepath.Join(userHome, ".bash_history")), shellQuote(workspace))

	if err := os.WriteFile(initPath, []byte(content), 0600); err != nil {
		return "", err
	}
	return initPath, nil
}

func ensureUserRuntimeDirs(username string) error {
	dirs := []string{
		config.UserDir(username),
		config.UserClaudeDir(username),
		config.UserWorkspace(username),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0700); err != nil {
			return err
		}
	}
	historyFile := filepath.Join(config.UserDir(username), ".bash_history")
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		if err := os.WriteFile(historyFile, []byte(""), 0600); err != nil {
			return err
		}
	}
	return nil
}

// Create 创建新的 PTY 会话
func (m *PtyManager) Create(id, username string, rows, cols uint16) (*PtySession, error) {
	if err := ensureUserRuntimeDirs(username); err != nil {
		log.Printf("用户目录初始化失败: user=%s err=%v", username, err)
		return nil, err
	}

	shell := os.Getenv("SHELL")
	if shell == "" || !strings.Contains(filepath.Base(shell), "bash") {
		shell = "/bin/bash"
	}

	bashInit, err := writeBashInit(username)
	if err != nil {
		log.Printf("写入 bash 初始化文件失败: user=%s err=%v", username, err)
		return nil, err
	}

	cmd := exec.Command(shell, "--noprofile", "--rcfile", bashInit, "-i")
	cmd.Dir = config.UserWorkspace(username)
	cmd.Env = mergeEnv(os.Environ(), map[string]string{
		"CLAUDE_CONFIG_DIR": config.UserClaudeDir(username),
		"HOME":              config.UserDir(username),
		"CCWT_USER":         username,
		"TERM":              "xterm-256color",
		"SHELL":             shell,
		"PWD":               config.UserWorkspace(username),
	})

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: rows, Cols: cols})
	if err != nil {
		log.Printf("PTY 创建失败: user=%s id=%s err=%v", username, id, err)
		return nil, err
	}

	sess := &PtySession{
		ID:       id,
		UserName: username,
		Cmd:      cmd,
		Pty:      ptmx,
		Buf:      NewRingBuffer(maxScrollback),
		CreateAt: time.Now(),
		done:     make(chan struct{}),
	}

	m.mu.Lock()
	m.sessions[id] = sess
	m.mu.Unlock()

	// 单一 reader goroutine：读 PTY → 写 buffer + 广播给订阅者
	go func() {
		defer close(sess.done)
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				sess.Buf.Write(buf[:n])
				sess.broadcast(buf[:n])
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("PTY 读取结束: id=%s err=%v", id, err)
				}
				break
			}
		}
	}()

	return sess, nil
}

func (m *PtyManager) Get(id string) *PtySession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[id]
}

func (m *PtyManager) List(username string) []*PtySession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []*PtySession
	for _, s := range m.sessions {
		if s.UserName == username {
			out = append(out, s)
		}
	}
	return out
}

func (m *PtyManager) Resize(id string, rows, cols uint16) error {
	sess := m.Get(id)
	if sess == nil {
		return nil
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	return pty.Setsize(sess.Pty, &pty.Winsize{Rows: rows, Cols: cols})
}

func (m *PtyManager) Close(id string) {
	m.mu.Lock()
	sess, ok := m.sessions[id]
	if ok {
		delete(m.sessions, id)
	}
	m.mu.Unlock()

	if sess != nil {
		sess.Pty.Close()
		sess.Cmd.Process.Kill()
		sess.Cmd.Wait()
	}
}

func (s *PtySession) Write(data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Pty.Write(data)
}
