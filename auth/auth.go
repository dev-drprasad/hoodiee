package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateBcryptHash(s string, p string) string {
	password := []byte(fmt.Sprintf("%s%s", s, p))
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CompareBcryptHash(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
}
