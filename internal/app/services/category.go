package services

import (
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
)

type CategoryRepository interface {
	FindOne(map[string]any, ...*builder.Builder) *models.Category
	Create(*models.Category) error
	Update(*models.Category) error
	Delete(uint) error
	Get(uint, ...*builder.Builder) (*models.Category, error)
	ListWithCount(int, int, ...*builder.Builder) ([]models.Category, int64, error)
}

type CategoryService interface {
	Create(r *request.CategoryUpsert) error
	Update(r *request.CategoryUpsert) error
	Delete(id uint) error
	Detail(id uint) (any, error)
	List(r *request.NormalSearch) (any, int64, error)
}

type categoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) CategoryService {
	return &categoryService{repository: repository}
}

func (s *categoryService) Create(r *request.CategoryUpsert) error {
	category := models.Category{
		Name: *r.Name,
	}
	if one := s.repository.FindOne(map[string]any{"name": category.Name}); one.ID != 0 {
		return apperr.BadRequest("name exists")
	}
	if err := s.repository.Create(&category); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *categoryService) Update(r *request.CategoryUpsert) error {
	category, err := s.repository.Get(*r.Id)
	if err != nil {
		return apperr.Internal(err)
	}
	if r.Name != nil {
		category.Name = *r.Name
	}
	if one := s.repository.FindOne(map[string]any{"name": category.Name}); one.ID != 0 && one.ID != category.ID {
		return apperr.BadRequest("name exists")
	}
	if err = s.repository.Update(category); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *categoryService) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *categoryService) Detail(id uint) (any, error) {
	category, err := s.repository.Get(id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return category.ToMap(), nil
}

func (s *categoryService) List(r *request.NormalSearch) (any, int64, error) {
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
		return nil, 0, apperr.Internal(err)
	}
	list := make([]map[string]any, 0, len(categories))
	for _, v := range categories {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}
