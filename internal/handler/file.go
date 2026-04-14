package handler

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
)

// GetFileTree 获取文件树
func GetFileTree(c *gin.Context) {
	username, _ := c.Get("username")
	path := c.DefaultQuery("path", ".")
	depth := 3

	tree, err := service.FileTree(username.(string), path, depth)
	if err != nil {
		log.Printf("GetFileTree 失败: user=%s path=%s err=%v", username, path, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取目录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tree": tree})
}

// ListDir 列出目录内容
func ListDir(c *gin.Context) {
	username, _ := c.Get("username")
	path := c.DefaultQuery("path", ".")
	nodes, err := service.ListDir(username.(string), path)
	if err != nil {
		log.Printf("ListDir 失败: user=%s path=%s err=%v", username, path, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取目录失败"})
		return
	}
	if nodes == nil {
		nodes = []*service.FileNode{}
	}
	c.JSON(http.StatusOK, gin.H{"files": nodes})
}

// ReadFile 读取文件内容
func ReadFile(c *gin.Context) {
	username, _ := c.Get("username")
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 path 参数"})
		return
	}
	data, err := service.ReadFile(username.(string), path)
	if err != nil {
		log.Printf("ReadFile 失败: user=%s path=%s err=%v", username, path, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": string(data), "path": path})
}

// WriteFile 写入文件
func WriteFile(c *gin.Context) {
	username, _ := c.Get("username")
	var req struct {
		Path    string `json:"path" binding:"required"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}
	if err := service.WriteFile(username.(string), req.Path, []byte(req.Content)); err != nil {
		log.Printf("WriteFile 失败: user=%s path=%s err=%v", username, req.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已保存"})
}

// CreateDir 创建目录
func CreateDir(c *gin.Context) {
	username, _ := c.Get("username")
	var req struct {
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}
	if err := service.CreateDir(username.(string), req.Path); err != nil {
		log.Printf("CreateDir 失败: user=%s path=%s err=%v", username, req.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已创建"})
}

// DeleteFile 删除文件或目录
func DeleteFile(c *gin.Context) {
	username, _ := c.Get("username")
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 path 参数"})
		return
	}
	if err := service.Remove(username.(string), path); err != nil {
		log.Printf("DeleteFile 失败: user=%s path=%s err=%v", username, path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

// RenameFile 重命名
func RenameFile(c *gin.Context) {
	username, _ := c.Get("username")
	var req struct {
		OldPath string `json:"old_path" binding:"required"`
		NewPath string `json:"new_path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}
	if err := service.Rename(username.(string), req.OldPath, req.NewPath); err != nil {
		log.Printf("RenameFile 失败: user=%s old=%s new=%s err=%v", username, req.OldPath, req.NewPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重命名失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已重命名"})
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	username, _ := c.Get("username")
	dir := c.DefaultPostForm("path", ".")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有文件"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	defer f.Close()

	relPath := filepath.Join(dir, file.Filename)
	if err := service.SaveUpload(username.(string), relPath, f); err != nil {
		log.Printf("UploadFile 失败: user=%s path=%s err=%v", username, relPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已上传", "path": relPath})
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {
	username, _ := c.Get("username")
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 path 参数"})
		return
	}
	full, err := service.SafePath(username.(string), path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法路径"})
		return
	}
	c.File(full)
}
