package  postgres

type UserAccount struct {
	Username string `json:"username" form:"username" query:"username"`
}