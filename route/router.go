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
	mux.Handle("POST /register", middleware.RateLimitMiddleware(1, 2)(http.HandlerFunc(userHandler.Register)))
	mux.Handle("POST /login", middleware.RateLimitMiddleware(2, 5)(http.HandlerFunc(userHandler.Login)))
	mux.Handle("POST /logout", middleware.RateLimitMiddleware(5, 10)(http.HandlerFunc(userHandler.Logout)))

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedProfile := middleware.RateLimitMiddleware(2, 5)(http.HandlerFunc(userHandler.GetProfile))
	protectedUpdate := middleware.RateLimitMiddleware(2, 5)(http.HandlerFunc(userHandler.UpdateUser))

	protectedMux.Handle("GET /profile/{id}", protectedProfile)
	protectedMux.Handle("PUT /profile/{id}", protectedUpdate)

	//adminOnly routes
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("GET /admin/users", userHandler.GetAllUsers)
	adminMux.HandleFunc("DELETE /admin/user/{id}", userHandler.DeleteUser)
	adminMux.HandleFunc("GET /admin/blacklist", userHandler.GetBlackList)

	//middleware
	protectedWithAuth := middleware.AuthMiddleware(redisClient, protectedMux)
	withAdminAuth := middleware.AuthMiddleware(redisClient, middleware.RateLimitMiddleware(10, 20)(middleware.AdminOnly(userRepo, adminMux)))
	mux.Handle("/", protectedWithAuth)

	handlerWithLogging := middleware.Logging(mux)
	mux.Handle("/admin/", withAdminAuth)
	return handlerWithLogging

}
