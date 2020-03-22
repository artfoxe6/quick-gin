package ArticleModel

import (
	"quick_gin/config/db"
	"quick_gin/model"
)

type Article struct {
	model.Base
	Uid      int    `db:"uid"`
	Title    string `db:"title"`
	Content  string `db:"content"`
	Favorite int    `db:"favorite"`
}

type Articles []Article

func (article *Article) Insert(r map[string]string) error {

	sql := "insert into article (title,content,uid) values (?,?,?)"

	_, err := db.Insert(sql, r["title"], r["content"], r["uid"])
	return err
}
