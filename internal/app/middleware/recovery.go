package middleware

import (
	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Code int
	Msg  string
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if apiErr, ok := err.(ApiError); ok {
					c.JSON(apiErr.Code, gin.H{"err": apiErr.Msg})
					c.Abort()
				} else {
					panic(err)
				}
			}
		}()
		c.Next()
	}
}
