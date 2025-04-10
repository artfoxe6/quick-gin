package services

import (
	"errors"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"gorm.io/gorm"
	"strings"
	"time"
)

type NewsService struct {
	repository *repositories.NewsRepository
}

func NewNewsService() *NewsService {
	return &NewsService{
		repository: repositories.NewNewsRepository(),
	}
}

func (s *NewsService) Create(r *request.NewsUpsert, user *models.User) (uint, error) {
	news := models.News{}
	if r.Slug != nil && *r.Slug != "" {
		news.Slug = *r.Slug
	} else {
		news.Slug = news.GetSlug()
	}
	if r.Title != nil {
		news.Title = *r.Title
	}
	if r.Content != nil {
		news.Content = *r.Content
	}
	if r.AuthorId != nil {
		news.AuthorId = r.AuthorId
	}
	if r.Source != nil {
		news.Source = *r.Source
	}
	if r.Status != nil {
		if news.Status != models.NEWS_STATUS_PUBLISHED && *r.Status == models.NEWS_STATUS_PUBLISHED {
			news.PublishTime = time.Now().Unix()
		}
		news.Status = *r.Status
	}
	if r.IsFeatured != nil {
		if *r.IsFeatured {
			news.IsFeatured = time.Now().Unix()
		} else {
			news.IsFeatured = 0
		}
	}
	if r.Summary != nil {
		news.Summary = *r.Summary
	}
	if r.ImageUrl != nil {
		news.ImageUrl = *r.ImageUrl
	}
	if r.VideoUrl != nil {
		news.VideoUrl = *r.VideoUrl
	}
	if r.AudioUrl != nil {
		news.AudioUrl = *r.AudioUrl
	}
	if r.TypeId != nil {
		news.TypeId = *r.TypeId
	}
	if r.Slug != nil {
		news.Slug = *r.Slug
	}

	if r.CategoryIds != nil {
		categoryRepository := repositories.NewCategoryRepository()
		tempNews, err := categoryRepository.List(0, 100, builder.New().In("id", *r.CategoryIds))
		if err != nil {
			return 0, err
		}
		for _, tempNew := range tempNews {
			news.Categories = append(news.Categories, &tempNew)
		}
	}
	if r.TagIds != nil {
		tagRepository := repositories.NewTagRepository()
		tempTags, err := tagRepository.List(0, 100, builder.New().In("id", *r.TagIds))
		if err != nil {
			return 0, err
		}
		for _, tempTag := range tempTags {
			news.Tags = append(news.Tags, &tempTag)
		}
	}
	if one := s.repository.FindOne(map[string]any{"slug": news.Slug}); one.ID != 0 {
		return 0, errors.New("slug exists")
	}
	if err := s.repository.Create(&news); err != nil {
		return 0, err
	}
	return news.ID, nil
}

func (s *NewsService) Update(r *request.NewsUpsert) error {
	news, err := s.repository.Get(*r.Id)
	if err != nil {
		return err
	}
	if r.Title != nil {
		news.Title = *r.Title
	}
	if r.Content != nil {
		news.Content = *r.Content
	}
	if r.AuthorId != nil {
		news.AuthorId = r.AuthorId
	}
	if r.Source != nil {
		news.Source = *r.Source
	}
	if r.Status != nil {
		if news.Status != models.NEWS_STATUS_PUBLISHED && *r.Status == models.NEWS_STATUS_PUBLISHED {
			news.PublishTime = time.Now().Unix()
		}
		news.Status = *r.Status
	}
	if r.IsFeatured != nil {
		if *r.IsFeatured {
			news.IsFeatured = time.Now().Unix()
		} else {
			news.IsFeatured = 0
		}
	}
	if r.Summary != nil {
		news.Summary = *r.Summary
	}
	if r.ImageUrl != nil {
		news.ImageUrl = *r.ImageUrl
	}
	if r.VideoUrl != nil {
		news.VideoUrl = *r.VideoUrl
	}
	if r.AudioUrl != nil {
		news.AudioUrl = *r.AudioUrl
	}
	if r.TypeId != nil {
		news.TypeId = *r.TypeId
	}
	if r.Slug != nil {
		news.Slug = *r.Slug
	}

	if r.CategoryIds != nil {
		categoryRepository := repositories.NewCategoryRepository()
		tempNews, err := categoryRepository.List(0, 100, builder.New().In("id", *r.CategoryIds))
		if err != nil {
			return err
		}
		for _, tempNew := range tempNews {
			news.Categories = append(news.Categories, &tempNew)
		}
	}
	if r.TagIds != nil {
		tagRepository := repositories.NewTagRepository()
		tempTags, err := tagRepository.List(0, 100, builder.New().In("id", *r.TagIds))
		if err != nil {
			return err
		}
		for _, tempTag := range tempTags {
			news.Tags = append(news.Tags, &tempTag)
		}
	}
	if one := s.repository.FindOne(map[string]any{"slug": news.Slug}); one.ID != 0 && one.ID != news.ID {
		return errors.New("slug exists")
	}
	if err = s.repository.Update(news); err != nil {
		return err
	}
	return nil
}

func (s *NewsService) Delete(id uint) error {
	return s.repository.Delete(id)
}
func (s *NewsService) Detail(id uint) (any, error) {
	b := builder.New()
	b.Preload("Categories").Preload("Author").Preload("Tags").Preload("User")
	news, err := s.repository.Get(id, b)
	if err != nil {
		return nil, err
	}
	return news.ToMap(), nil
}
func (s *NewsService) List(r *request.NewsSearch) (any, int64, error) {
	b := builder.New()
	if r.CategoryIds != nil && *r.CategoryIds != "" {
		b.Append(func(tx *gorm.DB) {
			tx.Joins("left join news_category ON news_category.news_id = news.id").Where("news_category.category_id in ?", strings.Split(*r.CategoryIds, ","))
		})
	}
	if r.TagIds != nil && *r.TagIds != "" {
		b.Append(func(tx *gorm.DB) {
			tx.Joins("left join news_tag ON news_tag.news_id = news.id").Where("news_tag.tag_id in ?", strings.Split(*r.TagIds, ","))
		})
	}
	if r.TypeIds != nil && *r.TypeIds != "" {
		b.In("type_id", strings.Split(*r.TypeIds, ","))
	}
	if r.Keyword != nil {
		b.Like("title", *r.Keyword)
	}
	if r.Status != nil {
		b.Eq("status", *r.Status)
	}
	total := s.repository.Count(b)
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
		2: "max(news.is_featured) desc",
	}
	b.Order(orderSet[r.Sort])
	b.Append(func(tx *gorm.DB) {
		tx.Select("id").Group("id")
	})
	newIds, err := s.repository.ListIds(r.Offset, r.Limit, b)
	if len(newIds) == 0 {
		return nil, 0, err
	}
	b2 := builder.New().In("id", newIds)
	b2.Preload("Categories").Preload("Author").Preload("Tags")
	news, err := s.repository.List(0, r.Limit, b2)
	if err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(news))
	for _, id := range newIds {
		for _, v := range news {
			if v.ID == id {
				list = append(list, v.ToMap())
			}
		}
	}
	return list, total, nil
}
