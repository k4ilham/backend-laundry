package database

import (
	"fmt"
	"log"

	"laundry-backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Get("DB_HOST"),
		config.Get("DB_USER"),
		config.Get("DB_PASSWORD"),
		config.Get("DB_NAME"),
		config.Get("DB_PORT"),
		config.Get("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	log.Println("Database connection established")
}
