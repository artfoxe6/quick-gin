package repositories

import (
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type TagRepository struct {
	Repository[models.Tag]
}

func NewTagRepository() *TagRepository {
	return &TagRepository{
		Repository[models.Tag]{
			db: db.Db(),
		},
	}
}
