package model

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ToArray interface {
	ToMap() map[string]any
}
