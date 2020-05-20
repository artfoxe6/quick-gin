package api

import (
	"github.com/artfoxe6/quick-gin/internal/app/service"
	"github.com/artfoxe6/quick-gin/internal/pkg/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeacherApi struct {
	service service.Teacher
}

func Teacher() TeacherApi {
	return TeacherApi{service: service.Teacher{
		Db: database.Db(),
	}}
}

func (t TeacherApi) Students(c *gin.Context) {
	uid := c.MustGet("uid")
	students := t.service.Students(uid.(int))
	c.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}

func (t TeacherApi) Info(c *gin.Context) {
	uid := c.MustGet("uid")
	info := t.service.Info(uid.(int))
	c.JSON(http.StatusOK, gin.H{
		"info": info,
	})
}
