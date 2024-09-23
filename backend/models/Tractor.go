package models

import (
	"github.com/google/uuid"
	"time"
)

type Tractor struct {
	Id                uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MaxUnits          uint       `json:"max_units" gorm:"not null"`
	CurrentUnits      uint       `json:"current_units" gorm:"not null"`
	CurrentCheckpoint Checkpoint `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	State             string     `json:"state" gorm:"not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	OwnerId           User       `json:"owner_id" gorm:"foreignKey:OwnerId"`
	MinPriceByKm      uint       `json:"min_price_by_km" gorm:"not null"`
	TrafficManagerId  User       `json:"traffic_manager_id" gorm:"foreignKey:TrafficManagerId"`
	TraderId          User       `json:"trader_id" gorm:"not null"`
	RouteId           Route      `json:"route_id"`
}
