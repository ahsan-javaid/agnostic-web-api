package utils

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Check(e error) {
	if e != nil {
			panic(e)
	}
}

func GetCollectionName(url string) string {
	return strings.Split(url, "/")[1]
}

func GetParams(url string) []string {
	return strings.Split(url, "/")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}