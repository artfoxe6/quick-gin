package model

import (
	"github.com/artfoxe6/quick-gin/internal/app/core/model"
)

type User struct {
	model.Base
	Name     string `gorm:"size:255;index:idx_name"`
	Password string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	Role     string `gorm:"size:255"`
	Avatar   string `gorm:"size:255"`
	LoginAt  int64
	LastIp   string `gorm:"size:255"`
}

func (m *User) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"name":       m.Name,
		"email":      m.Email,
		"role":       m.Role,
		"avatar":     m.Avatar,
		"login_at":   m.LoginAt,
		"last_ip":    m.LastIp,
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
	}
}

func (m User) TokenData() map[string]any {
	return map[string]any{
		"id":     m.ID,
		"name":   m.Name,
		"email":  m.Email,
		"role":   m.Role,
		"avatar": m.Avatar,
	}
}
