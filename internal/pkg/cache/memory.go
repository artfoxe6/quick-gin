package cache

import (
	"encoding/json"
	"sync"
	"time"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	store map[string]*cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	value      any
	expiration time.Time
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		store: make(map[string]*cacheItem),
	}

	// 启动清理过期缓存的协程
	go cache.startCleanup()

	return cache
}

// Set 设置缓存
func (m *MemoryCache) Set(key string, data any, ex int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var expiration time.Time
	if ex > 0 {
		expiration = time.Now().Add(time.Duration(ex) * time.Second)
	}

	m.store[key] = &cacheItem{
		value:      data,
		expiration: expiration,
	}

	return nil
}

// SetWithTTL 设置带 TTL 的缓存
func (m *MemoryCache) SetWithTTL(key string, data any, ttl time.Duration) error {
	return m.Set(key, data, int(ttl.Seconds()))
}

// Get 获取缓存
func (m *MemoryCache) Get(key string) (any, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.store[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Delete 删除缓存
func (m *MemoryCache) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.store, key)
	return nil
}

// Clear 清空所有缓存
func (m *MemoryCache) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store = make(map[string]*cacheItem)
	return nil
}

// Exists 检查缓存是否存在
func (m *MemoryCache) Exists(key string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.store[key]
	if !exists {
		return false
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		return false
	}

	return true
}

// startCleanup 定期清理过期缓存
func (m *MemoryCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanup()
	}
}

// cleanup 清理过期缓存
func (m *MemoryCache) cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for key, item := range m.store {
		if !item.expiration.IsZero() && now.After(item.expiration) {
			delete(m.store, key)
		}
	}
}

// MarshalJSON 实现 json.Marshaler 接口
func (m *MemoryCache) MarshalJSON() ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return json.Marshal(m.store)
}
