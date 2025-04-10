package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"size:255"`
}

func (m *Tag) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"name":       m.Name,
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
	}
}
