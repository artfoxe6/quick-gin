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
		dir string
		tpl string
	}{
		{dir: "repositories", tpl: repositoryTpl},
		{dir: "models", tpl: modelTpl},
		{dir: "services", tpl: serviceTpl},
		{dir: "handlers", tpl: handlerTpl},
	}
	for _, file := range files {
		path := filepath.Join("internal", "app", file.dir, fmt.Sprintf("%s.go", info.LowerName))
		if err := genCode(file.tpl, path, info); err != nil {
			fmt.Printf("error generating %s: %v\n", path, err)
		}
	}

	if err := genRoute(routeTpl, filepath.Join("internal", "app", "router", "route.go"), info); err != nil {
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
	routeInit := fmt.Sprintf("%s := handlers.New%sHandler()", info.LowerName, info.Name)
	if strings.Contains(string(oldContent), routeInit) {
		fmt.Printf("router already contains handlers for %s, skipping route update\n", info.LowerName)
		return nil
	}
	const returnMarker = "return r"
	newContent := strings.Replace(string(oldContent), returnMarker, buffer.String(), 1)
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
"strconv"

"{{.Package}}/internal/app"
"{{.Package}}/internal/app/apperr"
"{{.Package}}/internal/app/request"
"{{.Package}}/internal/app/services"
"github.com/gin-gonic/gin"
)

type {{.Name}}Handler struct {
service services.{{.Name}}Service
}

func New{{.Name}}Handler(service services.{{.Name}}Service) *{{.Name}}Handler {
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

{{.LowerName}} := handlers.New{{.Name}}Handler()
admin.POST("/{{.LowerName}}/create", {{.LowerName}}.Create)
admin.POST("/{{.LowerName}}/update", {{.LowerName}}.Update)
admin.POST("/{{.LowerName}}/delete", {{.LowerName}}.Delete)
admin.GET("/{{.LowerName}}/detail", {{.LowerName}}.Detail)
admin.GET("/{{.LowerName}}/list", {{.LowerName}}.List)
api.GET("/{{.LowerName}}/detail", {{.LowerName}}.Detail)
api.GET("/{{.LowerName}}/list", {{.LowerName}}.List)

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
