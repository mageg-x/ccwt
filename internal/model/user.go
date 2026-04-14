package model

import "time"

type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	PassHash string    `json:"-"`
	Role     string    `json:"role"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type RegisterReq struct {
	Username   string `json:"username" binding:"required,min=2,max=32,alphanum"`
	Password   string `json:"password" binding:"required,min=6,max=64"`
	InviteCode string `json:"invite_code"`
}
