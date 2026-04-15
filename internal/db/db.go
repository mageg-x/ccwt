package db

import (
	"database/sql"
	"log"
	"path/filepath"
	"sync"

	"github.com/ccwt/ccwt/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB   *sql.DB
	once sync.Once
)

// Init 初始化数据库连接并自动建表
func Init() {
	once.Do(func() {
		dbPath := filepath.Join(config.DataDir(), "ccwt.db")
		var err error
		DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
		if err != nil {
			log.Fatalf("打开数据库失败: %v", err)
		}
		DB.SetMaxOpenConns(1) // SQLite 单写

		migrate()
	})
}

func migrate() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			refresh_token TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			value TEXT NOT NULL,
			description TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, ddl := range tables {
		if _, err := DB.Exec(ddl); err != nil {
			log.Fatalf("建表失败: %v", err)
		}
	}
	seedSettings()
}

func seedSettings() {
	settings := []struct {
		key, value, desc string
	}{
		{"voice.enabled", "true", "是否启用语音识别功能"},
		{"voice.app_id", "", "百度语音 App ID"},
		{"voice.api_key", "", "百度语音 API Key"},
		{"voice.secret", "", "百度语音 Secret Key"},
		{"proxy.ip", "0.0.0.0", "SOCKS5 代理绑定 IP（0.0.0.0 表示所有网卡）"},
		{"proxy.port", "1080", "SOCKS5 代理默认端口"},
	}
	for _, s := range settings {
		_, _ = DB.Exec("INSERT OR IGNORE INTO settings (key, value, description) VALUES (?, ?, ?)",
			s.key, s.value, s.desc)
	}
}
