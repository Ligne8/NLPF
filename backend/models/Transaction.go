package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	Id              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	TransactionType string     `json:"transactionType" gorm:"not null" binding:"required"`
	Date            time.Time  `json:"date" gorm:"autoCreateTime"`
	LotId           uuid.UUID  `json:"lot_id" gorm:"not null"` // Foreign key for Lot
	Lot             Lot        `json:"lot" gorm:"foreignKey:LotId"`
	TractorId       uuid.UUID  `json:"tractor_id" gorm:"not null"` // Foreign key for Tractor
	Tractor         Tractor    `json:"tractor" gorm:"foreignKey:TractorId"`
	RouteId         uuid.UUID  `json:"route_id" gorm:"not null"` // Foreign key for Route
	Route           Route      `json:"route" gorm:"foreignKey:RouteId"`
	CheckpointId    uuid.UUID  `json:"checkpoint_id" gorm:"not null"` // Foreign key for Checkpoint
	Checkpoint      Checkpoint `json:"checkpoint" gorm:"foreignKey:CheckpointId"`
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	transaction.Id = uuid.New()
	return
}
