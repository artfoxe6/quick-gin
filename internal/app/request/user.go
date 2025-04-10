package request

import "mime/multipart"

type UserUpsert struct {
	Id       *uint   `json:"id"`
	Avatar   *string `json:"avatar"`
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	Password *string `json:"pass"`
	Email    *string `json:"email"`
	LoginAt  *int64  `json:"login_at"`
	LastIp   *string `json:"last_ip"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type UserCreate struct {
	Avatar   string `json:"avatar"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type UpdatePassword struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type Code struct {
	Email string `json:"email" binding:"required"`
	Type  int    `json:"type"`
}

type Upload struct {
	Type string                `form:"type" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required_without=Raw"`
	Raw  string                `form:"raw"`
}
