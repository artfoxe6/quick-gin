package cache

import "time"

// CacheInterface 定义缓存接口
type CacheInterface interface {
	// Set 设置缓存，ex 为过期时间（秒），0 表示永不过期
	Set(key string, data any, ex int) error

	// Get 获取缓存，返回值和是否存在标志
	Get(key string) (any, bool)

	// Delete 删除缓存
	Delete(key string) error

	// Clear 清空所有缓存
	Clear() error

	// Exists 检查缓存是否存在
	Exists(key string) bool

	// SetWithTTL 设置带 TTL 的缓存
	SetWithTTL(key string, data any, ttl time.Duration) error
}

// CacheType 缓存类型
type CacheType string

const (
	CacheTypeMemory CacheType = "memory"
	CacheTypeRedis  CacheType = "redis"
)