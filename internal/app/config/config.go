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
