package service

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ccwt/ccwt/internal/config"
)

// HistoryEntry 会话历史条目
type HistoryEntry struct {
	Type      string    `json:"type"`
	Message   any       `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// HistoryProject 项目的历史会话
type HistoryProject struct {
	Project  string           `json:"project"`
	Sessions []HistorySession `json:"sessions"`
}

// HistorySession 单次会话
type HistorySession struct {
	File    string         `json:"file"`
	Entries []HistoryEntry `json:"entries"`
	ModTime time.Time      `json:"mod_time"`
}

func isWithinPath(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// ListProjects 列出用户所有有历史记录的项目
func ListProjects(username string) ([]HistoryProject, error) {
	projDir := filepath.Join(config.UserClaudeDir(username), "projects")
	if _, err := os.Stat(projDir); os.IsNotExist(err) {
		return []HistoryProject{}, nil
	}

	projectMap := make(map[string][]HistorySession)

	// 遍历 projects 目录
	filepath.Walk(projDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".jsonl") {
			return nil
		}

		rel, _ := filepath.Rel(projDir, path)
		parts := strings.SplitN(rel, string(os.PathSeparator), 2)
		projName := parts[0]
		if len(parts) > 1 {
			projName = filepath.Join(parts[:len(parts)-1]...)
		}

		projectMap[projName] = append(projectMap[projName], HistorySession{
			File:    rel,
			ModTime: info.ModTime(),
		})
		return nil
	})

	var result []HistoryProject
	for proj, sessions := range projectMap {
		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].ModTime.After(sessions[j].ModTime)
		})
		result = append(result, HistoryProject{
			Project:  proj,
			Sessions: sessions,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if len(result[i].Sessions) == 0 || len(result[j].Sessions) == 0 {
			return false
		}
		return result[i].Sessions[0].ModTime.After(result[j].Sessions[0].ModTime)
	})

	return result, nil
}

// ReadSession 读取单个会话的 JSONL 文件
func ReadSession(username, relFile string) ([]HistoryEntry, error) {
	projDir := filepath.Join(config.UserClaudeDir(username), "projects")
	projAbs, err := filepath.Abs(projDir)
	if err != nil {
		return nil, err
	}

	full := filepath.Join(projAbs, filepath.Clean(relFile))
	abs, err := filepath.Abs(full)
	if err != nil {
		return nil, err
	}

	if !isWithinPath(projAbs, abs) {
		return nil, os.ErrPermission
	}

	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []HistoryEntry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB per line
	for scanner.Scan() {
		var entry HistoryEntry
		if json.Unmarshal(scanner.Bytes(), &entry) == nil {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}
