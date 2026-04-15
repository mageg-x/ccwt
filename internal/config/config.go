package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// Cfg 全局配置实例
var (
	Cfg  *AppConfig
	once sync.Once
)

type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Register RegisterConfig `mapstructure:"register"`
	Proxy    ProxyConfig    `mapstructure:"proxy"`
	Voice    VoiceConfig    `mapstructure:"voice"`
}

type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	DataDir string `mapstructure:"data_dir"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	AccessExpire  int    `mapstructure:"access_expire"`  // 分钟
	RefreshExpire int    `mapstructure:"refresh_expire"` // 天
}

type RegisterConfig struct {
	Open       bool   `mapstructure:"open"`        // 是否开放注册
	InviteCode string `mapstructure:"invite_code"` // 邀请码（为空则不限制）
}

type ProxyConfig struct {
	Port int `mapstructure:"port"` // SOCKS5 代理端口
}

type VoiceConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// DataDir 获取数据根目录
func DataDir() string {
	if Cfg != nil && Cfg.Server.DataDir != "" {
		return Cfg.Server.DataDir
	}
	if d := os.Getenv("CCWT_DATA_DIR"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ccwt")
}

// UserDir 获取用户专属目录
func UserDir(username string) string {
	return filepath.Join(DataDir(), "users", username)
}

// UserClaudeDir 获取用户 Claude 配置目录
func UserClaudeDir(username string) string {
	return filepath.Join(UserDir(username), ".claude")
}

// UserWorkspace 获取用户工作区目录
func UserWorkspace(username string) string {
	return filepath.Join(UserDir(username), "workspace")
}

// SetInviteCode 设置邀请码（命令行参数优先）
func SetInviteCode(code string) {
	viper.Set("register.invite_code", code)
	if Cfg != nil {
		Cfg.Register.InviteCode = code
	}
}

// SetAddr 设置监听地址（命令行参数优先）
func SetAddr(addr string) {
	if Cfg != nil {
		var host string
		var port int
		if _, err := fmt.Sscanf(addr, "%[^:]:%d", &host, &port); err == nil {
			Cfg.Server.Host = host
			Cfg.Server.Port = port
		}
	}
}

// Init 初始化配置
func Init() {
	once.Do(func() {
		Cfg = &AppConfig{}

		// 默认值
		viper.SetDefault("server.host", "0.0.0.0")
		viper.SetDefault("server.port", 3000)
		viper.SetDefault("server.data_dir", "")
		viper.SetDefault("server.log_level", "info")
		viper.SetDefault("jwt.secret", "ccwt-secret-change-me")
		viper.SetDefault("jwt.access_expire", 60*24) // 24小时
		viper.SetDefault("jwt.refresh_expire", 30)   // 30天
		viper.SetDefault("register.open", true)
		viper.SetDefault("register.invite_code", "")
		viper.SetDefault("proxy.port", 1080)
		viper.SetDefault("voice.enabled", true)

		// 环境变量前缀
		viper.SetEnvPrefix("CCWT")
		viper.AutomaticEnv()

		// 配置文件
		cfgDir := DataDir()
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(cfgDir)
		viper.AddConfigPath(".")

		_ = viper.ReadInConfig()
		_ = viper.Unmarshal(Cfg)

		// 确保数据目录存在
		os.MkdirAll(filepath.Join(DataDir(), "users"), 0755)
		os.MkdirAll(filepath.Join(DataDir(), "logs"), 0755)
	})
}
