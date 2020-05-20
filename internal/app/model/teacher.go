package model

type Teacher struct {
	Base
	Name     string `gorm:"size:20"`
	Students []Student
}
