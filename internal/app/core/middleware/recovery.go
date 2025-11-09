package middleware

import (
	"log"
	"net/http"

	"github.com/artfoxe6/quick-gin/internal/app/core/apperr"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case *apperr.Error:
					status := e.Code
					if status == 0 {
						status = http.StatusInternalServerError
					}
					message := e.Message
					if message == "" {
						message = http.StatusText(status)
					}
					c.JSON(status, gin.H{"err": message})
					c.Abort()
				case error:
					log.Printf("panic recovered: %v", e)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": http.StatusText(http.StatusInternalServerError)})
				default:
					log.Printf("panic recovered: %v", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": http.StatusText(http.StatusInternalServerError)})
				}
			}
		}()
		c.Next()
	}
}
