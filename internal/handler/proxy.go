package handler

import (
	"net/http"

	"github.com/ccwt/ccwt/internal/service"
	"github.com/gin-gonic/gin"
)

// GetProxyStatus 获取代理状态
func GetProxyStatus(c *gin.Context) {
	running, addr, bindHost, port := service.Proxy.Status()
	clientIP := "127.0.0.1"
	if bindHost != "" && bindHost != "0.0.0.0" && bindHost != "::" {
		clientIP = bindHost
	}
	c.JSON(http.StatusOK, gin.H{
		"running":    running,
		"address":    addr,
		"client_ip":  clientIP,
		"bind_host":  bindHost,
		"port":       port,
	})
}

// StartProxy 启动代理
func StartProxy(c *gin.Context) {
	var req struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := service.Proxy.Start(req.Host, req.Port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	running, addr, bindHost, port := service.Proxy.Status()
	clientIP := "127.0.0.1"
	if bindHost != "" && bindHost != "0.0.0.0" && bindHost != "::" {
		clientIP = bindHost
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "代理已启动",
		"running":   running,
		"address":   addr,
		"client_ip": clientIP,
		"bind_host": bindHost,
		"port":      port,
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
