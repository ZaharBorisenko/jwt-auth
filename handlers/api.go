package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
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

	loginUser, err := h.userService.LoginUser(r.Context(), &userReq)
	if err != nil {
		WriteERROR(w, http.StatusBadRequest, err.Error())
	}

	WriteJSON(w, http.StatusOK, loginUser)

}
