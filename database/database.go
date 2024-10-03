package database

import (
	"appartments-server/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDB() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	dsn := os.Getenv("DB_CONNECTION_URL")

	db, dbError := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if dbError != nil {
		log.Panic("Error connecting to database")
	}

	DB = db
	return DB
}

func performMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
	)
}

func InitializeDB() *gorm.DB {
	db := connectToDB()
	performMigrations(db)
	return db
}
