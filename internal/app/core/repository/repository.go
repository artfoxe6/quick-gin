package repository

import (
	"github.com/artfoxe6/quick-gin/internal/app/core/repository/builder"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	db      *gorm.DB
	builder *builder.Builder
}

func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

func (r *Repository[T]) Create(m *T) error {
	return r.db.Create(m).Error
}

func (r *Repository[T]) Update(m *T) error {
	return r.db.Save(m).Error
}

func (r *Repository[T]) Clear(m *T, associations ...string) error {
	for _, association := range associations {
		err := r.db.Model(m).Association(association).Clear()
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *Repository[T]) Replace(m *T, associations map[string]any) error {
	for key, association := range associations {
		err := r.db.Model(m).Association(key).Replace(association)
		if err != nil {
			return err
		}
	}
	return nil

}
func (r *Repository[T]) Get(id uint, builders ...*builder.Builder) (*T, error) {
	m := new(T)
	tx := r.db.Model(m)
	for _, b := range builders {
		b.Exec(tx)
	}
	err := tx.First(m, id).Error
	return m, err
}

func (r *Repository[T]) Delete(id uint) error {
	return r.db.Delete(new(T), id).Error
}

func (r *Repository[T]) FindOne(conditions map[string]any, builders ...*builder.Builder) *T {
	m := new(T)
	tx := r.db.Model(m)
	for k, v := range conditions {
		tx.Where(k+"=?", v)
	}
	for _, b := range builders {
		b.Exec(tx)
	}
	tx.Order("id desc").First(m)
	return m
}

func (r *Repository[T]) FindBy(conditions map[string]any, builders ...*builder.Builder) ([]T, error) {
	var list []T
	tx := r.db.Model(new(T))
	for k, v := range conditions {
		tx.Where(k+"=?", v)
	}
	for _, b := range builders {
		b.Exec(tx)
	}
	err := tx.Limit(500).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository[T]) List(offset int, limit int, builders ...*builder.Builder) ([]T, error) {
	var list []T
	tx := r.db.Model(new(T))
	for _, b := range builders {
		b.Exec(tx)
	}
	err := tx.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (r *Repository[T]) ListWithCount(offset int, limit int, builders ...*builder.Builder) ([]T, int64, error) {
	var list []T
	var count int64
	tx := r.db.Model(new(T))
	for _, b := range builders {
		b.Exec(tx)
	}
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := tx.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *Repository[T]) ListIds(offset int, limit int, builders ...*builder.Builder) ([]uint, error) {
	var list []uint
	tx := r.db.Model(new(T))
	for _, b := range builders {
		b.Exec(tx)
	}
	err := tx.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository[T]) ListIdWithCount(offset int, limit int, builders ...*builder.Builder) ([]uint, int64, error) {
	var list []uint
	var count int64
	tx := r.db.Model(new(T))
	for _, b := range builders {
		b.Exec(tx)
	}
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := tx.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (r *Repository[T]) Count(builders ...*builder.Builder) int64 {
	var count int64
	tx := r.db.Model(new(T))
	for _, b := range builders {
		b.Exec(tx)
	}
	if err := tx.Count(&count).Error; err != nil {
		return 0
	}
	return count
}
