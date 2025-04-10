package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string `gorm:"size:255"`
	Status int    `gorm:"default:0"`
}

func (m *Category) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"name":       m.Name,
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
	}
}
