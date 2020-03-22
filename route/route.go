package route

import (
	"github.com/gin-gonic/gin"
	"quick_gin/config/env"
	"quick_gin/route/api"
)

var Route *gin.Engine

func Init() {
	//设置调试模式
	gin.SetMode(env.Server().DebugMode)

	//新建一个空路由
	Route = gin.New()

	//使用日志中间件
	if env.Server().DebugMode == "debug" {
		Route.Use(gin.Logger())
	}
	//其他路由
	api.LoadUserRoute(Route)
	api.LoadArticleRoute(Route)

}
