package service

import (
	"bytes"
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
	// 环形覆盖后起始位置可能落在 ANSI 控制序列中间（如 "...[1;2m"），
	// 刷新重连回放时会把残片当普通文本显示。对满缓冲场景从下一行起回放，
	// 牺牲一小段最旧行内容，换取稳定可读的恢复效果。
	if nl := bytes.IndexByte(out, '\n'); nl >= 0 && nl+1 < len(out) {
		return append([]byte{}, out[nl+1:]...)
	}
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

func envToMap(base []string) map[string]string {
	out := make(map[string]string, len(base))
	for _, kv := range base {
		i := strings.IndexByte(kv, '=')
		if i <= 0 {
			continue
		}
		out[kv[:i]] = kv[i+1:]
	}
	return out
}

func mapToEnv(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func isAllowedPassthroughEnv(key string) bool {
	switch key {
	case "PATH", "LANG", "LANGUAGE", "TZ",
		"HTTP_PROXY", "HTTPS_PROXY", "NO_PROXY",
		"http_proxy", "https_proxy", "no_proxy",
		"SSL_CERT_FILE", "SSL_CERT_DIR":
		return true
	}
	return strings.HasPrefix(key, "LC_")
}

func filterEnvForSandbox(base []string) []string {
	m := envToMap(base)
	out := make(map[string]string)
	for k, v := range m {
		if isAllowedPassthroughEnv(k) {
			out[k] = v
		}
	}
	return mapToEnv(out)
}

func remapHomeInEnv(base []string, targetHome string) []string {
	sourceHome, _ := os.UserHomeDir()
	m := envToMap(base)
	delete(m, "PWD")
	delete(m, "OLDPWD")
	delete(m, "PROMPT_COMMAND")
	delete(m, "PS1")
	delete(m, "PS2")
	delete(m, "PS4")
	delete(m, "BASH_ENV")
	if sourceHome != "" && sourceHome != targetHome {
		for k, v := range m {
			m[k] = strings.ReplaceAll(v, sourceHome, targetHome)
		}
	}
	return mapToEnv(m)
}

func appendRoBindIfExists(args []string, hostPath string) []string {
	if _, err := os.Stat(hostPath); err == nil {
		args = append(args, "--ro-bind", hostPath, hostPath)
	}
	return args
}

func appendRoBindIntoIfExists(args []string, hostPath, sandboxPath string) []string {
	if _, err := os.Stat(hostPath); err == nil {
		args = append(args, "--ro-bind", hostPath, sandboxPath)
	}
	return args
}

func buildBubblewrapCommand(username, shell string) (*exec.Cmd, error) {
	bwrap, err := exec.LookPath("bwrap")
	if err != nil {
		return nil, err
	}

	userHome := config.UserDir(username)
	workspace := config.UserWorkspace(username)

	// 在沙箱内统一映射成固定路径，避免暴露宿主机路径结构
	sbHome := "/home/ccwt"
	sbWorkspace := sbHome + "/workspace"
	sbClaudeDir := sbHome + "/.claude"
	sbBashrc := sbHome + "/.ccwt_bashrc"
	sbHistory := sbHome + "/.bash_history"

	args := []string{
		"--die-with-parent",
		"--unshare-pid",
		"--unshare-ipc",
		"--unshare-uts",
		"--proc", "/proc",
		"--dev", "/dev",
		"--tmpfs", "/tmp",
		"--tmpfs", "/etc",
		"--bind", userHome, sbHome,
		"--chdir", sbWorkspace,
		"--setenv", "HOME", sbHome,
		"--setenv", "CLAUDE_CONFIG_DIR", sbClaudeDir,
		"--setenv", "HISTFILE", sbHistory,
		"--setenv", "__CCWT_WORKSPACE", sbWorkspace,
		"--setenv", "CCWT_USER", username,
		"--setenv", "USER", username,
		"--setenv", "LOGNAME", username,
		"--setenv", "TERM", "xterm-256color",
		"--setenv", "SHELL", shell,
	}
	args = appendRoBindIfExists(args, "/bin")
	args = appendRoBindIfExists(args, "/usr")
	args = appendRoBindIfExists(args, "/lib")
	args = appendRoBindIfExists(args, "/lib64")
	args = appendRoBindIntoIfExists(args, "/etc/ssl", "/etc/ssl")
	args = appendRoBindIntoIfExists(args, "/etc/ca-certificates", "/etc/ca-certificates")
	args = appendRoBindIntoIfExists(args, "/etc/alternatives", "/etc/alternatives")
	args = appendRoBindIntoIfExists(args, "/etc/resolv.conf", "/etc/resolv.conf")
	args = appendRoBindIntoIfExists(args, "/etc/hosts", "/etc/hosts")
	args = appendRoBindIntoIfExists(args, "/etc/nsswitch.conf", "/etc/nsswitch.conf")
	args = appendRoBindIntoIfExists(args, "/etc/manpath.config", "/etc/manpath.config")
	args = appendRoBindIntoIfExists(args, "/etc/man_db.conf", "/etc/man_db.conf")

	args = append(args, "--", shell, "--noprofile", "--rcfile", sbBashrc, "-i")

	cmd := exec.Command(bwrap, args...)
	cmd.Dir = workspace
	cmd.Env = mergeEnv(filterEnvForSandbox(os.Environ()), map[string]string{
		"HOME":              userHome,
		"CLAUDE_CONFIG_DIR": config.UserClaudeDir(username),
		"HISTFILE":          filepath.Join(userHome, ".bash_history"),
		"__CCWT_WORKSPACE":  workspace,
		"CCWT_USER":         username,
		"USER":              username,
		"LOGNAME":           username,
		"TERM":              "xterm-256color",
		"SHELL":             shell,
		"PWD":               workspace,
	})
	return cmd, nil
}

func writeBashInit(username string) (string, error) {
	userHome := config.UserDir(username)
	initPath := filepath.Join(userHome, ".ccwt_bashrc")

	content := `# CCWT managed bash init
export HOME="${HOME:-$PWD}"
export CLAUDE_CONFIG_DIR="${CLAUDE_CONFIG_DIR:-$HOME/.claude}"
export HISTFILE="${HISTFILE:-$HOME/.bash_history}"
export __CCWT_WORKSPACE="${__CCWT_WORKSPACE:-$HOME/workspace}"
export COLORTERM=truecolor
export TERM="${TERM:-xterm-256color}"
export HISTSIZE=10000
export HISTFILESIZE=20000
shopt -s histappend cmdhist checkwinsize
mkdir -p "$HOME" "$__CCWT_WORKSPACE" >/dev/null 2>&1 || true
touch "$HISTFILE" >/dev/null 2>&1 || true
export PROMPT_COMMAND="history -a; history -n"

# 兼容用户自定义初始化（如 nvm/goenv/sdkman 等）
if [ -f "$HOME/.bashrc" ]; then
  . "$HOME/.bashrc"
fi

# 彩色输出增强：目录、grep、分页器、提示符
if command -v dircolors >/dev/null 2>&1; then
  eval "$(dircolors -b 2>/dev/null)"
fi
if [ -z "${LS_COLORS:-}" ]; then
  export LS_COLORS='di=1;38;5;81:ln=38;5;117:so=38;5;213:pi=38;5;213:ex=1;38;5;120:bd=38;5;183:cd=38;5;183:su=37;41:sg=30;43:tw=30;42:ow=30;43'
fi
alias ls='ls --color=auto'
alias ll='ls -alF --color=auto'
alias la='ls -A --color=auto'
alias l='ls -CF --color=auto'
alias grep='grep --color=auto'
alias egrep='egrep --color=auto'
alias fgrep='fgrep --color=auto'
export LESS='-R'

__ccwt_guard_workspace() {
  case "$PWD/" in
    "$__CCWT_WORKSPACE"/|"$__CCWT_WORKSPACE"/*) ;;
    *)
      printf "\\033[33m[CCWT] 工作目录受限，已返回 workspace\\033[0m\\n"
      builtin cd "$__CCWT_WORKSPACE" || return
      ;;
  esac
}
cd() {
  builtin cd "$@" || return
  __ccwt_guard_workspace
}
pushd() { builtin pushd "$@" && __ccwt_guard_workspace; }
popd() { builtin popd "$@" && __ccwt_guard_workspace; }
PROMPT_COMMAND="__ccwt_guard_workspace; history -a; history -n"
case "$PWD/" in
  "$__CCWT_WORKSPACE"/|"$__CCWT_WORKSPACE"/*) ;;
  *) builtin cd "$__CCWT_WORKSPACE" ;;
esac
if [ -t 1 ]; then
  __CCWT_PROMPT_USER="${CCWT_USER:-${USER:-user}}"
  PS1="${__CCWT_PROMPT_USER}:\w\$ "
fi
`

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
	bashrc := filepath.Join(config.UserDir(username), ".bashrc")
	if _, err := os.Stat(bashrc); os.IsNotExist(err) {
		content := `# ~/.bashrc (CCWT user scope)
# Add user-specific shell init here.
`
		if err := os.WriteFile(bashrc, []byte(content), 0600); err != nil {
			return err
		}
	}
	profile := filepath.Join(config.UserDir(username), ".profile")
	if _, err := os.Stat(profile); os.IsNotExist(err) {
		content := `# ~/.profile (CCWT user scope)
if [ -n "$BASH_VERSION" ] && [ -f "$HOME/.bashrc" ]; then
  . "$HOME/.bashrc"
fi
`
		if err := os.WriteFile(profile, []byte(content), 0600); err != nil {
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
	if bwrapCmd, berr := buildBubblewrapCommand(username, shell); berr == nil {
		cmd = bwrapCmd
		log.Printf("PTY 隔离模式: user=%s mode=bwrap", username)
	} else {
		cmd.Dir = config.UserWorkspace(username)
		cmd.Env = mergeEnv(remapHomeInEnv(os.Environ(), config.UserDir(username)), map[string]string{
			"CLAUDE_CONFIG_DIR": config.UserClaudeDir(username),
			"HOME":              config.UserDir(username),
			"HISTFILE":          filepath.Join(config.UserDir(username), ".bash_history"),
			"__CCWT_WORKSPACE":  config.UserWorkspace(username),
			"CCWT_USER":         username,
			"USER":              username,
			"LOGNAME":           username,
			"TERM":              "xterm-256color",
			"SHELL":             shell,
			"PWD":               config.UserWorkspace(username),
		})
		log.Printf("PTY 隔离模式: user=%s mode=soft-shell-only reason=%v", username, berr)
	}

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
