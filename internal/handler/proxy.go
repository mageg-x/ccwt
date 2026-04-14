package handler

import (
	"net/http"

	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
)

// GetProxyStatus 获取代理状态
func GetProxyStatus(c *gin.Context) {
	running, addr := service.Proxy.Status()
	c.JSON(http.StatusOK, gin.H{
		"running": running,
		"address": addr,
	})
}

// StartProxy 启动代理
func StartProxy(c *gin.Context) {
	if err := service.Proxy.Start(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	running, addr := service.Proxy.Status()
	c.JSON(http.StatusOK, gin.H{
		"message": "代理已启动",
		"running": running,
		"address": addr,
	})
}

// StopProxy 停止代理
func StopProxy(c *gin.Context) {
	if err := service.Proxy.Stop(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "代理已停止", "running": false})
}
