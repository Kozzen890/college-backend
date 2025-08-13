package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"type:varchar(255);not null;unique" binding:"required"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
