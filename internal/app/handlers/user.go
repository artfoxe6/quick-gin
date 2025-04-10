package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/artfoxe6/quick-gin/internal/pkg/kit"
	"github.com/artfoxe6/quick-gin/internal/pkg/oss"
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: services.NewUserService(),
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	r := new(request.UserUpsert)
	api := app.New(c, r)
	err := h.service.Create(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *UserHandler) Update(c *gin.Context) {
	r := new(request.UserUpsert)
	api := app.New(c, r)
	err := h.service.Update(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *UserHandler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	err := h.service.Delete(r.Id)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *UserHandler) Detail(c *gin.Context) {
	api := app.New(c, nil)
	idStr := c.Query("id")
	if idStr == "" {
		api.Error("id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.Error("id is required")
	}
	news, err := h.service.Detail(uint(id))
	if err != nil {
		api.Error(err)
	}
	api.Json(news)
}
func (h *UserHandler) List(c *gin.Context) {
	r := new(request.NormalSearch)
	api := app.New(c, r)
	data, total, err := h.service.List(r)
	if err != nil {
		api.Error(err)
	}
	api.Json(map[string]any{
		"total": total,
		"data":  data,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	r := new(request.UserLogin)
	api := app.New(c, r)
	var token string
	var err error
	if r.Email == config.Super.Email && r.Password == config.Super.Password {
		token, err = h.service.SuperUserToken(r.Email, r.Password)
	} else {
		token, err = h.service.Login(r)
	}
	if err != nil {
		api.Error(err)
	}
	api.Json(token)
}

func (h *UserHandler) Register(c *gin.Context) {
	r := new(request.UserCreate)
	api := app.New(c, r)

	token, err := h.service.Register(r)
	if err != nil {
		api.Error(err)
	}
	api.Json(token)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	r := new(request.UpdatePassword)
	api := app.New(c, r)

	err := h.service.UpdatePassword(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *UserHandler) Code(c *gin.Context) {
	r := new(request.Code)
	api := app.New(c, r)
	err := h.service.SendCode(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *UserHandler) FreshToken(c *gin.Context) {
	api := app.New(c, nil)
	oldToken := c.GetHeader("Authorization")
	tokenStr, err := token.Refresh(oldToken)
	if err != nil {
		api.Error(err)
	}
	api.Json(tokenStr)
}

func (h *UserHandler) Upload(c *gin.Context) {
	r := new(request.Upload)
	api := app.New(c, r)
	// base64
	if r.Raw != "" {
		if strings.HasPrefix(r.Raw, "data:image/png;base64,") {
			r.Raw = r.Raw[22:]
		}
		data, err := base64.StdEncoding.DecodeString(r.Raw)
		if err != nil {
			api.Error(err.Error())
		}
		contentType := http.DetectContentType(data)
		if contentType == "image/png" || contentType == "image/jpeg" {
			if compressData, err := kit.Compress(data, contentType); err == nil {
				data = compressData
			}
		}
		fileName := fmt.Sprintf("%s/%s/%d%d%s", r.Type, time.Now().Format("20060102"), time.Now().Unix(), len(data), ".jpg")
		url := oss.GetClient().Upload(fileName, bytes.NewReader(data))
		if url == "" {
			api.Error("file upload error")
		}
		api.Json(url)
	} else if r.File != nil {
		//文件上传
		if r.File.Size == 0 || r.File.Size > 500*1024*1024 {
			api.Error("file_size_err")
		}
		fileName := fmt.Sprintf("%s/%s/%s/%s", r.Type, time.Now().Format("20060102"), time.Now().Format("150405"), r.File.Filename)
		f, err := r.File.Open()
		if err != nil {
			api.Error("file open error")
		}

		data, err := io.ReadAll(f)
		if r.File.Size > 10*1024*1024 {
			if err != nil {
				api.Error("file read error")
			}
			contentType := http.DetectContentType(data)
			if contentType == "image/png" || contentType == "image/jpeg" {
				if compressData, err := kit.Compress(data, contentType); err == nil {
					data = compressData
				}
			}
		}
		url := oss.GetClient().Upload(fileName, bytes.NewReader(data))
		if url == "" {
			api.Error("file upload error")
		}
		api.Json(url)
	} else {
		api.Error()
	}
}
