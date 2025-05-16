package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ServerPort     string
	JWTSecret      string
	JWTExpiration  time.Duration
	ResendAPIKey   string
	SenderEmail    string
	AllowedOrigins string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expiration, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		expiration = 24 * time.Hour // Default to 24 hours
	}

	return &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		ServerPort:     os.Getenv("SERVER_PORT"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiration:  expiration,
		ResendAPIKey:   os.Getenv("RESEND_API_KEY"),
		SenderEmail:    os.Getenv("SENDER_EMAIL"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}, nil
}
