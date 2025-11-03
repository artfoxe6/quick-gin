package handlers

import (
	"strconv"

	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/middleware"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	service services.NewsService
}

func NewNewsHandler(service services.NewsService) *NewsHandler {
	return &NewsHandler{service: service}
}

func (h *NewsHandler) Create(c *gin.Context) {
	r := new(request.NewsUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	user, ok := middleware.UserFromContext(c)
	if !ok {
		api.Error(apperr.Unauthorized("authorized invalid"))
		return
	}
	id, err := h.service.Create(r, user)
	if api.Error(err) {
		return
	}
	api.Json(id)
}

func (h *NewsHandler) Update(c *gin.Context) {
	r := new(request.NewsUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Update(r)) {
		return
	}
	api.Json()
}

func (h *NewsHandler) Delete(c *gin.Context) {
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

func (h *NewsHandler) Detail(c *gin.Context) {
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

func (h *NewsHandler) List(c *gin.Context) {
	r := new(request.NewsSearch)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	news, total, err := h.service.List(r)
	if api.Error(err) {
		return
	}
	api.Json(map[string]any{
		"total": total,
		"data":  news,
	})
}
