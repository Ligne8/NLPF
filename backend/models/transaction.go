package models

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	Id              uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	TransactionType string     `json:"transactionType" gorm:"not null" binding:"required"`
	Date            time.Time  `json:"date" gorm:"autoCreateTime"`
	LotId           Lot        `json:"lot_id" gorm:"not null" binding:"required"`
	TractorId       uint       `json:"tractorId" gorm:"not null" binding:"required"`
	RouteId         Route      `json:"route_id" gorm:"not null" binding:"required"`
	CheckpointId    Checkpoint `json:"checkpointId" gorm:"not null" binding:"required"`
}
