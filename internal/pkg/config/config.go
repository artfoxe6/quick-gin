package config

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	App = new(struct {
		AppMode      string
		Listen       string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		LogDir       string
	})
	Jwt = new(struct {
		Secret     string
		Exp        int
		RefreshExp int
	})
	Database = new(struct {
		Host     string
		Port     int
		User     string
		Password string
		Db       string
	})
	Redis = new(struct {
		Host      string
		Port      int
		Password  string
		Db        int
		MaxIdle   int
		MaxActive int
	})
)

// 加载配置信息
func Load(cfgPath string) {
	h, err := ini.Load(cfgPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	iniMap := map[string]interface{}{
		"app":      App,
		"jwt":      Jwt,
		"database": Database,
		"redis":    Redis,
	}
	for k, v := range iniMap {
		if err = h.Section(k).MapTo(v); err != nil {
			log.Fatal(err.Error())
		}
	}

}
