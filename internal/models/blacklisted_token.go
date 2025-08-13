package models

import (
	"time"

	"gorm.io/gorm"
)

// BlacklistedToken model untuk menyimpan token yang sudah logout
type BlacklistedToken struct {
	ID        uint      `gorm:"primaryKey"`
	TokenHash string    `gorm:"uniqueIndex;not null;size:64"` // SHA256 hash dari token
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`   // Waktu expired token
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook untuk set timestamps
func (bt *BlacklistedToken) BeforeCreate(tx *gorm.DB) error {
	bt.CreatedAt = time.Now()
	bt.UpdatedAt = time.Now()

	// Set default ExpiresAt jika belum diset
	if bt.ExpiresAt.IsZero() {
		bt.ExpiresAt = time.Now().Add(24 * time.Hour) // 24 jam dari sekarang
	}

	return nil
}

// BeforeUpdate hook untuk update timestamp
func (bt *BlacklistedToken) BeforeUpdate(tx *gorm.DB) error {
	bt.UpdatedAt = time.Now()
	return nil
}
