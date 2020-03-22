package api

import (
	"github.com/gin-gonic/gin"
	"quick_gin/service/ArticleService"
	"quick_gin/util/request"
)

func LoadArticleRoute(r *gin.Engine) {
	g := r.Group("/article")

	g.POST("/add", func(c *gin.Context) {
		ArticleService.Add(request.New(c))
	})
}
