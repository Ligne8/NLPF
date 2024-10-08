package database

import (
	"log"
	"tms-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	dsn := "host=localhost user=ligne8 password=secret dbname=tms_db port=5432 sslmode=disable TimeZone=Europe/Paris"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// AutoMigrate example for creating tables automatically
	err = db.AutoMigrate(&models.Checkpoint{}, &models.Lot{}, &models.Tractor{}, &models.User{}, &models.Route{}, &models.RouteCheckpoint{}, &models.Simulation{}, &models.Transaction{}, &models.Offer{}, &models.Bid{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}
	DB = db
	SeedDB(db)
	log.Println("Database connection established successfully.")
	return db
}

func SeedDB(db *gorm.DB) {
	models.CreateCheckpoints(db)
	SeedUsers(db)
	SeedTractors(db)
	SeedLots(DB)
}
