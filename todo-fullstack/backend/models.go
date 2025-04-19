package main

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
}

type Todo struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Task   string `json:"task"`
	Done   bool   `json:"done"`
	UserID uint   `json:"user_id"`
}
