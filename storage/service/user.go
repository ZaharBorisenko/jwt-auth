package service

import (
	"context"
	"github.com/ZaharBorisenko/jwt-auth/helper/password"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/storage/repositories"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *models.CreateUserRequestDTO) (*models.User, error) {
	// 1. Проверить, не существует ли уже пользователь с таким email/username
	// (это вызовет метод репозитория GetUserByEmail, который ты ещё не написал)

	hashPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Id:        uuid.New(),
		UserName:  req.UserName,
		Password:  hashPassword,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}
