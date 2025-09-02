package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `json:"id" db:"id"`
	UserName  string    `json:"userName" db:"username"`
	Password  string    `json:"password" db:"password"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateUserRequestDTO struct {
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequestDTO struct {
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Config struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}
