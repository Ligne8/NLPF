package database

import (
	"tms-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("api.sql"), &gorm.Config{})

	if err != nil {
			panic("Failed to connect database")
	}

	db.AutoMigrate(&models.Checkpoint{})
	db.AutoMigrate(&models.Lot{})
	db.AutoMigrate(&models.Tractor{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.Trip{})

	return db
}