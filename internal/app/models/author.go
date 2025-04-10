package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name string `gorm:"size:255"`
}

func (m *Author) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"name":       m.Name,
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
	}
}
