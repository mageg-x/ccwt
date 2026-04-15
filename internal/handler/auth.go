package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ccwt/ccwt/internal/config"
	"github.com/ccwt/ccwt/internal/db"
	"github.com/ccwt/ccwt/internal/middleware"
	"github.com/ccwt/ccwt/internal/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ensureUserBootstrapFiles(username string) {
	userHome := config.UserDir(username)
	files := map[string]string{
		filepath.Join(userHome, ".bashrc"): `# ~/.bashrc (CCWT user scope)
# User custom shell init for CCWT terminal.
# Example:
# export PATH="$HOME/bin:$PATH"
`,
		filepath.Join(userHome, ".profile"): `# ~/.profile (CCWT user scope)
# Load interactive bash settings when applicable.
if [ -n "$BASH_VERSION" ] && [ -f "$HOME/.bashrc" ]; then
  . "$HOME/.bashrc"
fi
`,
	}

	for path, content := range files {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if werr := os.WriteFile(path, []byte(content), 0600); werr != nil {
				log.Printf("创建用户初始化脚本失败: user=%s file=%s err=%v", username, path, werr)
			}
		}
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
		return
	}

	if !config.Cfg.Register.Open {
		c.JSON(http.StatusForbidden, gin.H{"error": "注册已关闭"})
		return
	}

	if config.Cfg.Register.InviteCode != "" && req.InviteCode != config.Cfg.Register.InviteCode {
		c.JSON(http.StatusForbidden, gin.H{"error": "邀请码错误"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Register bcrypt失败: user=%s err=%v", req.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 首个用户自动成为管理员
	role := "user"
	var count int
	db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		role = "admin"
	}

	res, err := db.DB.Exec(
		"INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		req.Username, string(hash), role,
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	uid, _ := res.LastInsertId()

	// 创建用户专属目录
	dirs := []string{
		config.UserClaudeDir(req.Username),
		config.UserWorkspace(req.Username),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0700); err != nil {
			log.Printf("Register 创建用户目录失败: user=%s dir=%s err=%v", req.Username, d, err)
		}
	}
	ensureUserBootstrapFiles(req.Username)

	// 自动登录
	token, _ := middleware.GenToken(uid, req.Username, role)
	setTokenCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"user": model.User{
			ID: uid, Username: req.Username, Role: role,
			CreateAt: time.Now(),
		},
		"token": token,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}

	var user model.User
	var passHash string
	err := db.DB.QueryRow(
		"SELECT id, username, password_hash, role, created_at FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &passHash, &user.Role, &user.CreateAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	if err != nil {
		log.Printf("Login 查询用户失败: user=%s err=%v", req.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 旧用户补齐初始化脚本，避免历史账户缺失 ~/.bashrc ~/.profile
	ensureUserBootstrapFiles(user.Username)

	token, _ := middleware.GenToken(user.ID, user.Username, user.Role)

	// 保存 refresh token
	rt := genRefreshToken()
	expire := time.Now().AddDate(0, 0, config.Cfg.JWT.RefreshExpire)
	db.DB.Exec(
		"INSERT INTO sessions (user_id, refresh_token, expires_at) VALUES (?, ?, ?)",
		user.ID, rt, expire,
	)

	setTokenCookie(c, token)
	c.SetCookie("refresh_token", rt, config.Cfg.JWT.RefreshExpire*86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

// RefreshToken 刷新 access token
func RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "需要重新登录"})
		return
	}

	var uid int64
	var expire time.Time
	err = db.DB.QueryRow(
		"SELECT user_id, expires_at FROM sessions WHERE refresh_token = ?", rt,
	).Scan(&uid, &expire)
	if err != nil || time.Now().After(expire) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "会话过期，请重新登录"})
		return
	}

	var user model.User
	db.DB.QueryRow(
		"SELECT id, username, role FROM users WHERE id = ?", uid,
	).Scan(&user.ID, &user.Username, &user.Role)

	token, _ := middleware.GenToken(user.ID, user.Username, user.Role)
	setTokenCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

// Logout 登出
func Logout(c *gin.Context) {
	rt, _ := c.Cookie("refresh_token")
	if rt != "" {
		db.DB.Exec("DELETE FROM sessions WHERE refresh_token = ?", rt)
	}
	c.SetCookie("token", "", -1, "/", "", false, false)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "已登出"})
}

// GetMe 获取当前用户信息
func GetMe(c *gin.Context) {
	uid, _ := c.Get("uid")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       uid,
			"username": username,
			"role":     role,
		},
	})
}

func setTokenCookie(c *gin.Context, token string) {
	maxAge := config.Cfg.JWT.AccessExpire * 60
	c.SetCookie("token", token, maxAge, "/", "", false, false)
}

func genRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
