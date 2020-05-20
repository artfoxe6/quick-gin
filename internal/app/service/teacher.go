package service

import (
	"github.com/artfoxe6/quick-gin/internal/app/model"
	"gorm.io/gorm"
)

type Teacher struct {
	Db *gorm.DB
}

func (t Teacher) Students(teacherId int) []model.Student {
	students := []model.Student{}
	t.Db.Where("teacherId= ? ", teacherId).Find(&students)
	return students
}

func (t Teacher) Info(teacherId int) model.Teacher {
	teacher := model.Teacher{}
	t.Db.Where("id= ? ", teacherId).First(&teacher)
	return teacher
}
