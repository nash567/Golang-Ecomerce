package model

type User struct {
	UserID    int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

type UpdateUser struct {
	UserID    int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UserCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}
type QueryParams struct {
	FirstName string
	LastName  string
	Email     string
	Archived  string
	Sort      string
	Order     string
	Limit     int
	Page      int
}
