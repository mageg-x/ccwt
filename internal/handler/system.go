package handler

import (
	"net/http"

	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
)

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	info := service.GetSystemInfo()
	c.JSON(http.StatusOK, gin.H{"system": info})
}
