package config

import (
	"time"
)

type Config struct {
	RedisURL      string
	ServerPort    string
	BaseURL       string
	SlugMinLength int
	SlugMaxLength int
}

func LoadConfig() *Config {
	return &Config{
		RedisURL:      "redis://default:CrRppNmTZsmlBvFYGFpAKujvahhNMyxb@maglev.proxy.rlwy.net:47225",
		ServerPort:    "8080",
		BaseURL:       "http://localhost:8080",
		SlugMinLength: 4,
		SlugMaxLength: 6,
	}
}

// Expiration times
const (
	URLExpiration = 24 * time.Hour * 30 // 30 days
	StatsTTL      = 24 * time.Hour * 90 // 90 days for stats
)
