package handlers

import (
	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TagHandler struct {
	service *services.TagService
}

func NewTagHandler() *TagHandler {
	return &TagHandler{
		service: services.NewTagService(),
	}
}
func (h *TagHandler) Create(c *gin.Context) {
	r := new(request.TagUpsert)
	api := app.New(c, r)
	err := h.service.Create(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *TagHandler) Update(c *gin.Context) {
	r := new(request.TagUpsert)
	api := app.New(c, r)
	err := h.service.Update(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *TagHandler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	err := h.service.Delete(r.Id)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *TagHandler) Detail(c *gin.Context) {
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
func (h *TagHandler) List(c *gin.Context) {
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
