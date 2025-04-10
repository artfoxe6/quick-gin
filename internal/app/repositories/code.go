package repositories

import (
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type CodeRepository struct {
	Repository[models.Code]
}

func NewCodeRepository() *CodeRepository {
	return &CodeRepository{
		Repository[models.Code]{
			db: db.Db(),
		},
	}
}
