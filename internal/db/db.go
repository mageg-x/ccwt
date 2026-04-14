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
	}
	for _, ddl := range tables {
		if _, err := DB.Exec(ddl); err != nil {
			log.Fatalf("建表失败: %v", err)
		}
	}
}
