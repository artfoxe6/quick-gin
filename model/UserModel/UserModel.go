package UserModel

import (
	"golang.org/x/crypto/bcrypt"
	"quick_gin/config/db"
	"quick_gin/model"
	"quick_gin/model/ArticleModel"
	"time"
)

type User struct {
	model.Base
	UserName    *string    `db:"user_name"`
	Age         *int       `db:"age"`
	Password    *string    `db:"password"`
	LastLoginAt *time.Time `db:"last_login_at"`
}

func (user *User) Add(r map[string]string) error {
	sql := "insert into user (user_name,age,password,created_at) values (?,?,?,?)"
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	b, err := bcrypt.GenerateFromPassword([]byte(r["password"]), bcrypt.MinCost)
	if err != nil {
		return err
	}
	_, err = db.Insert(sql, r["user_name"], r["age"], string(b), createdAt)
	return err
}

func (user *User) Find(r map[string]string) (map[string]interface{}, error) {
	sql := "select password from user where user_name=? limit 1"
	err := db.Get(user, sql, r["user_name"])
	return user.Source(), err
}

type UserList []User

func (userList *UserList) List(r map[string]string) ([]map[string]interface{}, error) {
	sql := "select user_name,last_login_at,age from user "
	err := db.Select(userList, sql)
	if err != nil {
		return nil, err
	}
	return userList.ToArray(), nil
}

type UserWithArticle struct {
	User
	Articles ArticleModel.Articles
}

func (userWithArticle *UserWithArticle) Info(r map[string]string) (map[string]interface{}, error) {
	userSql := "select * from user where id=? limit 1"
	err := db.Get(userWithArticle, userSql, r["id"])
	if err != nil {
		return nil, err
	}
	articleSql := "select uid,title,content from article where uid = ?"

	articles := new(ArticleModel.Articles)

	err = db.Select(articles, articleSql, r["id"])
	if err != nil {
		return nil, err
	}
	userWithArticle.Articles = *articles
	return userWithArticle.Source(), nil
}
