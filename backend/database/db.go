package database

import (
	"tms-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("api.sql"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to the database")
	}

	// Migration des modèles
	err = db.AutoMigrate(
		&models.Checkpoint{},
		&models.Lot{},
		&models.Tractor{},
		&models.Transaction{},
		&models.Route{},
		&models.RouteCheckpoint{},
		&models.User{},
	)

	if err != nil {
		panic("Failed to auto-migrate models")
	}

	return db
}
