package model

import "time"

type Session struct {
	ID      int64     `json:"id"`
	UserID  int64     `json:"user_id"`
	Token   string    `json:"refresh_token"`
	Expire  time.Time `json:"expires_at"`
	CreateAt time.Time `json:"created_at"`
}
