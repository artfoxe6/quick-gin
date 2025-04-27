package builder

import (
	"fmt"
	"gorm.io/gorm"
)

type BuilderFunc func(tx *gorm.DB)

type Builder struct {
	builders []BuilderFunc
}

func New() *Builder {
	return &Builder{
		builders: make([]BuilderFunc, 0, 10),
	}
}
func (b Builder) Exec(tx *gorm.DB) *gorm.DB {
	for _, builder := range b.builders {
		builder(tx)
	}
	return tx
}

func (b *Builder) set(key string, compare string, value any) *Builder {
	b.builders = append(b.builders, func(tx *gorm.DB) {
		tx.Where(fmt.Sprintf("%s %s ?", key, compare), value)
	})
	return b
}
func (b *Builder) Eq(key string, value any) *Builder {
	return b.set(key, "=", value)
}
func (b *Builder) Neq(key string, value any) *Builder {
	return b.set(key, "!=", value)
}
func (b *Builder) In(key string, value any) *Builder {
	return b.set(key, "in", value)
}
func (b *Builder) Nin(key string, value any) *Builder {
	return b.set(key, "not in", value)
}
func (b *Builder) Gt(key string, value any) *Builder {
	return b.set(key, ">", value)
}
func (b *Builder) Gte(key string, value any) *Builder {
	return b.set(key, ">=", value)
}
func (b *Builder) Lt(key string, value any) *Builder {
	return b.set(key, "<", value)
}
func (b *Builder) Lte(key string, value any) *Builder {
	return b.set(key, "<=", value)
}
func (b *Builder) Like(key string, value any) *Builder {
	return b.set(key, "like", "%"+value.(string)+"%")
}
func (b *Builder) Order(order any) *Builder {
	b.builders = append(b.builders, func(tx *gorm.DB) {
		tx.Order(order)
	})
	return b
}
func (b *Builder) Where(query interface{}, args ...interface{}) *Builder {
	b.builders = append(b.builders, func(tx *gorm.DB) {
		tx.Where(query, args...)
	})
	return b
}
func (b *Builder) Preload(query string, args ...interface{}) *Builder {
	b.builders = append(b.builders, func(tx *gorm.DB) {
		tx.Preload(query, args...)
	})
	return b
}
func (b *Builder) Append(builderFunc BuilderFunc) *Builder {
	b.builders = append(b.builders, builderFunc)
	return b
}
