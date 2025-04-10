package services

import (
	"errors"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
)

type AuthorService struct {
	repository *repositories.AuthorRepository
}

func NewAuthorService() *AuthorService {
	return &AuthorService{
		repository: repositories.NewAuthorRepository(),
	}
}

func (s *AuthorService) Create(r *request.AuthorUpsert) error {
	author := models.Author{
		Name: *r.Name,
	}
	if one := s.repository.FindOne(map[string]any{"name": author.Name}); one.ID != 0 {
		return errors.New("name exists")
	}
	if err := s.repository.Create(&author); err != nil {
		return err
	}
	return nil
}

func (s *AuthorService) Update(r *request.AuthorUpsert) error {
	author, err := s.repository.Get(*r.Id)
	if err != nil {
		return err
	}
	if r.Name != nil {
		author.Name = *r.Name
	}
	if one := s.repository.FindOne(map[string]any{"name": author.Name}); one.ID != 0 && one.ID != author.ID {
		return errors.New("name exists")
	}
	if err = s.repository.Update(author); err != nil {
		return err
	}
	return nil
}

func (s *AuthorService) Delete(id uint) error {
	return s.repository.Delete(id)
}
func (s *AuthorService) Detail(id uint) (any, error) {
	author, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return author.ToMap(), nil
}
func (s *AuthorService) List(r *request.NormalSearch) (any, int64, error) {
	b := builder.New()
	if r.Keyword != nil {
		b.Like("name", *r.Keyword)
	}
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
	}
	b.Order(orderSet[r.Sort])
	authors, total, err := s.repository.ListWithCount(r.Offset, r.Limit, b)
	if err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(authors))
	for _, v := range authors {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}
