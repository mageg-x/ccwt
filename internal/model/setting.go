package model

import "time"

type Setting struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	UpdateAt    time.Time `json:"updated_at"`
}

type SettingsReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}