package repo

import (
	"github.com/artfoxe6/quick-gin/internal/app/core/repository"
	"github.com/artfoxe6/quick-gin/internal/app/user/model"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type CodeRepository struct {
	*repository.Repository[model.Code]
}

func NewCodeRepository() *CodeRepository {
	return &CodeRepository{
		Repository: repository.New[model.Code](db.Db()),
	}
}
