package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Не удалось зашифровать пароль %v", err)
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(HashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashPassword), []byte(password))
	return err == nil
}
