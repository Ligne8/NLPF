package database

import (
	"log"
	"tms-backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) []models.User {
	if db == nil {
		log.Fatalf("db is nil")
	}
	users := []models.User{
		{Id: uuid.New(), Username: "admin", Password: "test", Role: models.Role("admin")},
		{Id: uuid.New(), Username: "tm", Password: "test", Role: models.Role("traffic_manager")},
		{Id: uuid.New(), Username: "client", Password: "test", Role: models.Role("client")},
		{Id: uuid.New(), Username: "trader", Password: "test", Role: models.Role("trader")},
	}

	for _, user := range users {
		var existingUser models.User
		if err := db.Where("username = ? AND role = ?", user.Username, user.Role).First(&existingUser).Error; err == nil {
			log.Printf("User already exists: %s, %s", user.Username, user.Role)
			continue
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Username, err)
			continue
		}
		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
		} else {
			log.Printf("User created: %s, %s", user.Username, user.Role)
		}
	}
	log.Printf("Users: %v", users)
	return users
}

func SeedTractors(db *gorm.DB, checkpoints []models.Checkpoint, users []models.User) {
	clientId := users[2].Id
	log.Printf("Client ID: %v", users[2])
	tractors := []models.Tractor{
		{Id: uuid.New(), Name: "Tractor A", ResourceType: "bulk", MaxVolume: 200, CurrentVolume: 2, State: "available", CurrentCheckpointId: &checkpoints[5].Id, StartCheckpointId: &checkpoints[3].Id, EndCheckpointId: &checkpoints[1].Id, MinPriceByKm: 0.8, OwnerId: clientId},
		{Id: uuid.New(), Name: "Tractor B", ResourceType: "bulk", MaxVolume: 100, CurrentVolume: 50, State: "pending", CurrentCheckpointId: &checkpoints[7].Id, StartCheckpointId: &checkpoints[2].Id, EndCheckpointId: &checkpoints[4].Id, MinPriceByKm: 8.0, OwnerId: clientId},
		{Id: uuid.New(), Name: "Tractor C", ResourceType: "solid", MaxVolume: 100, CurrentVolume: 20, State: "in_transit", CurrentCheckpointId: &checkpoints[0].Id, StartCheckpointId: &checkpoints[6].Id, EndCheckpointId: &checkpoints[8].Id, MinPriceByKm: 2.0, OwnerId: clientId},
		{Id: uuid.New(), Name: "Tractor D", ResourceType: "solid", MaxVolume: 100, CurrentVolume: 100, State: "archive", CurrentCheckpointId: &checkpoints[9].Id, StartCheckpointId: &checkpoints[11].Id, EndCheckpointId: &checkpoints[10].Id, MinPriceByKm: 1.5, OwnerId: clientId},
		{Id: uuid.New(), Name: "Tractor E", ResourceType: "liquide", MaxVolume: 100, CurrentVolume: 50, State: "on_market", CurrentCheckpointId: &checkpoints[12].Id, StartCheckpointId: &checkpoints[14].Id, EndCheckpointId: &checkpoints[13].Id, MinPriceByKm: 9.0, OwnerId: clientId},
		{Id: uuid.New(), Name: "Tractor F", ResourceType: "liquide", MaxVolume: 100, CurrentVolume: 50, State: "at_trader", CurrentCheckpointId: &checkpoints[16].Id, StartCheckpointId: &checkpoints[15].Id, EndCheckpointId: &checkpoints[17].Id, MinPriceByKm: 11.0, OwnerId: clientId},
	}

	for _, tractor := range tractors {
		var existingTractor models.Tractor
		log.Printf("Tractor: %v", tractor)
		log.Printf("Tractor.name: %v", tractor.Name)
		log.Printf("Tractor.resource_type: %v", tractor.ResourceType)
		log.Printf("Tractor.state: %v", tractor.State)
		log.Printf("Tractor.current_checkpoint_id: %v", tractor.CurrentCheckpointId)
		log.Printf("Tractor.start_checkpoint_id: %v", tractor.StartCheckpointId)
		log.Printf("Tractor.end_checkpoint_id: %v", tractor.EndCheckpointId)
		log.Printf("Tractor.owner_id: %v", tractor.OwnerId)
		log.Printf("Tractor.min_price_by_km: %v", tractor.MinPriceByKm)

		if err := db.Where("name = ? AND resource_type = ? AND state = ? AND current_checkpoint_id = ? AND start_checkpoint_id = ? AND end_checkpoint_id = ?",
			tractor.Name, tractor.ResourceType, tractor.State, tractor.CurrentCheckpointId, tractor.StartCheckpointId, tractor.EndCheckpointId).First(&existingTractor).Error; err == nil {
			log.Printf("Tractor already exists: %s, %s", tractor.Name, tractor.ResourceType)
			continue
		}

		if err := db.Create(&tractor).Error; err != nil {
			log.Fatalf("could not seed tractors: %v", err)
		} else {
			log.Printf("Tractor created: %s, %s", tractor.Name, tractor.ResourceType)
		}
	}
}

/*
func SeedLots(db *gorm.DB) {
	lots := []models.Lot{
		{ResourceType: "lot", Volume: 100},
		{ResourceType: "lot", Volume: 200},
	}

	for _, lot := range lots {
		if err := db.Create(&lot).Error; err != nil {
			log.Fatalf("could not seed lots: %v", err)
		}
	}
}

func SeedRoutes(db *gorm.DB) {
	routes := []models.Route{
		{Name: "Route 1", TrafficManagerId: uuid.New()},
		{Name: "Route 2", TrafficManagerId: uuid.New()},
	}

	for _, route := range routes {
		if err := db.Create(&route).Error; err != nil {
			log.Fatalf("could not seed routes: %v", err)
		}
	}
}
*/
