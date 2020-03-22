package UserModel

func (user *User) Source() map[string]interface{} {
	return map[string]interface{}{
		"user_name":     user.UserName,
		"age":           user.Age,
		"last_login_at": user.LastLoginAt,
		"password":      user.Password,
	}
}

func (userList *UserList) ToArray() []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, v := range *userList {
		res = append(res, map[string]interface{}{
			"user_name":     v.UserName,
			"age":           v.Age,
			"last_login_at": (*(v.LastLoginAt)).Format("2006-01-02 15:04:05"),
		})
	}
	return res
}

func (userWithArticle *UserWithArticle) Source() map[string]interface{} {
	return map[string]interface{}{
		"id":            userWithArticle.ID,
		"user_name":     userWithArticle.UserName,
		"age":           userWithArticle.Age,
		"last_login_at": (*(userWithArticle.LastLoginAt)).Format("2006-01-02 15:04:05"),
		"articles":      userWithArticle.Articles.ToArray(),
	}
}
