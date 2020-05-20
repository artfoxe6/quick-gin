package model

type Student struct {
	Base
	Name      string `gorm:"size:20"`
	Score     int
	TeacherId int
	Teacher   Teacher
}
