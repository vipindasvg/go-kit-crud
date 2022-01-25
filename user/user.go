package user

import "context"

type User struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
	EmailId	 string `json:"email_id"`
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	UserLogin(ctx context.Context, EmailId string, Password string) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
}