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
	UserName  string `json:"userName" validate:"required,min=3,max=50" example:"johndoe"`
	FirstName string `json:"firstName" validate:"required,min=2,max=50" example:"John"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50" example:"Doe"`
	Password  string `json:"password" validate:"required,min=6" example:"password123"`
	Email     string `json:"email" validate:"required,email" example:"john.doe@example.com"`
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
	Email           string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password        string `json:"password" validate:"required,min=6" example:"password123"`
	PasswordConfirm string `validate:"required,eqfield=Password"`
}

type UpdateUserRequestDTO struct {
	UserName  string `json:"userName" example:"johndoe_updated"`
	FirstName string `json:"firstName" example:"John_updated"`
	LastName  string `json:"lastName" example:"Doe_updated"`
	Email     string `json:"email" example:"john.updated@example.com"`
}

type VerificationEmailDto struct {
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Code  string `json:"code" validate:"required" example:"123456"`
}

type ChangePasswordDto struct {
	Email       string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	OldPassword string `json:"OldPassword" validate:"required,min=6" example:"oldPassword123"`
	NewPassword string `json:"NewPassword" validate:"required,min=6" example:"newPassword456"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
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
