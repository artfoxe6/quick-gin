package model

import (
	"github.com/artfoxe6/quick-gin/internal/app/core/model"
)

type Code struct {
	model.Base
	Email string `gorm:"size:255"`
	Code  string `gorm:"size:255"`
	Type  int
}

func (m *Code) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"email":      m.Email,
		"code":       m.Code,
		"type":       m.Type,
		"created_at": m.CreatedAt,
	}
}
