package handlers

import (
	"context"
	"github.com/ZaharBorisenko/jwt-auth/helpers/JSON"
	"github.com/ZaharBorisenko/jwt-auth/helpers/jwtToken"
	"github.com/ZaharBorisenko/jwt-auth/helpers/parseUUID"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"github.com/ZaharBorisenko/jwt-auth/validator"
	"net/http"
	"strconv"
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

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSON.WriteERROR(w, http.StatusMethodNotAllowed, "method not allowed!")
		return
	}
	defer r.Body.Close()

	userReq := models.CreateUserRequestDTO{}
	if !validator.ValidateRequest(w, r, &userReq) {
		return
	}

	createdUser, err := h.userService.RegisterUser(r.Context(), &userReq)
	if err != nil {
		JSON.WriteERROR(w, http.StatusInternalServerError, err.Error())
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

	JSON.WriteJSON(w, http.StatusCreated, response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSON.WriteERROR(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	defer r.Body.Close()

	userReq := models.UserLoginDTO{}
	if !validator.ValidateRequest(w, r, &userReq) {
		return
	}

	//authentication user
	loginUser, err := h.userService.LoginUser(r.Context(), &userReq)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	//generate JWT token
	token, err := jwtToken.GenerateJWTToken(loginUser.Id, loginUser.Email)
	if err != nil {
		JSON.WriteERROR(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	//return user and token
	response := map[string]interface{}{
		"user":  loginUser,
		"token": token,
	}

	JSON.WriteJSON(w, http.StatusOK, response)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSON.WriteERROR(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		JSON.WriteERROR(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		JSON.WriteERROR(w, http.StatusUnauthorized, "Invalid authorization format")
		return
	}

	tokenString := parts[1]

	expiration := 72 * time.Hour
	err := h.redisClient.AddToBlackList(r.Context(), tokenString, expiration)
	if err != nil {
		JSON.WriteERROR(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	JSON.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID.ParseUUID(r)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.ProfileUser(r.Context(), id)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON.WriteJSON(w, http.StatusOK, user)

}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	pagination := models.PaginationConfig{
		Offset: offset,
		Limit:  limit,
	}

	users, err := h.userService.AllUsers(r.Context(), pagination)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID.ParseUUID(r)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userService.DeleteUser(context.Background(), id)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON.WriteJSON(w, http.StatusOK, map[string]string{"status": "user successfully deleted"})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID.ParseUUID(r)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	userReq := models.UpdateUserRequestDTO{}
	if !validator.ValidateRequest(w, r, &userReq) {
		return
	}

	userUpdate, err := h.userService.UpdateUser(context.Background(), id, &userReq)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			JSON.WriteERROR(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "email already taken") {
			JSON.WriteERROR(w, http.StatusConflict, err.Error())
		} else {
			JSON.WriteERROR(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	JSON.WriteJSON(w, http.StatusOK, userUpdate)
}

func (h *UserHandler) GetBlackList(w http.ResponseWriter, r *http.Request) {
	keys, err := h.redisClient.GetBlackListKeys(r.Context())
	if err != nil {
		JSON.WriteERROR(w, http.StatusInternalServerError, "Failed to get blacklist")
		return
	}

	JSON.WriteJSON(w, http.StatusOK, keys)
}
