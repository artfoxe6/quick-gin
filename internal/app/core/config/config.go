package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	App = new(struct {
		AppMode      string
		Listen       string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		LogDir       string
		SignKey      string
	})
	Jwt = new(struct {
		Secret     string
		Exp        int
		RefreshExp int
	})
	Database = new(struct {
		Type            string // mysql, sqlite
		Host            string
		Port            int
		User            string
		Password        string
		Db              string
		DbFile          string // SQLite database file path
		ConnMaxLifeTime int    // 连接最大生命周期（分钟）
		MaxPoolSize     int    // 最大连接池大小
		MaxIdle         int    // 最大空闲连接数
	})
	Redis = new(struct {
		Host      string
		Port      int
		Password  string
		Db        int
		MaxIdle   int
		MaxActive int
	})
	Super = new(struct {
		Email    string
		Password string
	})
	Sendgrid = new(struct {
		Key string
	})
	Oss = new(struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
		Url             string
		ImageDir        string
		CdnUrl          string
	})
	Cache = new(struct {
		Type      string // memory, redis
		Host      string // Redis host
		Port      int    // Redis port
		Password  string // Redis password
		Db        int    // Redis db
		MaxIdle   int    // Redis max idle
		MaxActive int    // Redis max active
	})
)

func Setup(cfgPath string) {
	h, err := ini.Load(cfgPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	iniMap := map[string]interface{}{
		"app":      App,
		"jwt":      Jwt,
		"database": Database,
		"redis":    Redis,
		"cache":    Cache,
		"super":    Super,
		"sendgrid": Sendgrid,
		"oss":      Oss,
	}
	for k, v := range iniMap {
		if err = h.Section(k).MapTo(v); err != nil {
			log.Fatal(err.Error())
		}
	}
	if App.LogDir != "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = ensureDir(filepath.Join(pwd, App.LogDir))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		return nil
	}
	return err
}
