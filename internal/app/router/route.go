package router

import (
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/artfoxe6/quick-gin/internal/app/handlers"
	"github.com/artfoxe6/quick-gin/internal/app/middleware"
	"github.com/artfoxe6/quick-gin/internal/app/repositories"
	"github.com/artfoxe6/quick-gin/internal/app/services"
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

	userRepo := repositories.NewUserRepository()
	codeRepo := repositories.NewCodeRepository()
	categoryRepo := repositories.NewCategoryRepository()
	tagRepo := repositories.NewTagRepository()
	authorRepo := repositories.NewAuthorRepository()
	newsRepo := repositories.NewNewsRepository()

	userService := services.NewUserService(userRepo, codeRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	tagService := services.NewTagService(tagRepo)
	authorService := services.NewAuthorService(authorRepo)
	newsService := services.NewNewsService(newsRepo, categoryRepo, tagRepo)

	user := handlers.NewUserHandler(userService)
	category := handlers.NewCategoryHandler(categoryService)
	tag := handlers.NewTagHandler(tagService)
	author := handlers.NewAuthorHandler(authorService)
	news := handlers.NewNewsHandler(newsService)

	api := r.Group("/api", middleware.Sign(config.App.SignKey))
	admin := r.Group("/admin", middleware.Auth(userService, "admin"))
	api.POST("/user/login", user.Login)
	api.POST("/user/fresh-token", user.FreshToken)
	api.POST("/user/register", user.Register)
	api.POST("/user/password/update", user.UpdatePassword)
	api.POST("/code", user.Code)
	api.POST("/upload", user.Upload)

	api.GET("/news/detail", news.Detail)
	api.GET("/news/list", news.List)

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

	admin.POST("/news/tag/create", tag.Create)
	admin.POST("/news/tag/update", tag.Update)
	admin.POST("/news/tag/delete", tag.Delete)
	admin.GET("/news/tag/detail", tag.Detail)
	admin.GET("/news/tag/list", tag.List)

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
