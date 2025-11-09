package db

import (
	"fmt"
	"github.com/artfoxe6/quick-gin/internal/app/core/config"
	"github.com/artfoxe6/quick-gin/internal/app/user/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

var db *gorm.DB

func Db() *gorm.DB {
	if db == nil {
		setup()
	}
	return db
}

func setup() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if config.Database == nil {
		panic("database config not init")
	}

	// 设置数据库类型，默认为 mysql
	dbType := config.Database.Type
	if dbType == "" {
		dbType = "mysql"
	}

	var dialector gorm.Dialector
	switch dbType {
	case "sqlite":
		// 确保数据库文件目录存在
		dbFile := config.Database.DbFile
		if dbFile == "" {
			dbFile = "data/app.db"
		}

		// 创建数据库文件目录
		if err := os.MkdirAll(filepath.Dir(dbFile), 0755); err != nil {
			panic(fmt.Sprintf("failed to create database directory: %v", err))
		}

		dialector = sqlite.Open(dbFile)
		fmt.Printf("Connecting to SQLite database: %s\n", dbFile)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Db,
		)
		dialector = mysql.Open(dsn)
		fmt.Printf("Connecting to MySQL database: %s@%s:%d/%s\n",
			config.Database.User, config.Database.Host, config.Database.Port, config.Database.Db)
	default:
		panic(fmt.Sprintf("unsupported database type: %s (supported: mysql, sqlite)", dbType))
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)

	var err error
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get underlying sql.DB: %v", err))
	}

	if dbType == "sqlite" {
		// SQLite 建议使用单连接以避免锁定问题
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	} else {
		// MySQL 连接池配置
		sqlDB.SetMaxOpenConns(config.Database.MaxPoolSize)
		sqlDB.SetMaxIdleConns(config.Database.MaxIdle)
		sqlDB.SetConnMaxLifetime(time.Duration(config.Database.ConnMaxLifeTime) * time.Minute)
	}

	// 自动迁移数据库表
	tables := []any{
		&model.User{},
		&model.Code{},
	}

	fmt.Println("Running database migrations...")
	err = db.AutoMigrate(tables...)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	fmt.Println("Database connection established and migrations completed successfully!")
}
