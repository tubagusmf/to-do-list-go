package model

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextAuthKey string

const BearerAuthKey ContextAuthKey = "BearerAuth"

type IUserRepository interface {
	FindAll(ctx context.Context, user User) ([]*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) *User
	Create(ctx context.Context, user User) (*User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int64) error
}

type IUserUsecase interface {
	FindAll(ctx context.Context, user User) ([]*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	Login(ctx context.Context, in LoginInput) (token string, err error)
	Create(ctx context.Context, in CreateUserInput) (token string, err error)
	Update(ctx context.Context, id int64, in UpdateUserInput) error
	Delete(ctx context.Context, id int64) error
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
type User struct {
	Id        int64      `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type LoginInput struct {
	Id       int64  `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserInput struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type UpdateUserInput struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}
