package services

import (
	"errors"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
)

type CategoryService struct {
	repository *repositories.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		repository: repositories.NewCategoryRepository(),
	}
}

func (s *CategoryService) Create(r *request.CategoryUpsert) error {
	category := models.Category{
		Name: *r.Name,
	}
	if one := s.repository.FindOne(map[string]any{"name": category.Name}); one.ID != 0 {
		return errors.New("name exists")
	}
	if err := s.repository.Create(&category); err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) Update(r *request.CategoryUpsert) error {
	category, err := s.repository.Get(*r.Id)
	if err != nil {
		return err
	}
	if r.Name != nil {
		category.Name = *r.Name
	}
	if one := s.repository.FindOne(map[string]any{"name": category.Name}); one.ID != 0 && one.ID != category.ID {
		return errors.New("name exists")
	}
	if err = s.repository.Update(category); err != nil {
		return err
	}
	return nil
}

func (s *CategoryService) Delete(id uint) error {
	return s.repository.Delete(id)
}
func (s *CategoryService) Detail(id uint) (any, error) {
	category, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return category.ToMap(), nil
}
func (s *CategoryService) List(r *request.NormalSearch) (any, int64, error) {
	b := builder.New()
	if r.Keyword != nil {
		b.Like("name", *r.Keyword)
	}
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
	}
	b.Order(orderSet[r.Sort])
	categories, total, err := s.repository.ListWithCount(r.Offset, r.Limit, b)
	if err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(categories))
	for _, v := range categories {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}
