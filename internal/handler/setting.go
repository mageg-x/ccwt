package handler

import (
	"net/http"
	"strings"

	"github.com/ccwt/ccwt/internal/db"
	"github.com/ccwt/ccwt/internal/model"
	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, key, value, COALESCE(description, ''), updated_at FROM settings ORDER BY key")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设置失败"})
		return
	}
	defer rows.Close()

	var settings []model.Setting
	for rows.Next() {
		var s model.Setting
		if err := rows.Scan(&s.ID, &s.Key, &s.Value, &s.Description, &s.UpdateAt); err != nil {
			continue
		}
		settings = append(settings, s)
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func UpdateSetting(c *gin.Context) {
	var req model.SettingsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	key := strings.TrimSpace(req.Key)
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "设置项不能为空"})
		return
	}

	value := strings.TrimSpace(req.Value)

	_, err := db.DB.Exec(
		"UPDATE settings SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE key = ?",
		value, key,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已更新", "key": key, "value": value})
}

func GetSetting(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 key 参数"})
		return
	}

	var value string
	err := db.DB.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设置项不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}