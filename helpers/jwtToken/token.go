package jwtToken

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

func GenerateJWTToken(userID uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID.String(), // subject - user id
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // 3 day
		"iat":   time.Now().Unix(),                     // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}
