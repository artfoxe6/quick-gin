package db

import (
	"database/sql"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var (
	//配置文件中的标识
	section = "database"
	//配置信息结构
	c = struct {
		Connection string
		User       string
		Password   string
		Host       string
		DbName     string
	}{}
	//连接实例
	instance = new(sqlx.DB)
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
	var err error
	instance, err = sqlx.Open(c.Connection, c.User+":"+c.Password+"@tcp("+c.Host+")/"+
		c.DbName+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

//获取实例
func Instance() *sqlx.DB {
	if !isLoad {
		load()
	}
	return instance
}

//查询列表
func Select(dest interface{}, query string, args ...interface{}) error {
	return Instance().Select(dest, query, args...)
}

//查询单个
func Get(dest interface{}, query string, args ...interface{}) error {
	return Instance().Get(dest, query, args...)
}

//插入
func Insert(query string, args ...interface{}) (sql.Result, error) {
	return Instance().Exec(query, args...)
}

//更新
func Update(query string, args ...interface{}) (sql.Result, error) {
	return Instance().Exec(query, args...)
}
