package cache

import (
	"encoding/json"
	"github.com/artfoxe6/quick-gin/internal/pkg/config"
	"github.com/gomodule/redigo/redis"
	"log"
)

var pool *redis.Pool

func Conn() redis.Conn {
	return Pool().Get()
}

func Pool() *redis.Pool {
	if pool == nil {
		setup()
	}
	return pool
}

// 建立连接
func setup() {
	c := config.Redis
	pool = &redis.Pool{
		Dial: func() (con redis.Conn, err error) {
			con, err = redis.Dial("tcp", c.Host,
				redis.DialPassword(c.Password),
				redis.DialDatabase(c.Db))
			if err != nil {
				log.Fatalln("Redis连接错误", err.Error())
			}
			return con, err
		},
		MaxIdle:         c.MaxIdle,
		MaxActive:       c.MaxActive,
		Wait:            true,
		MaxConnLifetime: 0,
	}
}

// Set 缓存字节数组到缓存
func Set(key string, data interface{}, ex int) {
	temp, _ := json.Marshal(data)
	_, err := Conn().Do("set", key, temp, "EX", ex)
	if err != nil {
		log.Println(err)
	}
}

// Get 从缓存中取数据
func Get(key string) (interface{}, bool) {
	var dst interface{}
	ext, _ := redis.Bool(Conn().Do("exists", key))
	if ext {
		d, _ := redis.Bytes(Conn().Do("get", key))
		err := json.Unmarshal(d, &dst)
		if err == nil {
			return dst, true
		} else {
			log.Println(err)
			return nil, false
		}
	}
	return nil, false
}
