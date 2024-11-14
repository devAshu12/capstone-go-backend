package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var validate *validator.Validate

func ValidateRequest(req interface{}) error {
	return validate.Struct(req)
}

func LoadConfig() {
	validate = validator.New()
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
