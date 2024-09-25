package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type State string

const (
	StateAvailable State = "available"
	StateInTransit State = "in_transit"
	StateArchive   State = "archive"
	StateOnMarket  State = "on_market"
	StateAtTrader  State = "at_trader"
)

type Lot struct {
	Id                  uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	ResourceType        ResourceType `json:"resource_type" gorm:"not null" binding:"required"`
	Volume              float64      `json:"volume" gorm:"not null" binding:"required"`
	StartCheckpointId   uuid.UUID    `json:"start_checkpoint_id" gorm:"not null"` // Added explicit foreign key field
	StartCheckpoint     Checkpoint   `json:"start_checkpoint" gorm:"foreignKey:StartCheckpointId"`
	EndCheckpointId     uuid.UUID    `json:"end_checkpointId" gorm:"not null"` // Added explicit foreign key field
	EndCheckpoint       Checkpoint   `json:"end_checkpoint" gorm:"foreignKey:EndCheckpointId"`
	TractorId           uuid.UUID    `json:"tractor_id" gorm:""` // Changed from struct to UUID, added optional foreign key
	Tractor             Tractor      `json:"tractor" gorm:"foreignKey:TractorId"`
	CreatedAt           time.Time    `json:"created_at" gorm:"autoCreateTime"`
	CurrentCheckpointId uuid.UUID    `json:"current_checkpoint_id"` // Added current checkpoint foreign key
	CurrentCheckpoint   Checkpoint   `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	OwnerId             uuid.UUID    `json:"owner_id" gorm:"not null"` // Changed from User to UUID
	Owner               User         `json:"owner" gorm:"foreignKey:OwnerId"`
	State               State        `json:"state" gorm:"not null"`
	MaxPriceByKm        float64      `json:"max_price_by_km" gorm:"not null"`
	TrafficManagerId    uuid.UUID    `json:"traffic_manager_id" gorm:"not null"` // Changed from User to UUID
	TrafficManager      User         `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TraderId            uuid.UUID    `json:"trader_id" gorm:"not null"` // Changed from User to UUID
	Trader              User         `json:"trader" gorm:"foreignKey:TraderId"`
}

func (lot *Lot) BeforeCreate(tx *gorm.DB) (err error) {
	validTypes := map[ResourceType]bool{
		ResourceTypeBulk:   true,
		ResourceTypeSolid:  true,
		ResourceTypeLiquid: true,
	}
	if !validTypes[lot.ResourceType] {
		return errors.New("invalid resource type")
	}

	validState := map[State]bool{
		StateAvailable: true,
		StateArchive:   true,
		StateAtTrader:  true,
		StateInTransit: true,
		StateOnMarket:  true,
	}

	if !validState[lot.State] {
		return errors.New("invalid valid state")
	}

	lot.Id = uuid.New()
	return
}

func (lot *Lot) Save(db *gorm.DB) error {
	return db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Save(lot).Error
}
