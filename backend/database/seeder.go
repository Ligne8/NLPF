package database

import (
	"log"
	"tms-backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	if db == nil {
		log.Fatalf("db is nil")
	}
	users := []models.User{
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
}

func SeedTractors(db *gorm.DB) {
	var checkpoints []models.Checkpoint
	db.Find(&checkpoints)
	var client models.User
	db.Find(&client, "role = ?", "client")
	clientId := client.Id
	var traffic_manger models.User
	db.Find(&traffic_manger, "role = ?", "traffic_manager")
	trafficManagerId := traffic_manger.Id
	tractors := []models.Tractor{
		{Id: uuid.New(), Name: "Tractor A", ResourceType: "bulk", MaxVolume: 200, CurrentVolume: 2, State: "available", CurrentCheckpointId: &checkpoints[5].Id, StartCheckpointId: &checkpoints[3].Id, EndCheckpointId: &checkpoints[1].Id, MinPriceByKm: 0.8, OwnerId: clientId},
		{Id: uuid.New(), Name: "Tractor B", ResourceType: "bulk", MaxVolume: 100, CurrentVolume: 50, State: "pending", CurrentCheckpointId: &checkpoints[2].Id, StartCheckpointId: &checkpoints[2].Id, EndCheckpointId: &checkpoints[4].Id, MinPriceByKm: 8.0, OwnerId: clientId, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), Name: "Tractor C", ResourceType: "solid", MaxVolume: 100, CurrentVolume: 20, State: "in_transit", CurrentCheckpointId: &checkpoints[0].Id, StartCheckpointId: &checkpoints[6].Id, EndCheckpointId: &checkpoints[8].Id, MinPriceByKm: 2.0, OwnerId: clientId, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), Name: "Tractor D", ResourceType: "solid", MaxVolume: 100, CurrentVolume: 100, State: "archive", CurrentCheckpointId: &checkpoints[9].Id, StartCheckpointId: &checkpoints[11].Id, EndCheckpointId: &checkpoints[10].Id, MinPriceByKm: 1.5, OwnerId: clientId, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), Name: "Tractor E", ResourceType: "liquid", MaxVolume: 100, CurrentVolume: 50, State: "on_market", CurrentCheckpointId: &checkpoints[12].Id, StartCheckpointId: &checkpoints[14].Id, EndCheckpointId: &checkpoints[13].Id, MinPriceByKm: 9.0, OwnerId: clientId, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), Name: "Tractor F", ResourceType: "liquid", MaxVolume: 100, CurrentVolume: 50, State: "at_trader", CurrentCheckpointId: &checkpoints[16].Id, StartCheckpointId: &checkpoints[15].Id, EndCheckpointId: &checkpoints[17].Id, MinPriceByKm: 11.0, OwnerId: clientId, TrafficManagerId: &trafficManagerId},
	}

	for _, tractor := range tractors {
		var existingTractor models.Tractor

		if err := db.Where("name = ? AND resource_type = ? AND current_checkpoint_id = ? AND start_checkpoint_id = ? AND end_checkpoint_id = ? AND owner_id = ?",
			tractor.Name, tractor.ResourceType, tractor.CurrentCheckpointId, tractor.StartCheckpointId, tractor.EndCheckpointId, tractor.OwnerId).First(&existingTractor).Error; err == nil {
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

func SeedLots(db *gorm.DB) {
	var checkpoints []models.Checkpoint
	db.Find(&checkpoints)
	var client models.User
	db.Find(&client, "role = ?", "client")
	clientId := client.Id
	var traffic_manger models.User
	db.Find(&traffic_manger, "role = ?", "traffic_manager")
	trafficManagerId := traffic_manger.Id
	lots := []models.Lot{
		{Id: uuid.New(), ResourceType: "bulk", Volume: 200, StartCheckpointId: &checkpoints[5].Id, EndCheckpointId: &checkpoints[1].Id, CurrentCheckpointId: &checkpoints[5].Id, OwnerId: clientId, State: "available", MaxPriceByKm: 0.8},
		{Id: uuid.New(), ResourceType: "bulk", Volume: 100, StartCheckpointId: &checkpoints[2].Id, EndCheckpointId: &checkpoints[4].Id, CurrentCheckpointId: &checkpoints[2].Id, OwnerId: clientId, State: "pending", MaxPriceByKm: 8.0, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), ResourceType: "solid", Volume: 100, StartCheckpointId: &checkpoints[6].Id, EndCheckpointId: &checkpoints[8].Id, CurrentCheckpointId: &checkpoints[6].Id, OwnerId: clientId, State: "in_transit", MaxPriceByKm: 2.0, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), ResourceType: "solid", Volume: 100, StartCheckpointId: &checkpoints[11].Id, EndCheckpointId: &checkpoints[10].Id, CurrentCheckpointId: &checkpoints[11].Id, OwnerId: clientId, State: "archive", MaxPriceByKm: 1.5, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), ResourceType: "liquid", Volume: 100, StartCheckpointId: &checkpoints[14].Id, EndCheckpointId: &checkpoints[13].Id, CurrentCheckpointId: &checkpoints[14].Id, OwnerId: clientId, State: "on_market", MaxPriceByKm: 9.0, TrafficManagerId: &trafficManagerId},
		{Id: uuid.New(), ResourceType: "liquid", Volume: 100, StartCheckpointId: &checkpoints[15].Id, EndCheckpointId: &checkpoints[17].Id, CurrentCheckpointId: &checkpoints[20].Id, OwnerId: clientId, State: "at_trader", MaxPriceByKm: 11.0, TrafficManagerId: &trafficManagerId},
	}

	for _, lot := range lots {
		var existingLot models.Lot

		if err := db.Where("resource_type = ? AND volume = ? AND start_checkpoint_id = ? AND end_checkpoint_id = ? AND owner_id = ?",
			lot.ResourceType, lot.Volume, lot.StartCheckpointId, lot.EndCheckpointId, lot.OwnerId).First(&existingLot).Error; err == nil {
			log.Printf("Lot already exists: %s, %s", lot.ResourceType, lot.Volume)
			continue
		}

		if err := db.Create(&lot).Error; err != nil {
			log.Fatalf("could not seed lots: %v", err)
		} else {
			log.Printf("Lot created: %s, %s", lot.ResourceType, lot.Volume)
		}
	}
}
