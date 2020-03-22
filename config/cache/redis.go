package cache

import (
	"encoding/json"
	"github.com/go-ini/ini"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var (
	//配置文件中的标识
	section = "redis"
	//配置信息结构
	c = struct {
		Host        string
		Password    string
		MaxIdle     int
		MaxActive   int
		IdleTimeout time.Duration
		Db          int
		Timeout     int
	}{}
	//连接实例
	instance = new(redis.Pool)
	//防止重复加载
	isLoad = false
)

// 加载配置信息
func load() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = cfg.Section(section).MapTo(&c)
	if err != nil {
		log.Fatalf("%v", err)
	}
	connection()
	isLoad = true
}

// 建立连接
func connection() {
	instance = &redis.Pool{
		Dial: func() (con redis.Conn, err error) {
			con, err = redis.Dial("tcp", c.Host,
				redis.DialPassword(c.Password),
				redis.DialDatabase(c.Db),
				redis.DialConnectTimeout(time.Second*time.Duration(c.Timeout)),
				redis.DialReadTimeout(time.Second*time.Duration(c.Timeout)),
				redis.DialWriteTimeout(time.Second*time.Duration(c.Timeout)))
			if err != nil {
				log.Fatalln(err.Error())
			}
			return con, err
		},
		MaxIdle:         c.MaxIdle,
		MaxActive:       c.MaxActive,
		IdleTimeout:     c.IdleTimeout,
		Wait:            true,
		MaxConnLifetime: 0,
	}
}

//获取实例
func Instance() redis.Conn {
	if !isLoad {
		load()
	}
	return instance.Get()
}

//缓存字节数组到缓存
func Set(key string, data interface{}, ex int) {
	temp, _ := json.Marshal(data)
	_, err := Instance().Do("set", key, temp, "EX", ex)
	if err != nil {
		log.Printf("%v", err)
	}
}

//从缓存中取数据
func Get(key string) (interface{}, bool) {
	var dst interface{}
	ext, _ := redis.Bool(Instance().Do("exists", key))
	if ext {
		d, _ := redis.Bytes(Instance().Do("get", key))
		err := json.Unmarshal(d, &dst)
		if err == nil {
			return dst, true
		} else {
			log.Fatalf("%v", err)
		}
	}
	return nil, false
}
