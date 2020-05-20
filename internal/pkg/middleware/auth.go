package middleware

import (
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Auth 基于角色认证
func Auth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "no authorized"})
			return
		}
		data, err := token.Parse(auth)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "authorized invalid"})
			return
		}
		// 是否为指定的角色
		if len(roles) > 0 {
			ok := false
			for _, role := range roles {
				if role == data["Role"].(string) {
					ok = true
				}
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "permission denied"})
			}
		}
		c.Set("uid", int(data["Id"].(float64)))
		c.Next()
	}
}
