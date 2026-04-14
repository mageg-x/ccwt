package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ccwt/ccwt/internal/config"
)

// FileNode 文件树节点
type FileNode struct {
	Name  string      `json:"name"`
	Path  string      `json:"path"` // 相对于 workspace 的路径
	IsDir bool        `json:"is_dir"`
	Size  int64       `json:"size,omitempty"`
	Kids  []*FileNode `json:"children,omitempty"`
}

func isWithin(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// SafePath 安全地解析用户路径，防止目录穿越
func SafePath(username, relPath string) (string, error) {
	workspace := config.UserWorkspace(username)
	wsAbs, err := filepath.Abs(workspace)
	if err != nil {
		return "", err
	}

	clean := filepath.Clean(relPath)
	if clean == "." {
		return wsAbs, nil
	}

	full := filepath.Join(wsAbs, clean)
	abs, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if !isWithin(wsAbs, abs) {
		return "", fmt.Errorf("非法路径: %s", relPath)
	}
	return abs, nil
}

// ListDir 列出目录内容（非递归）
func ListDir(username, relPath string) ([]*FileNode, error) {
	dir, err := SafePath(username, relPath)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var nodes []*FileNode
	for _, e := range entries {
		info, _ := e.Info()
		node := &FileNode{
			Name:  e.Name(),
			Path:  filepath.Join(relPath, e.Name()),
			IsDir: e.IsDir(),
		}
		if info != nil && !e.IsDir() {
			node.Size = info.Size()
		}
		nodes = append(nodes, node)
	}

	// 目录在前，文件在后，各自按名称排序
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].IsDir != nodes[j].IsDir {
			return nodes[i].IsDir
		}
		return nodes[i].Name < nodes[j].Name
	})

	return nodes, nil
}

// FileTree 递归获取文件树（限制深度）
func FileTree(username, relPath string, depth int) (*FileNode, error) {
	dir, err := SafePath(username, relPath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	root := &FileNode{
		Name:  info.Name(),
		Path:  relPath,
		IsDir: info.IsDir(),
	}
	if info.IsDir() && depth > 0 {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			childPath := filepath.Join(relPath, e.Name())
			child, _ := FileTree(username, childPath, depth-1)
			if child != nil {
				root.Kids = append(root.Kids, child)
			}
		}
		// 排序
		sort.Slice(root.Kids, func(i, j int) bool {
			if root.Kids[i].IsDir != root.Kids[j].IsDir {
				return root.Kids[i].IsDir
			}
			return root.Kids[i].Name < root.Kids[j].Name
		})
	}
	return root, nil
}

// ReadFile 读取文件内容
func ReadFile(username, relPath string) ([]byte, error) {
	full, err := SafePath(username, relPath)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(full)
}

// WriteFile 写入文件内容
func WriteFile(username, relPath string, data []byte) error {
	full, err := SafePath(username, relPath)
	if err != nil {
		return err
	}
	dir := filepath.Dir(full)
	os.MkdirAll(dir, 0755)
	return os.WriteFile(full, data, 0644)
}

// CreateDir 创建目录
func CreateDir(username, relPath string) error {
	full, err := SafePath(username, relPath)
	if err != nil {
		return err
	}
	return os.MkdirAll(full, 0755)
}

// Remove 删除文件或目录
func Remove(username, relPath string) error {
	full, err := SafePath(username, relPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(full)
}

// Rename 重命名
func Rename(username, oldPath, newPath string) error {
	oldFull, err := SafePath(username, oldPath)
	if err != nil {
		return err
	}
	newFull, err := SafePath(username, newPath)
	if err != nil {
		return err
	}
	return os.Rename(oldFull, newFull)
}

// SaveUpload 保存上传的文件
func SaveUpload(username, relPath string, reader io.Reader) error {
	full, err := SafePath(username, relPath)
	if err != nil {
		return err
	}
	dir := filepath.Dir(full)
	os.MkdirAll(dir, 0755)
	f, err := os.Create(full)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}
