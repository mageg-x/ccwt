package handler

import (
	"log"
	"net/http"

	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
)

// GetHistoryProjects 获取用户的历史项目列表
func GetHistoryProjects(c *gin.Context) {
	username, _ := c.Get("username")
	projects, err := service.ListProjects(username.(string))
	if err != nil {
		log.Printf("GetHistoryProjects 失败: user=%s err=%v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取历史失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// GetHistorySession 获取单个会话详情
func GetHistorySession(c *gin.Context) {
	username, _ := c.Get("username")
	file := c.Query("file")
	if file == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 file 参数"})
		return
	}
	entries, err := service.ReadSession(username.(string), file)
	if err != nil {
		log.Printf("GetHistorySession 失败: user=%s file=%s err=%v", username, file, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"entries": entries})
}
