package middleware

import (
	"github.com/gin-gonic/gin"
	"quick_gin/util/request"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		r := request.New(c)

		// 检查header中是否存在 Authorization 字段
		token := r.Header("Authorization", "")
		if token == "" {
			c.JSON(200, "auth failed")
			c.Abort()
		}
		//fmt.Println("middleware before")
		c.Next()
		//fmt.Println("middleware after")
	}
}
