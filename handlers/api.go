package handlers

import (
	"context"
	"encoding/json"
	"github.com/ZaharBorisenko/jwt-auth/helpers"
	"github.com/ZaharBorisenko/jwt-auth/helpers/JSON"
	"github.com/ZaharBorisenko/jwt-auth/helpers/jwtToken"
	"github.com/ZaharBorisenko/jwt-auth/helpers/parseUUID"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/services"
	"github.com/ZaharBorisenko/jwt-auth/validator"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	userService *services.UserService
	redisClient *storage.RedisClient
}

func NewUserHandler(userService *services.UserService, redisClient *storage.RedisClient) *UserHandler {
	return &UserHandler{userService: userService, redisClient: redisClient}
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param request body models.CreateUserRequestDTO true "User registration data"
// @Success 201 {object} models.UserResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param request body models.UserLoginDTO true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
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

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logout [post]
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

// GetProfile godoc
// @Summary Get user profile
// @Description Get user profile by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile/{id} [get]
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Get paginated list of all users (Admin only)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort order" default(desc)
// @Success 200 {array} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	params := helpers.ParseQueryParams(r)

	users, err := h.userService.AllUsers(r.Context(), params)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON.WriteJSON(w, http.StatusOK, users)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID (Admin only)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/user/{id} [delete]
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

// UpdateUser godoc
// @Summary Update user profile
// @Description Update user profile by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequestDTO true "User update data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile/{id} [put]
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

// GetBlackList godoc
// @Summary Get blacklisted tokens
// @Description Get list of blacklisted JWT tokens (Admin only)
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/blacklist [get]
func (h *UserHandler) GetBlackList(w http.ResponseWriter, r *http.Request) {
	keys, err := h.redisClient.GetBlackListKeys(r.Context())
	if err != nil {
		JSON.WriteERROR(w, http.StatusInternalServerError, "Failed to get blacklist")
		return
	}

	JSON.WriteJSON(w, http.StatusOK, keys)
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify user email with verification code
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param request body models.VerificationEmailDto true "Verification data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /verify-email [post]
func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := models.VerificationEmailDto{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	err := h.userService.VerifyEmail(r.Context(), request.Email, request.Code)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Email successfully verified",
	})
}

// ResendVerificationCode godoc
// @Summary Resend verification code
// @Description Resend email verification code
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param request body models.ResendVerificationRequest true "Email address"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /resend-verification [post]
func (h *UserHandler) ResendVerificationCode(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request models.ResendVerificationRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	err := h.userService.ResendVerificationCode(r.Context(), request.Email)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Verification code resent",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for authenticated user
// @Tags authentication
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body models.ChangePasswordDto true "Change password request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /change-password [post]
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	req := models.ChangePasswordDto{}
	if !validator.ValidateRequest(w, r, &req) {
		return
	}

	log.Println("Request data:", req)

	err := h.userService.ChangePassword(r.Context(), &req)
	if err != nil {
		JSON.WriteERROR(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Password changed for email:", req.Email)

	JSON.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}
