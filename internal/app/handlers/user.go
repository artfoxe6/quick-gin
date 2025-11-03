package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/artfoxe6/quick-gin/internal/pkg/kit"
	"github.com/artfoxe6/quick-gin/internal/pkg/oss"
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	r := new(request.UserUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Create(r)) {
		return
	}
	api.Json()
}

func (h *UserHandler) Update(c *gin.Context) {
	r := new(request.UserUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Update(r)) {
		return
	}
	api.Json()
}

func (h *UserHandler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Delete(r.Id)) {
		return
	}
	api.Json()
}

func (h *UserHandler) Detail(c *gin.Context) {
	api := app.New(c, nil)
	if api.HasError() {
		return
	}
	idStr := c.Query("id")
	if idStr == "" {
		api.Error(apperr.BadRequest("id is required"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.Error(apperr.BadRequest("id is required"))
		return
	}
	news, err := h.service.Detail(uint(id))
	if api.Error(err) {
		return
	}
	api.Json(news)
}

func (h *UserHandler) List(c *gin.Context) {
	r := new(request.NormalSearch)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	data, total, err := h.service.List(r)
	if api.Error(err) {
		return
	}
	api.Json(map[string]any{
		"total": total,
		"data":  data,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	r := new(request.UserLogin)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	var tokenStr string
	var err error
	if r.Email == config.Super.Email && r.Password == config.Super.Password {
		tokenStr, err = h.service.SuperUserToken(r.Email, r.Password)
	} else {
		tokenStr, err = h.service.Login(r)
	}
	if api.Error(err) {
		return
	}
	api.Json(tokenStr)
}

func (h *UserHandler) Register(c *gin.Context) {
	r := new(request.UserCreate)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	tokenStr, err := h.service.Register(r)
	if api.Error(err) {
		return
	}
	api.Json(tokenStr)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	r := new(request.UpdatePassword)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.UpdatePassword(r)) {
		return
	}
	api.Json()
}

func (h *UserHandler) Code(c *gin.Context) {
	r := new(request.Code)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.SendCode(r)) {
		return
	}
	api.Json()
}

func (h *UserHandler) FreshToken(c *gin.Context) {
	api := app.New(c, nil)
	if api.HasError() {
		return
	}
	oldToken := c.GetHeader("Authorization")
	tokenStr, err := token.Refresh(oldToken)
	if api.Error(err) {
		return
	}
	api.Json(tokenStr)
}

func (h *UserHandler) Upload(c *gin.Context) {
	r := new(request.Upload)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if r.Raw != "" {
		raw := r.Raw
		if strings.HasPrefix(raw, "data:image/png;base64,") {
			raw = raw[22:]
		}
		data, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			api.Error(apperr.BadRequest(err.Error()))
			return
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
			api.Error(apperr.BadRequest("file upload error"))
			return
		}
		api.Json(url)
		return
	}
	if r.File != nil {
		if r.File.Size == 0 || r.File.Size > 500*1024*1024 {
			api.Error(apperr.BadRequest("file_size_err"))
			return
		}
		fileName := fmt.Sprintf("%s/%s/%s/%s", r.Type, time.Now().Format("20060102"), time.Now().Format("150405"), r.File.Filename)
		f, err := r.File.Open()
		if err != nil {
			api.Error(apperr.BadRequest("file open error"))
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			api.Error(apperr.BadRequest("file read error"))
			return
		}
		if r.File.Size > 10*1024*1024 {
			contentType := http.DetectContentType(data)
			if contentType == "image/png" || contentType == "image/jpeg" {
				if compressData, err := kit.Compress(data, contentType); err == nil {
					data = compressData
				}
			}
		}
		url := oss.GetClient().Upload(fileName, bytes.NewReader(data))
		if url == "" {
			api.Error(apperr.BadRequest("file upload error"))
			return
		}
		api.Json(url)
		return
	}
	api.Error(apperr.BadRequest("invalid upload payload"))
}
