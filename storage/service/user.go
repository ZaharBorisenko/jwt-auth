package service

import (
	"context"
	"fmt"
	"github.com/ZaharBorisenko/jwt-auth/helpers/password"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/repositories"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	userRepo    *repositories.UserRepository
	redisClient *storage.RedisClient
}

func NewUserService(userRepo *repositories.UserRepository, redisClient *storage.RedisClient) *UserService {
	return &UserService{userRepo: userRepo, redisClient: redisClient}
}

func (s *UserService) RegisterUser(ctx context.Context, req *models.CreateUserRequestDTO) (*models.User, error) {

	exists, err := s.userRepo.UserExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

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
		Role:      "user",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *models.UserLoginDTO) (*models.UserResponseDTO, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}

	if !password.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &models.UserResponseDTO{
		Id:        user.Id,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	}, nil
}

func (s *UserService) ProfileUser(ctx context.Context, id uuid.UUID) (*models.UserResponseDTO, error) {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}
	return &models.UserResponseDTO{
		Id:        user.Id,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	}, nil
}

func (s *UserService) AllUsers(ctx context.Context) (*[]models.User, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("users not found %w", err)
	}

	return users, nil
}
