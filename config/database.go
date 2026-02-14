package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf(
	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASS"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PORT"),
)


	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
	log.Fatal("Failed to connect database ❌: ", err)
}


	DB = database
	fmt.Println("Database Connected ✅")
}
