package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tms-backend/models"
)

var DB *gorm.DB

func InitDb() *gorm.DB {
	dsn := "host=localhost user=ligne8 password=secret dbname=tms_db port=5432 sslmode=disable TimeZone=Europe/Paris"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(
		&models.Checkpoint{},
		&models.Lot{},
		&models.Tractor{},
		&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}

	DB = db
	log.Println("Database connection established successfully.")
	return db
}
