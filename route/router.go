package route

import (
	"github.com/ZaharBorisenko/jwt-auth/handlers"
	"github.com/ZaharBorisenko/jwt-auth/middleware"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"net/http"
)

func MakeHTTPHandler(userService *service.UserService) http.Handler {
	mux := http.NewServeMux()
	userHandler := handlers.NewUserHandler(userService)

	// Public routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /profile/{id}", userHandler.GetProfile)

	// Middleware for private routes
	mux.Handle("/", middleware.AuthMiddleware(protectedMux))

	return mux
}
