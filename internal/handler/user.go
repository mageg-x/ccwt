package handler

import (
	"log"
	"net/http"

	"github.com/ccwt/ccwt/internal/db"
	"github.com/gin-gonic/gin"
)

// ListUsers 管理员获取用户列表
func ListUsers(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, username, role, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		log.Printf("ListUsers 查询失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var id int64
		var username, role, createdAt, updatedAt string
		rows.Scan(&id, &username, &role, &createdAt, &updatedAt)
		users = append(users, gin.H{
			"id": id, "username": username, "role": role,
			"created_at": createdAt, "updated_at": updatedAt,
		})
	}
	if users == nil {
		users = []gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// DeleteUser 管理员删除用户
func DeleteUser(c *gin.Context) {
	uid := c.Param("id")
	currentUID, _ := c.Get("uid")
	if uid == "1" || uid == currentUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己或超级管理员"})
		return
	}
	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", uid)
	if err != nil {
		log.Printf("DeleteUser 失败: id=%s err=%v", uid, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

// UpdateUserRole 管理员修改用户角色
func UpdateUserRole(c *gin.Context) {
	uid := c.Param("id")
	var req struct {
		Role string `json:"role" binding:"required,oneof=admin user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}
	_, err := db.DB.Exec("UPDATE users SET role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Role, uid)
	if err != nil {
		log.Printf("UpdateUserRole 失败: id=%s err=%v", uid, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已更新"})
}
