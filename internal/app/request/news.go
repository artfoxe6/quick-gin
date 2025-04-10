package request

type NewsUpsert struct {
	Id          *uint   `json:"id"`
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	Source      *string `json:"source"`
	Status      *int    `json:"status"`
	Slug        *string `json:"slug"`
	Summary     *string `json:"summary"`
	ImageUrl    *string `json:"image_url"`
	VideoUrl    *string `json:"video_url"`
	AudioUrl    *string `json:"audio_url"`
	TypeId      *int    `json:"type_id"`
	CategoryIds *[]uint `json:"category_ids"`
	TagIds      *[]uint `json:"tag_ids"`
	AuthorId    *uint   `json:"author_id"`
	IsFeatured  *bool   `json:"is_featured"`
}

type DeleteId struct {
	Id uint `json:"id"`
}

type NormalSearch struct {
	Offset  int     `form:"offset,default=0"`
	Limit   int     `form:"limit,default=15"`
	Sort    int     `form:"sort,default=0"`
	Keyword *string `form:"keyword"`
	Ids     *string `form:"ids"`
	Status  *int    `form:"status"`
}

type NewsSearch struct {
	NormalSearch
	TypeIds     *string `form:"type_ids"`
	CategoryIds *string `form:"category_ids"`
	TagIds      *string `form:"tag_ids"`
}

type CategoryUpsert struct {
	Id   *uint   `json:"id"`
	Name *string `json:"name"`
}

type TagUpsert struct {
	Id   *uint   `json:"id"`
	Name *string `json:"name"`
}
type AuthorUpsert struct {
	Id   *uint   `json:"id"`
	Name *string `json:"name"`
}

type BaseUpsert struct {
	Id   *uint   `json:"id"`
	Name *string `json:"name"`
}
