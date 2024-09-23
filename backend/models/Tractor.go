package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Tractor struct {
	Id                  uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ResourceType        string     `json:"resource_type" gorm:"not null" binding:"required"`
	MaxVolume           float64    `json:"max_units" gorm:"not null"`
	CurrentVolume       float64    `json:"current_units" gorm:"not null"`
	CurrentCheckpointId uuid.UUID  `json:"current_checkpoint_id" gorm:"type:uuid"` // Foreign key for Checkpoint
	CurrentCheckpoint   Checkpoint `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	State               string     `json:"state" gorm:"not null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"autoCreateTime"`
	OwnerId             uuid.UUID  `json:"owner_id" gorm:"type:uuid"` // Foreign key for User
	Owner               User       `json:"owner" gorm:"foreignKey:OwnerId"`
	MinPriceByKm        uint       `json:"min_price_by_km" gorm:"not null"`
	TrafficManagerId    uuid.UUID  `json:"traffic_manager_id" gorm:"type:uuid"` // Foreign key for User
	TrafficManager      User       `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TraderId            uuid.UUID  `json:"trader_id" gorm:"type:uuid"` // Foreign key for User
	Trader              User       `json:"trader" gorm:"foreignKey:TraderId"`
	RouteId             uuid.UUID  `json:"route_id" gorm:"type:uuid"` // Foreign key for Route
	Route               Route      `json:"route" gorm:"foreignKey:RouteId"`
}

func (tractor *Tractor) BeforeCreate(tx *gorm.DB) (err error) {
	tractor.Id = uuid.New()
	return
}
