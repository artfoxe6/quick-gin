package db

import (
	"fmt"
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/artfoxe6/quick-gin/internal/app/models"
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if config.Database == nil {
		panic("database config not init")
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
		panic(err)
	}

	tables := []any{
		&models.Author{},
		&models.Category{},
		&models.News{},
		&models.Tag{},
		&models.User{},
	}
	err = db.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}

}
