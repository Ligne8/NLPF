package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Lot struct {
	Id                uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Type              string     `json:"type" gorm:"not null" binding:"required"`
	Units             uint       `json:"units" gorm:"not null" binding:"required"`
	StartCheckpoint   Checkpoint `json:"startCheckpoint" gorm:"foreignKey:StartCheckpointId"`
	EndCheckpoint     Checkpoint `json:"endCheckpoint" gorm:"foreignKey:EndCheckpointId"`
	TractorId         Tractor    `json:"tractorId"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CurrentCheckpoint Checkpoint `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	OwnerId           User       `json:"owner_id" gorm:"not null"`
	State             string     `json:"state" gorm:"not null"`
	MaxPriceByKm      float64    `json:"max_price_by_km" gorm:"not null"`
	TrafficManagerId  User       `json:"traffic_manager_id" gorm:"not null"`
	TraderId          User       `json:"trader_id" gorm:"not null"`
}

func (u *Lot) Save(db *gorm.DB) error {
	return db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Save(u).Error
}
