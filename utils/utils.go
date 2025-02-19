package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	log.Println("[Info] Hashing password...")
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 15)
	if err != nil {
		log.Println("[Error] Error while trying to encrypt password")
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	log.Println("[Info] Comparing password hashes...")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	log.Println("[Info] Generating token...")
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("[Error] Failed to generate token: %v", err)
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	log.Printf("[Info] Created token: %s", token)
	return token
}
