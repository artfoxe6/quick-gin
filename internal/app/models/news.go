package models

import (
	"github.com/artfoxe6/quick-gin/internal/pkg/kit"
	"gorm.io/gorm"
)

const (
	NEWS_STATUS_DRAFT     = 0
	NEWS_STATUS_PUBLISHED = 1
)

type News struct {
	gorm.Model
	Title       string `gorm:"size:255"`
	Content     string `gorm:"type:text"`
	AuthorId    *uint
	Author      *Author
	UserId      *uint
	User        *User
	Source      string      `gorm:"size:255"`
	Status      int         `gorm:"default:0"`
	Views       int         `gorm:"default:0"`
	Likes       int         `gorm:"default:0"`
	Comments    int         `gorm:"default:0"`
	IsFeatured  int64       `gorm:"default:0"`
	Slug        string      `gorm:"size:255"`
	Summary     string      `gorm:"size:255"`
	ImageUrl    string      `gorm:"size:255"`
	VideoUrl    string      `gorm:"size:255"`
	AudioUrl    string      `gorm:"size:255"`
	TypeId      int         `gorm:"default:0"`
	Categories  []*Category `gorm:"many2many:news_category"`
	Tags        []*Tag      `gorm:"many2many:news_tag"`
	PublishTime int64
}

func (n News) GetSlug() string {
	if n.Slug != "" {
		return n.Slug
	}
	return kit.Slug(n.Title)
}

func (n News) ToMap() map[string]any {
	isFeatured := false
	if n.IsFeatured > 0 {
		isFeatured = true
	}
	data := map[string]any{
		"id":          n.ID,
		"type_id":     n.TypeId,
		"title":       n.Title,
		"content":     n.Content,
		"source":      n.Source,
		"status":      n.Status,
		"is_featured": isFeatured,
		"summary":     n.Summary,
		"image_url":   n.ImageUrl,
		"video_url":   n.VideoUrl,
		"audio_url":   n.AudioUrl,
		"slug":        n.Slug,
		"created_at":  n.CreatedAt,
		"updated_at":  n.UpdatedAt,
		"author":      nil,
		"tags":        nil,
		"categories":  nil,
	}
	if n.Author != nil {
		data["author"] = n.Author.ToMap()
	}
	if n.Tags != nil {
		tags := []map[string]any{}
		for _, tag := range n.Tags {
			tags = append(tags, tag.ToMap())
		}
		data["tags"] = tags
	}
	if n.Categories != nil {
		categories := []map[string]any{}
		for _, category := range n.Categories {
			categories = append(categories, category.ToMap())
		}
		data["categories"] = categories
	}
	return data
}
