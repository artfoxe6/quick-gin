package handlers

import (
	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type AuthorHandler struct {
	service *services.AuthorService
}

func NewAuthorHandler() *AuthorHandler {
	return &AuthorHandler{
		service: services.NewAuthorService(),
	}
}
func (h *AuthorHandler) Create(c *gin.Context) {
	r := new(request.AuthorUpsert)
	api := app.New(c, r)
	err := h.service.Create(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *AuthorHandler) Update(c *gin.Context) {
	r := new(request.AuthorUpsert)
	api := app.New(c, r)
	err := h.service.Update(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *AuthorHandler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	err := h.service.Delete(r.Id)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *AuthorHandler) Detail(c *gin.Context) {
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
func (h *AuthorHandler) List(c *gin.Context) {
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
