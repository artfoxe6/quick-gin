package repositories

import (
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type NewsRepository struct {
	Repository[models.News]
}

func NewNewsRepository() *NewsRepository {
	return &NewsRepository{
		Repository[models.News]{
			db: db.Db(),
		},
	}
}
