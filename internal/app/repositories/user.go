package repositories

import (
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/db"
)

type UserRepository struct {
	Repository[models.User]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Repository[models.User]{
			db: db.Db(),
		},
	}
}

func (r *UserRepository) GetByEmail(email string) *models.User {
	user := new(models.User)
	err := r.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil
	}
	return user
}
