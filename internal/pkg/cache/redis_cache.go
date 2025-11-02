package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/gomodule/redigo/redis"
)

// RedisCache Redis 缓存实现
type RedisCache struct {
	pool *redis.Pool
}

// NewRedisCache 创建 Redis 缓存实例
func NewRedisCache() *RedisCache {
	// 先尝试获取连接池，如果失败则返回无效实例
	pool, err := getRedisPool()
	if err != nil {
		// 如果连接失败，返回一个总是失败的缓存实例
		return &RedisCache{pool: nil}
	}

	// 测试连接
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("ping"); err != nil {
		// 如果连接失败，返回一个总是失败的缓存实例
		return &RedisCache{pool: nil}
	}

	return &RedisCache{pool: pool}
}

// Set 设置缓存
func (r *RedisCache) Set(key string, data any, ex int) error {
	if r.pool == nil {
		return fmt.Errorf("Redis 不可用")
	}
	conn := r.pool.Get()
	defer conn.Close()

	temp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if ex > 0 {
		_, err = conn.Do("set", key, temp, "EX", ex)
	} else {
		_, err = conn.Do("set", key, temp)
	}

	return err
}

// SetWithTTL 设置带 TTL 的缓存
func (r *RedisCache) SetWithTTL(key string, data any, ttl time.Duration) error {
	return r.Set(key, data, int(ttl.Seconds()))
}

// Get 获取缓存
func (r *RedisCache) Get(key string) (any, bool) {
	if r.pool == nil {
		return nil, false
	}
	conn := r.pool.Get()
	defer conn.Close()

	ext, err := redis.Bool(conn.Do("exists", key))
	if err != nil || !ext {
		return nil, false
	}

	d, err := redis.Bytes(conn.Do("get", key))
	if err != nil {
		return nil, false
	}

	var dst any
	err = json.Unmarshal(d, &dst)
	if err != nil {
		return nil, false
	}

	return dst, true
}

// Delete 删除缓存
func (r *RedisCache) Delete(key string) error {
	if r.pool == nil {
		return fmt.Errorf("Redis 不可用")
	}
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("del", key)
	return err
}

// Clear 清空所有缓存
func (r *RedisCache) Clear() error {
	if r.pool == nil {
		return fmt.Errorf("Redis 不可用")
	}
	conn := r.pool.Get()
	defer conn.Close()

	// 清空当前数据库
	_, err := conn.Do("flushdb")
	return err
}

// Exists 检查缓存是否存在
func (r *RedisCache) Exists(key string) bool {
	if r.pool == nil {
		return false
	}
	conn := r.pool.Get()
	defer conn.Close()

	ext, err := redis.Bool(conn.Do("exists", key))
	return err == nil && ext
}

// getRedisPool 获取 Redis 连接池（复用原有的连接池逻辑）
func getRedisPool() (*redis.Pool, error) {
	c := config.Redis

	// 检查配置是否完整
	if c.Host == "" {
		return nil, fmt.Errorf("Redis 主机地址未配置")
	}
	if c.Port <= 0 {
		return nil, fmt.Errorf("Redis 端口未配置")
	}

	address := fmt.Sprintf("%s:%d", c.Host, c.Port)

	pool := &redis.Pool{
		Dial: func() (con redis.Conn, err error) {
			con, err = redis.Dial("tcp", address,
				redis.DialPassword(c.Password),
				redis.DialDatabase(c.Db))
			if err != nil {
				return nil, err
			}
			return con, err
		},
		MaxIdle:         c.MaxIdle,
		MaxActive:       c.MaxActive,
		Wait:            true,
		MaxConnLifetime: 0,
	}

	return pool, nil
}