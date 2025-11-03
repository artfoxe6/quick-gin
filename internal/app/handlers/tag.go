package handlers

import (
	"strconv"

	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	service services.TagService
}

func NewTagHandler(service services.TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) Create(c *gin.Context) {
	r := new(request.TagUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Create(r)) {
		return
	}
	api.Json()
}

func (h *TagHandler) Update(c *gin.Context) {
	r := new(request.TagUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Update(r)) {
		return
	}
	api.Json()
}

func (h *TagHandler) Delete(c *gin.Context) {
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

func (h *TagHandler) Detail(c *gin.Context) {
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

func (h *TagHandler) List(c *gin.Context) {
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
