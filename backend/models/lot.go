package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Lot struct {
	Id                  uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Type                string     `json:"type" gorm:"not null" binding:"required"`
	Units               uint       `json:"units" gorm:"not null" binding:"required"`
	StartCheckpointId   uuid.UUID  `json:"startCheckpointId" gorm:"not null"` // Added explicit foreign key field
	StartCheckpoint     Checkpoint `json:"startCheckpoint" gorm:"foreignKey:StartCheckpointId"`
	EndCheckpointId     uuid.UUID  `json:"endCheckpointId" gorm:"not null"` // Added explicit foreign key field
	EndCheckpoint       Checkpoint `json:"endCheckpoint" gorm:"foreignKey:EndCheckpointId"`
	TractorId           uuid.UUID  `json:"tractorId" gorm:""` // Changed from struct to UUID, added optional foreign key
	Tractor             Tractor    `json:"tractor" gorm:"foreignKey:TractorId"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CurrentCheckpointId uuid.UUID  `json:"current_checkpoint_id"` // Added current checkpoint foreign key
	CurrentCheckpoint   Checkpoint `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	OwnerId             uuid.UUID  `json:"owner_id" gorm:"not null"` // Changed from User to UUID
	Owner               User       `json:"owner" gorm:"foreignKey:OwnerId"`
	State               string     `json:"state" gorm:"not null"`
	MaxPriceByKm        float64    `json:"max_price_by_km" gorm:"not null"`
	TrafficManagerId    uuid.UUID  `json:"traffic_manager_id" gorm:"not null"` // Changed from User to UUID
	TrafficManager      User       `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TraderId            uuid.UUID  `json:"trader_id" gorm:"not null"` // Changed from User to UUID
	Trader              User       `json:"trader" gorm:"foreignKey:TraderId"`
}

func (lot *Lot) BeforeCreate(tx *gorm.DB) (err error) {
	lot.Id = uuid.New()
	return
}

func (lot *Lot) Save(db *gorm.DB) error {
	return db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Save(lot).Error
}
