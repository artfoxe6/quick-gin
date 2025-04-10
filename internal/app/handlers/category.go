package handlers

import (
	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		service: services.NewCategoryService(),
	}
}
func (h *CategoryHandler) Create(c *gin.Context) {
	r := new(request.CategoryUpsert)
	api := app.New(c, r)
	err := h.service.Create(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *CategoryHandler) Update(c *gin.Context) {
	r := new(request.CategoryUpsert)
	api := app.New(c, r)
	err := h.service.Update(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	err := h.service.Delete(r.Id)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *CategoryHandler) Detail(c *gin.Context) {
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
func (h *CategoryHandler) List(c *gin.Context) {
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
