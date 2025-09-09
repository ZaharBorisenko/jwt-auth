package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id         uuid.UUID `json:"id" db:"id"`
	UserName   string    `json:"userName" db:"username"`
	Password   string    `json:"password" db:"password"`
	Email      string    `json:"email" db:"email"`
	FirstName  string    `json:"firstName" db:"first_name"`
	LastName   string    `json:"lastName" db:"last_name"`
	Role       string    `json:"role" db:"role"`
	IsVerified bool      `json:"isVerified" db:"is_verified"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateUserRequestDTO struct {
	UserName  string `json:"userName" validate:"required,min=3,max=50"`
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	Email     string `json:"email" validate:"required,email"`
}

type UserResponseDTO struct {
	Id        uuid.UUID `json:"id"`
	UserName  string    `json:"userName"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}
type UserLoginDTO struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	PasswordConfirm string `validate:"required,eqfield=Password"`
}

type UpdateUserRequestDTO struct {
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type VerificationEmailDto struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type Config struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}

type ConfigURLParams struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

type ConfigEmail struct {
	Username string
	Password string
	Host     string
	Port     string
}
