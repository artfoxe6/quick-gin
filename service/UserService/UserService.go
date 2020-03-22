package UserService

import (
	"golang.org/x/crypto/bcrypt"
	"quick_gin/model/UserModel"
	"quick_gin/util/request"
)

//添加用户
func Add(r *request.Request) {
	model := new(UserModel.User)
	err := model.Add(r.Inputs())
	if err != nil {
		r.Error(err.Error())
		return
	}
	r.Success(nil)
}

//登录
func Login(r *request.Request) {
	err := r.Validate([]string{"user_name", "password"})
	if err != nil {
		r.Error(err.Error())
		return
	}
	model := new(UserModel.User)
	data, err := model.Find(r.Inputs())
	if err != nil {
		r.Error("账号不存在")
		return
	}
	password := []byte(r.Get("password", ""))
	err = bcrypt.CompareHashAndPassword([]byte(data["password"].(string)), password)
	if err != nil {
		r.Error("密码错误")
		return
	}
	r.Success(nil)
}

//用户列表
func List(r *request.Request) {
	userList := new(UserModel.UserList)
	data, err := userList.List(r.Gets())
	if err != nil {
		r.Error(err.Error())
		return
	}
	r.Success(data)
}

//用户信息以及用户发布的文章
func InfoWithArticle(r *request.Request) {

	err := r.Validate([]string{"id"})
	if err != nil {
		r.Error(err.Error())
		return
	}
	m := new(UserModel.UserWithArticle)
	list, err := m.Info(r.Gets())
	if err != nil {
		r.Error(err.Error())
		return
	}
	r.Success(list)
}
