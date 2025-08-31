package main

import (
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/ZaharBorisenko/jwt-auth/route"
	"github.com/ZaharBorisenko/jwt-auth/storage"
	"github.com/ZaharBorisenko/jwt-auth/storage/repositories"
	"github.com/ZaharBorisenko/jwt-auth/storage/service"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := redisHost + ":" + redisPort

	config := models.Config{
		HOST:     os.Getenv("DB_HOST"),
		PORT:     os.Getenv("DB_PORT"),
		USER:     os.Getenv("DB_USER"),
		PASSWORD: os.Getenv("DB_PASSWORD"),
		DBNAME:   os.Getenv("DB_NAME"),
	}
	db, err := storage.ConnectDB(&config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := storage.NewRedisClient(redisAddr)
	userRepo := repositories.NewUserRepository(db)
	userService := service.NewUserService(userRepo, redisClient)
	r := route.MakeHTTPHandler(userService, userRepo, redisClient)

	log.Println("Server starting on :8080")
	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
