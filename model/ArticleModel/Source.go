package ArticleModel

func (articles *Articles) ToArray() []map[string]interface{} {
	list := []map[string]interface{}{}
	for _, v := range *articles {
		list = append(list, map[string]interface{}{
			"title":   v.Title,
			"content": v.Content,
		})
	}
	return list
}
