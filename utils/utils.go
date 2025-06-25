package utils

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func MustHash(password string) string {
	hash, _ := HashPassword(password)
	return hash
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSlug(name string) string {
	// Simple slug generator for appointments/patients
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
