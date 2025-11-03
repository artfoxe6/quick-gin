package services

import (
	"strings"
	"time"

	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"gorm.io/gorm"
)

type NewsRepository interface {
	FindOne(map[string]any, ...*builder.Builder) *models.News
	Create(*models.News) error
	Update(*models.News) error
	Delete(uint) error
	Get(uint, ...*builder.Builder) (*models.News, error)
	List(int, int, ...*builder.Builder) ([]models.News, error)
	ListIds(int, int, ...*builder.Builder) ([]uint, error)
	Count(...*builder.Builder) int64
	Replace(*models.News, map[string]any) error
}

type NewsCategoryRepository interface {
	List(int, int, ...*builder.Builder) ([]models.Category, error)
}

type NewsTagRepository interface {
	List(int, int, ...*builder.Builder) ([]models.Tag, error)
}

type NewsService interface {
	Create(r *request.NewsUpsert, user *models.User) (uint, error)
	Update(r *request.NewsUpsert) error
	Delete(id uint) error
	Detail(id uint) (any, error)
	List(r *request.NewsSearch) (any, int64, error)
}

type newsService struct {
	newsRepo     NewsRepository
	categoryRepo NewsCategoryRepository
	tagRepo      NewsTagRepository
	now          func() time.Time
}

func NewNewsService(newsRepo NewsRepository, categoryRepo NewsCategoryRepository, tagRepo NewsTagRepository) NewsService {
	return &newsService{
		newsRepo:     newsRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		now:          time.Now,
	}
}

func (s *newsService) Create(r *request.NewsUpsert, user *models.User) (uint, error) {
	assembler := newNewsAssembler(s.now)
	news := assembler.buildForCreate(r, user)

	if err := s.attachAssociations(&news, r); err != nil {
		return 0, err
	}
	if existing := s.newsRepo.FindOne(map[string]any{"slug": news.Slug}); existing.ID != 0 {
		return 0, apperr.BadRequest("slug exists")
	}
	if err := s.newsRepo.Create(&news); err != nil {
		return 0, apperr.Internal(err)
	}
	return news.ID, nil
}

func (s *newsService) Update(r *request.NewsUpsert) error {
	if r.Id == nil {
		return apperr.BadRequest("id is required")
	}
	news, err := s.newsRepo.Get(*r.Id)
	if err != nil {
		return apperr.Internal(err)
	}

	assembler := newNewsAssembler(s.now)
	assembler.apply(r, news)

	if err := s.attachAssociations(news, r); err != nil {
		return err
	}

	if one := s.newsRepo.FindOne(map[string]any{"slug": news.Slug}); one.ID != 0 && one.ID != news.ID {
		return apperr.BadRequest("slug exists")
	}
	if err = s.newsRepo.Update(news); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *newsService) Delete(id uint) error {
	if err := s.newsRepo.Delete(id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *newsService) Detail(id uint) (any, error) {
	b := builder.New()
	b.Preload("Categories").Preload("Author").Preload("Tags").Preload("User")
	news, err := s.newsRepo.Get(id, b)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return news.ToMap(), nil
}

func (s *newsService) List(r *request.NewsSearch) (any, int64, error) {
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
	total := s.newsRepo.Count(b)
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
		2: "max(news.is_featured) desc",
	}
	b.Order(orderSet[r.Sort])
	b.Append(func(tx *gorm.DB) {
		tx.Select("id").Group("id")
	})
	newIds, err := s.newsRepo.ListIds(r.Offset, r.Limit, b)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	if len(newIds) == 0 {
		return []map[string]any{}, total, nil
	}
	b2 := builder.New().In("id", newIds)
	b2.Preload("Categories").Preload("Author").Preload("Tags")
	news, err := s.newsRepo.List(0, r.Limit, b2)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	list := make([]map[string]any, 0, len(news))
	for _, id := range newIds {
		for i := range news {
			if news[i].ID == id {
				list = append(list, news[i].ToMap())
				break
			}
		}
	}
	return list, total, nil
}

func (s *newsService) attachAssociations(news *models.News, r *request.NewsUpsert) error {
	associations := make(map[string]any)
	if r.CategoryIds != nil {
		ids := *r.CategoryIds
		if len(ids) == 0 {
			associations["Categories"] = []*models.Category{}
		} else {
			categories, err := s.categoryRepo.List(0, len(ids), builder.New().In("id", ids))
			if err != nil {
				return apperr.Internal(err)
			}
			cats := make([]*models.Category, 0, len(categories))
			for i := range categories {
				cats = append(cats, &categories[i])
			}
			associations["Categories"] = cats
		}
	}
	if r.TagIds != nil {
		ids := *r.TagIds
		if len(ids) == 0 {
			associations["Tags"] = []*models.Tag{}
		} else {
			tags, err := s.tagRepo.List(0, len(ids), builder.New().In("id", ids))
			if err != nil {
				return apperr.Internal(err)
			}
			ts := make([]*models.Tag, 0, len(tags))
			for i := range tags {
				ts = append(ts, &tags[i])
			}
			associations["Tags"] = ts
		}
	}
	if len(associations) == 0 {
		return nil
	}
	// Replace associations only when updating existing record to avoid gorm creating duplicates during create.
	if news.ID == 0 {
		if cats, ok := associations["Categories"].([]*models.Category); ok {
			news.Categories = cats
		}
		if ts, ok := associations["Tags"].([]*models.Tag); ok {
			news.Tags = ts
		}
		return nil
	}
	if err := s.newsRepo.Replace(news, associations); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

// newNewsAssembler constructs a reusable assembler for news entities.
func newNewsAssembler(now func() time.Time) newsAssembler {
	if now == nil {
		now = time.Now
	}
	return newsAssembler{now: now}
}

type newsAssembler struct {
	now func() time.Time
}

func (a newsAssembler) buildForCreate(r *request.NewsUpsert, user *models.User) models.News {
	news := models.News{}
	a.apply(r, &news)
	if user != nil {
		userID := user.ID
		news.UserId = &userID
	}
	if news.Slug == "" {
		news.Slug = news.GetSlug()
	}
	return news
}

func (a newsAssembler) apply(r *request.NewsUpsert, news *models.News) {
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
			news.PublishTime = a.now().Unix()
		}
		news.Status = *r.Status
	}
	if r.IsFeatured != nil {
		if *r.IsFeatured {
			news.IsFeatured = a.now().Unix()
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
	if r.Slug != nil && *r.Slug != "" {
		news.Slug = *r.Slug
	}
}
