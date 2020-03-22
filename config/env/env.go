package env

import (
	"github.com/go-ini/ini"
	"log"
)

// 定义配置结构
type (
	serverConfig struct {
		Port            int
		DebugMode       string
		ReadTimeout     int
		WriteTimeout    int
		ShutdownTimeout int
	}
)

var (
	// 防止重复加载
	isLoad = false
	// 配置文件标志映射到struct结构
	envMap = map[string]interface{}{
		"server": new(serverConfig),
	}
)

// 加载配置信息
func load() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for key, value := range envMap {
		if err := cfg.Section(key).MapTo(value); err != nil {
			log.Fatal(err.Error())
		}
	}
	isLoad = true
}

// 暴露给外部使用
func Server() *serverConfig {
	if !isLoad {
		load()
	}
	return envMap["server"].(*serverConfig)
}
