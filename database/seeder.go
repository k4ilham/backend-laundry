package database

import (
	"laundry-backend/models"
	"laundry-backend/utils"
	"log"
)

func Seed() {
	// Seed Admin
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hash, _ := utils.HashPassword("admin123")
		DB.Create(&models.User{
			Name:     "Admin",
			Email:    "admin@laundry.com",
			Password: hash,
			Role:     "admin",
		})
		log.Println("Seeded Admin User")
	}

	// Seed Services
	var serviceCount int64
	DB.Model(&models.Service{}).Count(&serviceCount)
	if serviceCount == 0 {
		services := []models.Service{
			{Name: "Cuci Komplit", Description: "Cuci setrika rapi", Unit: "kg", Price: 7000},
			{Name: "Cuci Kering", Description: "Hanya cuci dan kering", Unit: "kg", Price: 5000},
			{Name: "Setrika Saja", Description: "Hanya jasa setrika", Unit: "kg", Price: 4000},
			{Name: "Cuci Bedcover", Description: "Cuci bedcover ukuran besar", Unit: "pcs", Price: 25000},
		}
		DB.Create(&services)
		log.Println("Seeded Default Services")
	}
}
