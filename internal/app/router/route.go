package router

import (
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/artfoxe6/quick-gin/internal/app/handlers"
	"github.com/artfoxe6/quick-gin/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func Handler() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(gin.LoggerWithWriter(middleware.Log{Dir: config.App.LogDir}))
	r.Use(middleware.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api", middleware.Sign(config.App.SignKey))
	admin := r.Group("/admin", middleware.Auth("admin"))
	user := handlers.NewUserHandler()
	api.POST("/user/login", user.Login)
	api.POST("/user/fresh-token", user.FreshToken)
	api.POST("/user/register", user.Register)
	api.POST("/user/password/update", user.UpdatePassword)
	api.POST("/code", user.Code)
	api.POST("/upload", user.Upload)

	news := handlers.NewNewsHandler()
	api.GET("/news/detail", news.Detail)
	api.GET("/news/list", news.List)

	category := handlers.NewCategoryHandler()
	api.GET("/news/category/list", category.List)

	admin.POST("/news/create", news.Create)
	admin.POST("/news/update", news.Update)
	admin.POST("/news/delete", news.Delete)
	admin.GET("/news/detail", news.Detail)
	admin.GET("/news/list", news.List)

	admin.POST("/news/category/create", category.Create)
	admin.POST("/news/category/update", category.Update)
	admin.POST("/news/category/delete", category.Delete)
	admin.GET("/news/category/detail", category.Detail)
	admin.GET("/news/category/list", category.List)

	tag := handlers.NewTagHandler()
	admin.POST("/news/tag/create", tag.Create)
	admin.POST("/news/tag/update", tag.Update)
	admin.POST("/news/tag/delete", tag.Delete)
	admin.GET("/news/tag/detail", tag.Detail)
	admin.GET("/news/tag/list", tag.List)

	author := handlers.NewAuthorHandler()
	admin.POST("/news/author/create", author.Create)
	admin.POST("/news/author/update", author.Update)
	admin.POST("/news/author/delete", author.Delete)
	admin.GET("/news/author/detail", author.Detail)
	admin.GET("/news/author/list", author.List)

	admin.POST("/user/create", user.Create)
	admin.POST("/user/update", user.Update)
	admin.POST("/user/delete", user.Delete)
	admin.GET("/user/detail", user.Detail)
	admin.GET("/user/list", user.List)

	return r
}
