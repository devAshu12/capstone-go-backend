package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	db_url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	// dsn := "host=localhost user=postgres password=asqwASQW12@ dbname=capstone port=5432 sslmode=disable TimeZone=UTC"
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// ----- replace &gorm.Config{} with gorm.Config{Logger: newLogger} to enable logging
	DB, err = gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to Database")
	}

	log.Print("Database connected successfully")
}
