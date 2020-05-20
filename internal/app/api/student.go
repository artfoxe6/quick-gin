package api

import (
	"github.com/artfoxe6/quick-gin/internal/app/service"
	"github.com/artfoxe6/quick-gin/internal/pkg/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StudentApi struct {
	service service.Student
}

func Student() StudentApi {
	return StudentApi{service: service.Student{
		Db: database.Db(),
	}}
}

func (s StudentApi) Info(c *gin.Context) {
	uid := c.MustGet("uid")
	info := s.service.Info(uid.(int))
	c.JSON(http.StatusOK, gin.H{
		"info": info,
	})
}

// Test 模拟业务错误
func (s StudentApi) Test(c *gin.Context) {

	// 错误会被我们自定义的 Recovery 中间件捕获，返回统一格式的错误信息
	// c.JSON(apiErr.Code, gin.H{"err": apiErr.Msg})
	ReturnApiError("操作失败")
}
