package cache

import (
	"encoding/json"
	"github.com/artfoxe6/quick-gin/internal/app/config"
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

func setup() {
	c := config.Redis
	pool = &redis.Pool{
		Dial: func() (con redis.Conn, err error) {
			con, err = redis.Dial("tcp", c.Host,
				redis.DialPassword(c.Password),
				redis.DialDatabase(c.Db))
			if err != nil {
				log.Fatalln(err.Error())
			}
			return con, err
		},
		MaxIdle:         c.MaxIdle,
		MaxActive:       c.MaxActive,
		Wait:            true,
		MaxConnLifetime: 0,
	}
}

func Set(key string, data interface{}, ex int) {
	temp, _ := json.Marshal(data)
	_, err := Conn().Do("set", key, temp, "EX", ex)
	if err != nil {
		log.Println(err)
	}
}

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
