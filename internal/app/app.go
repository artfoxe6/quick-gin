package app

import (
	"github.com/artfoxe6/quick-gin/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type Api struct {
	context *gin.Context
}

func New(c *gin.Context, request any) *Api {
	api := &Api{context: c}
	if request != nil {
		api.Bind(request)
	}
	return api
}

func (a *Api) Bind(r any) {
	if err := a.context.Bind(r); err != nil {
		a.Error(err)
	}
}

func (a *Api) Error(params ...any) {
	var err = middleware.ApiError{
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

func (a *Api) Json(data ...any) {
	if data == nil {
		a.context.JSON(http.StatusOK, gin.H{})
		return
	}
	if len(data) == 0 {
		a.context.JSON(http.StatusOK, gin.H{})
		return
	}
	if len(data) == 1 && data[0] == nil {
		a.context.JSON(http.StatusOK, gin.H{})
		return
	}
	for _, d := range data {
		t := reflect.TypeOf(d).Kind()
		switch t {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			a.context.JSON(http.StatusOK, gin.H{"data": d})
		default:
			a.context.JSON(http.StatusOK, d)
		}
		break
	}

}
