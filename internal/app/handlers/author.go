package handlers

import (
	"strconv"

	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	service services.AuthorService
}

func NewAuthorHandler(service services.AuthorService) *AuthorHandler {
	return &AuthorHandler{service: service}
}

func (h *AuthorHandler) Create(c *gin.Context) {
	r := new(request.AuthorUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Create(r)) {
		return
	}
	api.Json()
}

func (h *AuthorHandler) Update(c *gin.Context) {
	r := new(request.AuthorUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Update(r)) {
		return
	}
	api.Json()
}

func (h *AuthorHandler) Delete(c *gin.Context) {
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

func (h *AuthorHandler) Detail(c *gin.Context) {
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

func (h *AuthorHandler) List(c *gin.Context) {
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
