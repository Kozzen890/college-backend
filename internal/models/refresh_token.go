package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index"`
	Token     string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RefreshToken) TableName() string { return "refresh_tokens" }
