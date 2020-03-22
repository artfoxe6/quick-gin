package ArticleService

import (
	"quick_gin/model/ArticleModel"
	"quick_gin/util/request"
)

func Add(r *request.Request) {

	err := r.Validate([]string{"title", "content", "uid"})

	if err != nil {
		r.Error(err.Error())
		return
	}

	articleModel := new(ArticleModel.Article)

	err = articleModel.Insert(r.Posts())

	if err != nil {
		r.Error(err.Error())
		return
	}
	r.Success(nil)
}
