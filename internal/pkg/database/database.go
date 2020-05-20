package database

import (
	"fmt"
	"github.com/artfoxe6/quick-gin/internal/app/model"
	"github.com/artfoxe6/quick-gin/internal/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var db *gorm.DB

func Db() *gorm.DB {
	if db == nil {
		setup()
	}
	return db
}

func setup() {
	if config.Database == nil {
		log.Fatalln("database config not init")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Db,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalln("数据库连接失败", err)
	}

	// 迁移表
	tables := []interface{}{
		&model.Student{},
		&model.Teacher{},
	}
	err = db.AutoMigrate(tables...)
	if err != nil {
		log.Fatalln(err)
	}
}
