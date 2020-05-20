package app

import (
	"github.com/artfoxe6/quick-gin/internal/app/api"
	"github.com/artfoxe6/quick-gin/internal/pkg/config"
	"github.com/artfoxe6/quick-gin/internal/pkg/middleware"
	"github.com/artfoxe6/quick-gin/internal/pkg/reqLog"
	"github.com/gin-gonic/gin"
)

func Handler() *gin.Engine {
	r := gin.New()

	// 加载中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(gin.LoggerWithWriter(reqLog.Log{Dir: config.App.LogDir}))
	//捕获自定义业务错误
	r.Use(middleware.Recovery())

	//基本的jwt认证，可指定角色，在创建jwt的时候设定的角色
	authStudent := middleware.Auth("student")
	authTeacher := middleware.Auth("teacher")

	r.GET("/", api.Index)

	studentGroup := r.Group("/student", authStudent)
	{
		studentApi := api.Student()
		studentGroup.GET("/info", studentApi.Info)
		// 模拟业务错误捕获
		studentGroup.GET("/test", studentApi.Test)
	}

	teacherGroup := r.Group("/teacher", authTeacher)
	{
		teacherApi := api.Teacher()
		teacherGroup.GET("/students", teacherApi.Students)
		teacherGroup.GET("/info", teacherApi.Info)
	}

	return r
}
