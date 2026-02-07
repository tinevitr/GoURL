package utils

import (
	"math/rand"
	"time"
	"gourl/config"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateSlug() string {
	cfg := config.LoadConfig()
	length := seededRand.Intn(cfg.SlugMaxLength-cfg.SlugMinLength+1) + cfg.SlugMinLength
	
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func IsValidURL(url string) bool {
	// You can add more validation logic here
	// For now, just check if it's not empty
	return len(url) > 0
}

func IsValidSlug(slug string) bool {
	if len(slug) < 3 || len(slug) > 20 {
		return false
	}
	
	// Check if slug contains only allowed characters
	for _, c := range slug {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}
