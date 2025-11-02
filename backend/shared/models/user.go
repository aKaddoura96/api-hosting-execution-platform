package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                   string    `json:"id"`
	Email                string    `json:"email"`
	PasswordHash         string    `json:"-"`
	Name                 string    `json:"name"`
	Role                 string    `json:"role"` // "developer" or "consumer"
	EmailVerified        bool      `json:"email_verified"`
	VerificationToken    string    `json:"-"`
	PasswordResetToken   string    `json:"-"`
	PasswordResetExpires *time.Time `json:"-"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
