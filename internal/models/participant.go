package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Participant struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null" binding:"required,min=2,max=100"`
	Place     string    `json:"place" gorm:"type:varchar(255);not null" binding:"required,min=2,max=255"`
	BirthDate time.Time `json:"birth_date" gorm:"type:date;not null" binding:"required"`
	Kampus    string    `json:"kampus" gorm:"type:varchar(255);not null" binding:"required"`
	Jurusan   string    `json:"jurusan" gorm:"type:varchar(255);not null" binding:"required,min=2,max=255"`
	Angkatan  string    `json:"angkatan" gorm:"type:varchar(255);not null" binding:"required,min=2,max=255"`
	Phone     string    `json:"phone" gorm:"type:varchar(20);not null" binding:"required,min=8,max=20"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook untuk generate UUID dan set timestamps
func (p *Participant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

// BeforeUpdate hook untuk update timestamp
func (p *Participant) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

func (Participant) TableName() string { return "participants" }
