package route

import (
	"github.com/ZaharBorisenko/jwt-auth/handlers"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"net/http"
)

func MakeHTTPHandler(userService *service.UserService) http.Handler {
	userHandler := handlers.NewUserHandler(userService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/", userHandler.Register)
	return mux
}
