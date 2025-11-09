package app

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/artfoxe6/quick-gin/internal/app/core/apperr"
	"github.com/gin-gonic/gin"
)

type Api struct {
	context  *gin.Context
	hasError bool
}

func New(c *gin.Context, request any) *Api {
	api := &Api{context: c}
	if request != nil {
		api.Bind(request)
	}
	return api
}

func (a *Api) Bind(r any) bool {
	if err := a.context.ShouldBind(r); err != nil {
		a.Error(apperr.BadRequest(err.Error()))
		return true
	}
	return false
}

func (a *Api) Error(err error) bool {
	if err == nil {
		return false
	}
	a.hasError = true

	status := http.StatusBadRequest
	message := http.StatusText(status)

	var appErr *apperr.Error
	if errors.As(err, &appErr) {
		if appErr.Code != 0 {
			status = appErr.Code
		}
		if appErr.Message != "" {
			message = appErr.Message
		}
		if appErr.Err != nil {
			err = appErr.Err
		}
	} else {
		if err.Error() != "" {
			message = err.Error()
		}
	}

	if message == "" {
		message = http.StatusText(status)
	}

	a.context.AbortWithStatusJSON(status, gin.H{"err": message})
	return true
}

func (a *Api) HasError() bool {
	return a.hasError
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
