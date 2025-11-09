package cache

import (
	"fmt"
	"sync"

	"github.com/artfoxe6/quick-gin/internal/app/core/config"
)

var (
	globalCache CacheInterface
	once        sync.Once
)

// GetCache 获取全局缓存实例（单例模式）
func GetCache() CacheInterface {
	once.Do(func() {
		cacheType := CacheType(config.Cache.Type)
		switch cacheType {
		case CacheTypeMemory:
			globalCache = NewMemoryCache()
			fmt.Printf("使用内存缓存\n")
		case CacheTypeRedis:
			globalCache = NewRedisCache()
			fmt.Printf("使用 Redis 缓存\n")
		default:
			// 默认使用内存缓存
			globalCache = NewMemoryCache()
			fmt.Printf("缓存类型 '%s' 不支持，默认使用内存缓存\n", cacheType)
		}
	})
	return globalCache
}

// NewCache 创建指定类型的缓存实例
func NewCache(cacheType CacheType) CacheInterface {
	switch cacheType {
	case CacheTypeMemory:
		return NewMemoryCache()
	case CacheTypeRedis:
		return NewRedisCache()
	default:
		return NewMemoryCache()
	}
}
