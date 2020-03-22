package api

import (
	"github.com/gin-gonic/gin"
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

	//登录
	g.GET("/login", func(c *gin.Context) {
		UserService.Login(request.New(c))
	})

	//获取用户信息以及用户发表的文章
	g.GET("/info_with_article", func(c *gin.Context) {
		UserService.InfoWithArticle(request.New(c))
	})
}
