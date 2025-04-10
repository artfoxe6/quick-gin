package request

type ContactUpsert struct {
	Id      *uint   `json:"id"`
	Email   *string `json:"email"`
	Message *string `json:"message"`
	Name    *string `json:"name"`
}
