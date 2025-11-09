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

	userService := services.NewUserService(userRepo, codeRepo)

	user := handlers.NewUserHandler(userService)

	api := r.Group("/api", middleware.Sign(config.App.SignKey))
	admin := r.Group("/admin", middleware.Auth(userService, "admin"))
	api.POST("/user/login", user.Login)
	api.POST("/user/fresh-token", user.FreshToken)
	api.POST("/user/register", user.Register)
	api.POST("/user/password/update", user.UpdatePassword)
	api.POST("/code", user.Code)
	api.POST("/upload", user.Upload)

	admin.POST("/user/create", user.Create)
	admin.POST("/user/update", user.Update)
	admin.POST("/user/delete", user.Delete)
	admin.GET("/user/detail", user.Detail)
	admin.GET("/user/list", user.List)

	return r
}
