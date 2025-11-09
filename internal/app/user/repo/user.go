package repo

import (
	"github.com/artfoxe6/quick-gin/internal/app/core/repository"
	"github.com/artfoxe6/quick-gin/internal/app/user/model"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type UserRepository struct {
	*repository.Repository[model.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Repository: repository.New[model.User](db.Db()),
	}
}

func (r *UserRepository) GetByEmail(email string) *model.User {
	user := new(model.User)
	err := r.DB().Where("email = ?", email).First(user).Error
	if err != nil {
		return nil
	}
	return user
}
