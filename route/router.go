package route

import (
	"github.com/ZaharBorisenko/jwt-auth/handlers"
	"github.com/ZaharBorisenko/jwt-auth/middleware"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"net/http"
)

func MakeHTTPHandler(userService *service.UserService, redisClient *storage.RedisClient) http.Handler {
	mux := http.NewServeMux()
	userHandler := handlers.NewUserHandler(userService, redisClient)

	// Public routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /logout", userHandler.Logout)

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /profile/{id}", userHandler.GetProfile)

	// Middleware for private routes
	mux.Handle("/", middleware.AuthMiddleware(redisClient, protectedMux))

	return mux
}
