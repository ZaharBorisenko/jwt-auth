package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ZaharBorisenko/jwt-auth/helpers/jwtToken"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	userService *service.UserService
	redisClient *storage.RedisClient
}

func NewUserHandler(userService *service.UserService, redisClient *storage.RedisClient) *UserHandler {
	return &UserHandler{userService: userService, redisClient: redisClient}
}

type Error struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Printf("could not write json %v", err)
	}
}

func WriteERROR(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, Error{Error: message})
}
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteERROR(w, http.StatusMethodNotAllowed, "method not allowed!")
		return
	}
	defer r.Body.Close()

	userReq := models.CreateUserRequestDTO{}
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if userReq.UserName == "" {
		WriteERROR(w, http.StatusBadRequest, "username is required")
		return
	}
	if userReq.Email == "" {
		WriteERROR(w, http.StatusBadRequest, "email is required")
		return
	}
	if userReq.Password == "" {
		WriteERROR(w, http.StatusBadRequest, "password is required")
		return
	}

	createdUser, err := h.userService.RegisterUser(r.Context(), &userReq)
	if err != nil {
		WriteERROR(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.UserResponseDTO{
		Id:        createdUser.Id,
		UserName:  createdUser.UserName,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
		Role:      createdUser.Role,
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteERROR(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	defer r.Body.Close()

	userReq := models.UserLoginDTO{}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if userReq.Email == "" || userReq.Password == "" {
		WriteERROR(w, http.StatusBadRequest, "email or password is required")
		return
	}

	//authentication user
	loginUser, err := h.userService.LoginUser(r.Context(), &userReq)
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	//generate JWT token
	token, err := jwtToken.GenerateJWTToken(loginUser.Id, loginUser.Email)
	if err != nil {
		WriteERROR(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	//return user and token
	response := map[string]interface{}{
		"user":  loginUser,
		"token": token,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteERROR(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		WriteERROR(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		WriteERROR(w, http.StatusUnauthorized, "Invalid authorization format")
		return
	}

	tokenString := parts[1]

	expiration := 72 * time.Hour
	err := h.redisClient.AddToBlackList(r.Context(), tokenString, expiration)
	if err != nil {
		WriteERROR(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		WriteERROR(w, http.StatusBadRequest, "ID parameter is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	user, err := h.userService.ProfileUser(r.Context(), id)
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, user)

}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.AllUsers(r.Context())
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, users)
}
