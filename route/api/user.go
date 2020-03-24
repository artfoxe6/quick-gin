package api

import (
	"github.com/gin-gonic/gin"
	"quick_gin/middleware"
	"quick_gin/service/UserService"
	"quick_gin/util/request"
)

// 此函数需要在 route.go 调用
func LoadUserRoute(r *gin.Engine) {
	g := r.Group("/user")

	//添加用户
	g.POST("/add", func(context *gin.Context) {

		UserService.Add(request.New(context))
	})

	//用户列表
	g.GET("/list", func(context *gin.Context) {
		UserService.List(request.New(context))
	})

	//用户列表以及发表的文章
	g.GET("/list_with_article", func(c *gin.Context) {
		UserService.ListWithArticles(request.New(c))
	})

	//获取token
	g.GET("/token", func(context *gin.Context) {

		UserService.CreateToken(request.New(context))
	})
	//需要token认证才能请求的路由
	auth := r.Group("/user/info")
	auth.Use(middleware.Auth())
	{
		//用户列表
		auth.GET("/", func(context *gin.Context) {
			UserService.Info(request.New(context))
		})
	}

}
