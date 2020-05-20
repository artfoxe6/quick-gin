package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ApiError struct {
	Code int
	Msg  string
}

func ReturnApiError(params ...any) {
	var err = ApiError{
		Code: http.StatusBadRequest,
		Msg:  "",
	}
	for _, param := range params {
		t := reflect.TypeOf(param).String()
		switch t {
		case "int":
			err.Code = param.(int)
		case "string":
			err.Msg = param.(string)
		default:
			if e, ok := param.(error); ok {
				s := e.Error()
				err.Msg = s
			}
		}
	}
	panic(err)
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "welcome quick-gin",
	})
}
