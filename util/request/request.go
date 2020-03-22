package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Request gin.Context

//通过gin.Context转为Request
func New(c *gin.Context) *Request {
	return (*Request)(c)
}

//将Request反转为 gin.Context
func Ctx(r *Request) *gin.Context {
	return (*gin.Context)(r)
}

//获取GET参数
func (r *Request) Gets() map[string]string {
	values := r.Request.URL.Query()
	temp := map[string]string{}
	for k, v := range values {
		temp[k] = v[0]
	}
	return temp
}

//获取POST参数
func (r *Request) Posts() map[string]string {
	err := r.Request.ParseForm()
	if err != nil {
		log.Fatalf("%v", err)
	}
	values := r.Request.PostForm
	temp := map[string]string{}
	for k, v := range values {
		temp[k] = v[0]
	}
	return temp
}

//获取请求参数
func (r *Request) Inputs() map[string]string {
	if r.Request.Method == "POST" {
		return r.Posts()
	}
	return r.Gets()
}

//获取指定Headers
func (r *Request) Headers() map[string]string {
	temp := map[string]string{}
	for k, v := range r.Request.Header {
		temp[k] = v[0]
	}
	return temp
}

//获取指定Header
func (r *Request) Header(k, d string) string {
	temp := r.Headers()

	value, ok := temp[k]
	if ok {
		return value
	}
	return d
}

//从map中取出一个值，如果key不存在返回默认值
func (r *Request) Get(k, d string) string {
	temp := r.Inputs()

	value, ok := temp[k]
	if ok {
		return value
	}
	return d
}

//错误返回请求
func (r *Request) Error(message interface{}) {

	Ctx(r).JSON(http.StatusOK, map[string]interface{}{
		"data":       "",
		"message":    message,
		"statusCode": http.StatusBadRequest,
	})
}

//成功返回请求
func (r *Request) Success(data interface{}) {

	Ctx(r).JSON(http.StatusOK, map[string]interface{}{
		"data":       data,
		"message":    "",
		"statusCode": http.StatusOK,
	})
}

// 仅验证字段是否缺少
func (r *Request) Validate(list []string) error {
	ctx := Ctx(r)
	if strings.ToUpper(ctx.Request.Method) == "POST" {
		for i := 0; i < len(list); i++ {
			if _, b := ctx.GetPostForm(list[i]); b == false {
				return errors.New("缺少字段：" + list[i])
			}
		}
		return nil
	}
	for i := 0; i < len(list); i++ {
		if _, b := ctx.GetQuery(list[i]); b == false {
			return errors.New("缺少字段：" + list[i])
		}
	}
	return nil
}

//------------ 常见参数封装 ------------

// 获取id参数
func (r *Request) Id() string {
	return r.Get("id", "0")
}

//获取请求中的分页信息
func (r *Request) Page() int {
	temp := r.Get("page", "1")
	page, err := strconv.Atoi(temp)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return page - 1
}

//获取请求中的分页信息
func (r *Request) PerPage() int {
	temp := r.Get("per_page", "10")
	perPage, err := strconv.Atoi(temp)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return perPage
}

//--------------- end ----------------------
