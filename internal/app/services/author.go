package services

import (
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
)

type AuthorRepository interface {
	FindOne(map[string]any, ...*builder.Builder) *models.Author
	Create(*models.Author) error
	Update(*models.Author) error
	Delete(uint) error
	Get(uint, ...*builder.Builder) (*models.Author, error)
	ListWithCount(int, int, ...*builder.Builder) ([]models.Author, int64, error)
}

type AuthorService interface {
	Create(r *request.AuthorUpsert) error
	Update(r *request.AuthorUpsert) error
	Delete(id uint) error
	Detail(id uint) (any, error)
	List(r *request.NormalSearch) (any, int64, error)
}

type authorService struct {
	repository AuthorRepository
}

func NewAuthorService(repository AuthorRepository) AuthorService {
	return &authorService{repository: repository}
}

func (s *authorService) Create(r *request.AuthorUpsert) error {
	author := models.Author{
		Name: *r.Name,
	}
	if one := s.repository.FindOne(map[string]any{"name": author.Name}); one.ID != 0 {
		return apperr.BadRequest("name exists")
	}
	if err := s.repository.Create(&author); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *authorService) Update(r *request.AuthorUpsert) error {
	author, err := s.repository.Get(*r.Id)
	if err != nil {
		return apperr.Internal(err)
	}
	if r.Name != nil {
		author.Name = *r.Name
	}
	if one := s.repository.FindOne(map[string]any{"name": author.Name}); one.ID != 0 && one.ID != author.ID {
		return apperr.BadRequest("name exists")
	}
	if err = s.repository.Update(author); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *authorService) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *authorService) Detail(id uint) (any, error) {
	author, err := s.repository.Get(id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return author.ToMap(), nil
}

func (s *authorService) List(r *request.NormalSearch) (any, int64, error) {
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
		return nil, 0, apperr.Internal(err)
	}
	list := make([]map[string]any, 0, len(authors))
	for _, v := range authors {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}
