package repositories

import (
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type CategoryRepository struct {
	Repository[models.Category]
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		Repository[models.Category]{
			db: db.Db(),
		},
	}
}
