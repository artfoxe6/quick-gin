package handlers

import (
	"github.com/artfoxe6/quick-gin/internal/app"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/app/services"
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NewsHandler struct {
	service *services.NewsService
}

func NewNewsHandler() *NewsHandler {
	return &NewsHandler{
		service: services.NewNewsService(),
	}
}
func (h *NewsHandler) Create(c *gin.Context) {
	r := new(request.NewsUpsert)
	api := app.New(c, r)
	user, err := token.GetUserByToken(c.GetHeader("Authorization"))
	if err != nil {
		api.Error(err)
	}
	id, err := h.service.Create(r, user)
	if err != nil {
		api.Error(err)
	}
	api.Json(id)
}
func (h *NewsHandler) Update(c *gin.Context) {
	r := new(request.NewsUpsert)
	api := app.New(c, r)
	err := h.service.Update(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *NewsHandler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	err := h.service.Delete(r.Id)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *NewsHandler) Detail(c *gin.Context) {
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
func (h *NewsHandler) List(c *gin.Context) {
	r := new(request.NewsSearch)
	api := app.New(c, r)
	news, total, err := h.service.List(r)
	if err != nil {
		api.Error(err)
	}
	api.Json(map[string]any{
		"total": total,
		"data":  news,
	})
}
