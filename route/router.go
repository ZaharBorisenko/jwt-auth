package route

import (
	"github.com/ZaharBorisenko/jwt-auth/handlers"
	"github.com/ZaharBorisenko/jwt-auth/middleware"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/repositories"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"net/http"
)

func MakeHTTPHandler(userService *service.UserService, userRepo *repositories.UserRepository, redisClient *storage.RedisClient) http.Handler {
	mux := http.NewServeMux()
	userHandler := handlers.NewUserHandler(userService, redisClient)

	// Public routes
	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("POST /logout", userHandler.Logout)

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /profile/{id}", userHandler.GetProfile)
	protectedMux.HandleFunc("PUT /profile/{id}", userHandler.UpdateUser)

	//adminOnly routes
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /admin/users", userHandler.GetAllUsers)
	adminMux.HandleFunc("DELETE /admin/user/{id}", userHandler.DeleteUser)

	//middleware
	protectedWithAuth := middleware.AuthMiddleware(redisClient, protectedMux)
	withAdminAuth := middleware.AuthMiddleware(redisClient, middleware.AdminOnly(userRepo, adminMux))
	mux.Handle("/", protectedWithAuth)

	handlerWithLogging := middleware.Logging(mux)
	mux.Handle("/admin/", withAdminAuth)
	return handlerWithLogging

}
