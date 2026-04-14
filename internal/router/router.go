package router

import (
	"io/fs"
	"net/http"

	"github.com/ccwt/ccwt/internal/handler"
	"github.com/ccwt/ccwt/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由
func Setup(r *gin.Engine, webFS fs.FS) {
	// API 路由
	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)
		api.POST("/auth/refresh", handler.RefreshToken)

		// 需要认证
		auth := api.Group("", middleware.AuthRequired())
		{
			auth.POST("/auth/logout", handler.Logout)
			auth.GET("/auth/me", handler.GetMe)

			// 终端
			auth.GET("/terminals", handler.ListTerminals)
			auth.DELETE("/terminals/:id", handler.CloseTerminal)

			// 文件系统
			auth.GET("/files/tree", handler.GetFileTree)
			auth.GET("/files/list", handler.ListDir)
			auth.GET("/files/read", handler.ReadFile)
			auth.POST("/files/write", handler.WriteFile)
			auth.POST("/files/mkdir", handler.CreateDir)
			auth.DELETE("/files", handler.DeleteFile)
			auth.POST("/files/rename", handler.RenameFile)
			auth.POST("/files/upload", handler.UploadFile)
			auth.GET("/files/download", handler.DownloadFile)

			// 会话历史
			auth.GET("/history/projects", handler.GetHistoryProjects)
			auth.GET("/history/session", handler.GetHistorySession)

			// 代理
			auth.GET("/proxy/status", handler.GetProxyStatus)
			auth.POST("/proxy/start", handler.StartProxy)
			auth.POST("/proxy/stop", handler.StopProxy)

			// 系统信息
			auth.GET("/system/info", handler.GetSystemInfo)

			// 语音
			auth.POST("/voice/recognize", handler.VoiceRecognize)
			auth.GET("/voice/status", handler.VoiceStatus)

			// 管理员接口
			admin := auth.Group("/admin", middleware.AdminRequired())
			{
				admin.GET("/users", handler.ListUsers)
				admin.DELETE("/users/:id", handler.DeleteUser)
				admin.PUT("/users/:id/role", handler.UpdateUserRole)
			}
		}
	}

	// WebSocket 终端（需要 token 认证，在 handler 内处理）
	r.GET("/ws/terminal", handler.TerminalWS)

	// 前端静态资源
	if webFS != nil {
		// 服务前端文件
		fileServer := http.FileServer(http.FS(webFS))
		r.NoRoute(func(c *gin.Context) {
			// 尝试静态文件
			path := c.Request.URL.Path
			f, err := webFS.Open(path[1:]) // 去掉开头的 /
			if err == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			// SPA fallback: 返回 index.html
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	}
}
