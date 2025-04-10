package services

import (
	"errors"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
)

type TagService struct {
	repository *repositories.TagRepository
}

func NewTagService() *TagService {
	return &TagService{
		repository: repositories.NewTagRepository(),
	}
}

func (s *TagService) Create(r *request.TagUpsert) error {
	tag := models.Tag{
		Name: *r.Name,
	}
	if one := s.repository.FindOne(map[string]any{"name": tag.Name}); one.ID != 0 {
		return errors.New("name exists")
	}
	if err := s.repository.Create(&tag); err != nil {
		return err
	}
	return nil
}

func (s *TagService) Update(r *request.TagUpsert) error {
	tag, err := s.repository.Get(*r.Id)
	if err != nil {
		return err
	}
	if r.Name != nil {
		tag.Name = *r.Name
	}
	if one := s.repository.FindOne(map[string]any{"name": tag.Name}); one.ID != 0 && one.ID != tag.ID {
		return errors.New("name exists")
	}
	if err = s.repository.Update(tag); err != nil {
		return err
	}
	return nil
}

func (s *TagService) Delete(id uint) error {
	return s.repository.Delete(id)
}
func (s *TagService) Detail(id uint) (any, error) {
	tag, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return tag.ToMap(), nil
}
func (s *TagService) List(r *request.NormalSearch) (any, int64, error) {
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
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(tags))
	for _, v := range tags {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}
