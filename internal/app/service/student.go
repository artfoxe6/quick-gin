package service

import (
	"github.com/artfoxe6/quick-gin/internal/app/model"
	"gorm.io/gorm"
)

type Student struct {
	Db *gorm.DB
}

func (s Student) Info(studentId int) model.Student {
	student := model.Student{}
	s.Db.Where("id= ? ", studentId).First(&student)
	return student
}
