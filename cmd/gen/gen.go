package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	pkg := strings.TrimSpace(getPackage())
	if pkg == "" {
		fmt.Println("unable to resolve module path")
		return
	}
	info := GenInfo{
		Name:      capitalizeFirstLetter(module),
		LowerName: strings.ToLower(module),
		Package:   pkg,
	}

	if packageExists(info) {
		fmt.Printf("module %q already exists, skipping\n", info.LowerName)
		return
	}

	files := []struct {
		path string
		tpl  string
	}{
		{path: filepath.Join("internal", "app", info.LowerName, "repo", "repo.go"), tpl: repositoryTpl},
		{path: filepath.Join("internal", "app", info.LowerName, "model", "model.go"), tpl: modelTpl},
		{path: filepath.Join("internal", "app", info.LowerName, "service", "service.go"), tpl: serviceTpl},
		{path: filepath.Join("internal", "app", info.LowerName, "handler", "handler.go"), tpl: handlerTpl},
	}
	for _, file := range files {
		if err := genCode(file.tpl, file.path, info); err != nil {
			fmt.Printf("error generating %s: %v\n", file.path, err)
		}
	}

	if err := genRoute(routeTpl, filepath.Join("internal", "app", "core", "router", "route.go"), info); err != nil {
		fmt.Printf("error updating router: %v\n", err)
	}
}

func genRoute(tpl, path string, info GenInfo) error {
	tmpl, err := template.New(info.LowerName).Parse(tpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, info)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	fp, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer fp.Close()
	oldContent, err := io.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	routeInit := fmt.Sprintf("%sHandler := %sHandler.New%sHandler", info.LowerName, info.LowerName, info.Name)
	if strings.Contains(string(oldContent), routeInit) {
		fmt.Printf("router already contains handlers for %s, skipping route update\n", info.LowerName)
		return nil
	}
	newContent := string(oldContent)

	imports := []string{
		fmt.Sprintf("%sHandler \"%s/internal/app/%s/handler\"", info.LowerName, info.Package, info.LowerName),
		fmt.Sprintf("%sRepo \"%s/internal/app/%s/repo\"", info.LowerName, info.Package, info.LowerName),
		fmt.Sprintf("%sService \"%s/internal/app/%s/service\"", info.LowerName, info.Package, info.LowerName),
	}

	for _, imp := range imports {
		if strings.Contains(newContent, imp) {
			continue
		}
		const importMarker = "import ("
		idx := strings.Index(newContent, importMarker)
		if idx == -1 {
			return fmt.Errorf("route file missing import block")
		}
		idx += len(importMarker)
		newContent = newContent[:idx] + "\n\t" + imp + newContent[idx:]
	}

	const returnMarker = "return r"
	newContent = strings.Replace(newContent, returnMarker, buffer.String(), 1)
	if newContent == string(oldContent) {
		return fmt.Errorf("route file missing marker %q", returnMarker)
	}

	if _, err = fp.Seek(0, 0); err != nil {
		return fmt.Errorf("seek file: %w", err)
	}
	if err = fp.Truncate(0); err != nil {
		return fmt.Errorf("truncate file: %w", err)
	}

	_, err = io.WriteString(fp, newContent)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	if err := exec.Command("gofmt", "-w", path).Run(); err != nil {
		return fmt.Errorf("format route: %w", err)
	}
	fmt.Println(path)
	return nil
}

func genCode(tpl, path string, info GenInfo) error {
	if fileExists(path) {
		fmt.Printf("%s already exists, skipping\n", path)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	tmpl, err := template.New(info.LowerName).Parse(tpl)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()
	err = tmpl.Execute(file, info)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	fmt.Println(path)
	return nil
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
package repo

import (
	"{{.Package}}/internal/app/core/repository"
	"{{.Package}}/internal/app/{{.LowerName}}/model"
	"{{.Package}}/internal/pkg/db"
)

type {{.Name}}Repository struct {
	repository.Repository[model.{{.Name}}]
}

func New{{.Name}}Repository() *{{.Name}}Repository {
	return &{{.Name}}Repository{
		repository.Repository[model.{{.Name}}]{
			db: db.Db(),
		},
	}
}
`
var modelTpl = `
package model

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
package service

import (
	"{{.Package}}/internal/app/core/repository/builder"
	"{{.Package}}/internal/app/core/request"
	"{{.Package}}/internal/app/{{.LowerName}}/model"
	"{{.Package}}/internal/app/{{.LowerName}}/repo"
)

type {{.Name}}Service struct {
	repository *repo.{{.Name}}Repository
}

func New{{.Name}}Service(repository *repo.{{.Name}}Repository) *{{.Name}}Service {
	return &{{.Name}}Service{
		repository: repository,
	}
}

func (s *{{.Name}}Service) Create(r *request.BaseUpsert) error {
	m := model.{{.Name}}{
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
	data, total, err := s.repository.ListWithCount(r.Offset(), r.Limit, b)
	if err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(data))
	for _, v := range data {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}`

var handlerTpl = `package handler

import (
	"strconv"

	app "{{.Package}}/internal/app/core"
	"{{.Package}}/internal/app/core/apperr"
	"{{.Package}}/internal/app/core/request"
	"{{.Package}}/internal/app/{{.LowerName}}/service"
	"github.com/gin-gonic/gin"
)

type {{.Name}}Handler struct {
	service *service.{{.Name}}Service
}

func New{{.Name}}Handler(service *service.{{.Name}}Service) *{{.Name}}Handler {
	return &{{.Name}}Handler{service: service}
}
func (h *{{.Name}}Handler) Create(c *gin.Context) {
	r := new(request.BaseUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Create(r)) {
		return
	}
	api.Json()
}
func (h *{{.Name}}Handler) Update(c *gin.Context) {
	r := new(request.BaseUpsert)
	api := app.New(c, r)
	if api.HasError() {
		return
	}
	if api.Error(h.service.Update(r)) {
		return
	}
	api.Json()
}

func (h *{{.Name}}Handler) Delete(c *gin.Context) {
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
func (h *{{.Name}}Handler) Detail(c *gin.Context) {
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
	list, err := h.service.Detail(uint(id))
	if api.Error(err) {
		return
	}
	api.Json(list)
}
func (h *{{.Name}}Handler) List(c *gin.Context) {
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
}`

var routeTpl = `

	{{.LowerName}}Repository := {{.LowerName}}Repo.New{{.Name}}Repository()
	{{.LowerName}}Service := {{.LowerName}}Service.New{{.Name}}Service({{.LowerName}}Repository)
	{{.LowerName}}Handler := {{.LowerName}}Handler.New{{.Name}}Handler({{.LowerName}}Service)
	admin.POST("/{{.LowerName}}/create", {{.LowerName}}Handler.Create)
	admin.POST("/{{.LowerName}}/update", {{.LowerName}}Handler.Update)
	admin.POST("/{{.LowerName}}/delete", {{.LowerName}}Handler.Delete)
	admin.GET("/{{.LowerName}}/detail", {{.LowerName}}Handler.Detail)
	admin.GET("/{{.LowerName}}/list", {{.LowerName}}Handler.List)
	api.GET("/{{.LowerName}}/detail", {{.LowerName}}Handler.Detail)
	api.GET("/{{.LowerName}}/list", {{.LowerName}}Handler.List)

return r`

func packageExists(info GenInfo) bool {
	dir := filepath.Join("internal", "app", info.LowerName)
	stat, err := os.Stat(dir)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return stat.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return !info.IsDir()
}
