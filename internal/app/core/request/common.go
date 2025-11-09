package request

// BaseUpsert is a lightweight DTO used for simple CRUD modules.
type BaseUpsert struct {
	Id   *uint   `json:"id"`
	Name *string `json:"name" binding:"required"`
}

type DeleteId struct {
	Id uint `json:"id" binding:"required"`
}

type NormalSearch struct {
	Page    int     `form:"page"`
	Limit   int     `form:"limit"`
	Keyword *string `form:"keyword"`
	Sort    int     `form:"sort"`
}

func (r *NormalSearch) Offset() int {
	if r.Page <= 1 {
		return 0
	}
	return (r.Page - 1) * r.Limit
}
