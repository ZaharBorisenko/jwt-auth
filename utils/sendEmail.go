package utils

import (
	"fmt"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

func SendEmail(email string, code string) error {
	c := models.ConfigEmail{
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
	}

	from := c.Username
	to := []string{email}

	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to[0])
	message += fmt.Sprintf("Subject: %s\r\n", "Email Verification Code")
	message += fmt.Sprintf("\r\nYour verification code is: %s\r\n", code)
	message += "This code will expire in 15 minutes.\r\n"

	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)

	err := smtp.SendMail(c.Host+":"+c.Port, auth, from, to, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully Gmail.")
	return nil
}

func GenerateCodeEmail() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
