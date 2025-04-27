package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"text/template"
	"unicode"
)

type GenInfo struct {
	Package   string
	Name      string
	LowerName string
}

func main() {
	var module string
	flag.StringVar(&module, "module", "", "module name")
	flag.Parse()

	if module == "" {
		fmt.Println("module name is required, eg: --module=test")
		return
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Println("must in project root dir")
		return
	}
	info := GenInfo{
		Name:      capitalizeFirstLetter(module),
		LowerName: strings.ToLower(module),
		Package:   getPackage(),
	}

	if isPackageExists(getPackage() + "internal/app/" + info.LowerName) {
		log.Println("Package exists")
		return
	}

	temp := map[string]string{
		"repositories": repositoryTpl,
		"models":       modelTpl,
		"services":     serviceTpl,
		"handlers":     handlerTpl,
	}
	for k, v := range temp {
		genCode(v, "./internal/app/"+k+"/"+strings.ToLower(module)+".go", info)
	}

	genRoute(routeTpl, "./internal/app/router/route.go", info)
}

func genRoute(tpl, path string, info GenInfo) {
	tmpl, err := template.New(info.LowerName).Parse(tpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, info)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
	fp, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fp.Close()
	oldContent, err := io.ReadAll(fp)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	newContent := strings.Replace(string(oldContent), "return r", buffer.String(), 1)

	if _, err = fp.Seek(0, 0); err != nil {
		fmt.Println("Error seeking file:", err)
		return
	}
	if err = fp.Truncate(0); err != nil {
		fmt.Println("Error truncating file:", err)
		return
	}

	_, err = io.WriteString(fp, newContent)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	cmd := exec.Command("go", "fmt", "./internal/app/router/route.go")
	_ = cmd.Run()
	fmt.Println(path)
}

func genCode(tpl, path string, info GenInfo) {
	tmpl, err := template.New(info.LowerName).Parse(tpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	fmt.Println(path)
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	err = tmpl.Execute(file, info)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}
func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func getPackage() string {
	info, _ := debug.ReadBuildInfo()
	if info.Main.Path != "" {
		return info.Main.Path
	}
	out, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(out), "\n", "")
}

var repositoryTpl = `
package repositories

import (
	"{{.Package}}/internal/app/models"
	"{{.Package}}/internal/pkg/db"
)

type {{.Name}}Repository struct {
	Repository[models.{{.Name}}]
}

func New{{.Name}}Repository() *{{.Name}}Repository {
	return &{{.Name}}Repository{
		Repository[models.{{.Name}}]{
			db: db.Db(),
		},
	}
}
`
var modelTpl = `
package models

import "gorm.io/gorm"

type {{.Name}} struct {
	gorm.Model
	Name string ` + "`gorm:\"size:255\"`" + `
}

func (m *{{.Name}}) ToMap() map[string]any {
	return map[string]any{
		"id":         m.ID,
		"name":       m.Name,
		"created_at": m.CreatedAt,
		"updated_at": m.UpdatedAt,
	}
}`

var serviceTpl = `
package services

import (
	"{{.Package}}/internal/app/models"
	"{{.Package}}/internal/app/repositories"
	"{{.Package}}/internal/app/repositories/builder"
	"{{.Package}}/internal/app/request"
)

type {{.Name}}Service struct {
	repository *repositories.{{.Name}}Repository
}

func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{
		repository: repositories.New{{.Name}}Repository(),
	}
}

func (s *{{.Name}}Service) Create(r *request.BaseUpsert) error {
	m := models.{{.Name}}{
		Name: *r.Name,
	}
	if err := s.repository.Create(&m); err != nil {
		return err
	}
	return nil
}

func (s *{{.Name}}Service) Update(r *request.BaseUpsert) error {
	m, err := s.repository.Get(*r.Id)
	if err != nil {
		return err
	}
	if r.Name != nil {
		m.Name = *r.Name
	}
	if err = s.repository.Update(m); err != nil {
		return err
	}
	return nil
}

func (s *{{.Name}}Service) Delete(id uint) error {
	return s.repository.Delete(id)
}
func (s *{{.Name}}Service) Detail(id uint) (any, error) {
	m, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return m.ToMap(), nil
}
func (s *{{.Name}}Service) List(r *request.NormalSearch) (any, int64, error) {
	b := builder.New()
	if r.Keyword != nil {
		b.Like("name", *r.Keyword)
	}
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
	}
	b.Order(orderSet[r.Sort])
	data, total, err := s.repository.ListWithCount(r.Offset, r.Limit, b)
	if err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(data))
	for _, v := range data {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}`

var handlerTpl = `package handlers

import (
	"{{.Package}}/internal/app"
	"{{.Package}}/internal/app/request"
	"{{.Package}}/internal/app/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type {{.Name}}Handler struct {
	service *services.{{.Name}}Service
}

func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{
		service: services.New{{.Name}}Service(),
	}
}
func (h *{{.Name}}Handler) Create(c *gin.Context) {
	r := new(request.BaseUpsert)
	api := app.New(c, r)
	err := h.service.Create(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *{{.Name}}Handler) Update(c *gin.Context) {
	r := new(request.BaseUpsert)
	api := app.New(c, r)
	err := h.service.Update(r)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}

func (h *{{.Name}}Handler) Delete(c *gin.Context) {
	r := new(request.DeleteId)
	api := app.New(c, r)
	err := h.service.Delete(r.Id)
	if err != nil {
		api.Error(err)
	}
	api.Json()
}
func (h *{{.Name}}Handler) Detail(c *gin.Context) {
	api := app.New(c, nil)
	idStr := c.Query("id")
	if idStr == "" {
		api.Error("id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.Error("id is required")
	}
	list, err := h.service.Detail(uint(id))
	if err != nil {
		api.Error(err)
	}
	api.Json(list)
}
func (h *{{.Name}}Handler) List(c *gin.Context) {
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
}`

var routeTpl = `

{{.LowerName}} := handlers.New{{.Name}}Handler()
admin.POST("/{{.LowerName}}/create", {{.LowerName}}.Create)
admin.POST("/{{.LowerName}}/update", {{.LowerName}}.Update)
admin.POST("/{{.LowerName}}/delete", {{.LowerName}}.Delete)
admin.GET("/{{.LowerName}}/detail", {{.LowerName}}.Detail)
admin.GET("/{{.LowerName}}/list", {{.LowerName}}.List)
api.GET("/{{.LowerName}}/detail", {{.LowerName}}.Detail)
api.GET("/{{.LowerName}}/list", {{.LowerName}}.List)

return r`

func isPackageExists(pkgPath string) bool {
	_, err := build.Import(pkgPath, "", build.FindOnly)
	return err == nil
}
