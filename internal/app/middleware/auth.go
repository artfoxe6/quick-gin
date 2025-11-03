package middleware

import (
	"log"
	"net/http"

	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"github.com/gin-gonic/gin"
)

const ContextUserKey = "current_user"

type UserProvider interface {
	GetByID(id uint) (*models.User, error)
}

func Auth(userProvider UserProvider, roles ...string) gin.HandlerFunc {
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
		if len(roles) > 0 {
			ok := false
			userRole, _ := data["role"].(string)
			for _, role := range roles {
				if role == userRole {
					ok = true
					break
				}
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "permission denied"})
				return
			}
		}

		uid := uint(data["id"].(float64))
		c.Set("uid", int(uid))

		if userProvider != nil {
			user, err := userProvider.GetByID(uid)
			if err != nil {
				log.Println(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "authorized invalid"})
				return
			}
			c.Set(ContextUserKey, user)
		}

		c.Next()
	}
}

func UserFromContext(c *gin.Context) (*models.User, bool) {
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return nil, false
	}
	user, ok := value.(*models.User)
	return user, ok && user != nil
}
