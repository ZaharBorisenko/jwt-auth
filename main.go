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

	config := models.Config{
		HOST:     os.Getenv("HOST"),
		PORT:     os.Getenv("PORT"),
		USER:     os.Getenv("USER"),
		PASSWORD: os.Getenv("PASSWORD"),
		DBNAME:   os.Getenv("DBNAME"),
	}
	db, err := storage.ConnectDB(&config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	r := route.MakeHTTPHandler(userService)

	log.Println("Server starting on :8080")
	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
