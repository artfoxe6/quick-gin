package middleware

import (
	"github.com/artfoxe6/quick-gin/internal/app/api"
	"github.com/gin-gonic/gin"
)

// Recovery 统一业务错误处理
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if apiErr, ok := err.(api.ApiError); ok {
					c.JSON(apiErr.Code, gin.H{"err": apiErr.Msg})
					c.Abort()
				} else {
					// 非业务错误，继续上抛给框架
					panic(err)
				}
			}
		}()
		c.Next()
	}
}
