package repositories

import (
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type AuthorRepository struct {
	Repository[models.Author]
}

func NewAuthorRepository() *AuthorRepository {
	return &AuthorRepository{
		Repository[models.Author]{
			db: db.Db(),
		},
	}
}
