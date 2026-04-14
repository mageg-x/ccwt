package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"

	"github.com/ccwt/ccwt/internal/config"
	"github.com/ccwt/ccwt/internal/db"
	"github.com/ccwt/ccwt/internal/router"
	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var webEmbed embed.FS

func main() {
	code := flag.String("code", "", "邀请码")
	flag.Parse()

	config.Init()

	if *code != "" {
		config.SetInviteCode(*code)
	}

	// 初始化数据库
	db.Init()

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 前端静态资源
	webFS, err := fs.Sub(webEmbed, "web/dist")
	if err != nil {
		log.Printf("前端资源加载失败（开发模式）: %v", err)
		webFS = nil
	}

	// 注册路由
	router.Setup(r, webFS)

	// 启动服务
	addr := fmt.Sprintf("0.0.0.0:%d", config.Cfg.Server.Port)
	log.Printf("CCWT 服务启动: http://%s", addr)
	log.Printf("数据目录: %s", config.DataDir())
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
