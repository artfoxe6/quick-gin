package services

import (
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
)

type TagRepository interface {
	FindOne(map[string]any, ...*builder.Builder) *models.Tag
	Create(*models.Tag) error
	Update(*models.Tag) error
	Delete(uint) error
	Get(uint, ...*builder.Builder) (*models.Tag, error)
	ListWithCount(int, int, ...*builder.Builder) ([]models.Tag, int64, error)
}

type TagService interface {
	Create(r *request.TagUpsert) error
	Update(r *request.TagUpsert) error
	Delete(id uint) error
	Detail(id uint) (any, error)
	List(r *request.NormalSearch) (any, int64, error)
}

type tagService struct {
	repository TagRepository
}

func NewTagService(repository TagRepository) TagService {
	return &tagService{repository: repository}
}

func (s *tagService) Create(r *request.TagUpsert) error {
	tag := models.Tag{
		Name: *r.Name,
	}
	if one := s.repository.FindOne(map[string]any{"name": tag.Name}); one.ID != 0 {
		return apperr.BadRequest("name exists")
	}
	if err := s.repository.Create(&tag); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *tagService) Update(r *request.TagUpsert) error {
	tag, err := s.repository.Get(*r.Id)
	if err != nil {
		return apperr.Internal(err)
	}
	if r.Name != nil {
		tag.Name = *r.Name
	}
	if one := s.repository.FindOne(map[string]any{"name": tag.Name}); one.ID != 0 && one.ID != tag.ID {
		return apperr.BadRequest("name exists")
	}
	if err = s.repository.Update(tag); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *tagService) Delete(id uint) error {
	if err := s.repository.Delete(id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *tagService) Detail(id uint) (any, error) {
	tag, err := s.repository.Get(id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return tag.ToMap(), nil
}

func (s *tagService) List(r *request.NormalSearch) (any, int64, error) {
	b := builder.New()
	if r.Keyword != nil {
		b.Like("name", *r.Keyword)
	}
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
	}
	b.Order(orderSet[r.Sort])
	tags, total, err := s.repository.ListWithCount(r.Offset, r.Limit, b)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	list := make([]map[string]any, 0, len(tags))
	for _, v := range tags {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}
