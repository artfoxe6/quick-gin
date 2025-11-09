package model

import "gorm.io/gorm"

type Code struct {
	gorm.Model
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
