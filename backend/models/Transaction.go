package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionState string

const (
	TransactionStateIn State = "in"
	TransactionStateOut State = "out"
)

type Transaction struct {
	Id              uuid.UUID             `json:"id" gorm:"type:uuid;primaryKey"`
	TransactionType TransactionState      `json:"transaction_type" gorm:"not null" binding:"required"`
	CreateAt            time.Time  						`json:"create_at" gorm:"autoCreateTime"`
	LotId           uuid.UUID  						`json:"lot_id" gorm:"not null"` // Foreign key for Lot
	Lot             Lot        						`json:"lot" gorm:"foreignKey:LotId"`
	TractorId       uuid.UUID  						`json:"tractor_id" gorm:"not null"` // Foreign key for Tractor
	Tractor         Tractor    						`json:"tractor" gorm:"foreignKey:TractorId"`
	RouteId         uuid.UUID  						`json:"route_id" gorm:"not null"` // Foreign key for Route
	Route           Route      						`json:"route" gorm:"foreignKey:RouteId"`
	CheckpointId    uuid.UUID  						`json:"checkpoint_id" gorm:"not null"` // Foreign key for Checkpoint
	Checkpoint      Checkpoint 						`json:"checkpoint" gorm:"foreignKey:CheckpointId"`
	TrafficManager User			 						`json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TrafficManagerId uuid.UUID 						`json:"traffic_manager_id" gorm:"not null"` // Foreign key for Traffic Manager (User)
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	transaction.Id = uuid.New()
	return
}

func (transaction *Transaction) FindByRouteId(db *gorm.DB, routeId uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	if err := db.Preload("Lot").Preload("Tractor").Preload("Route").Preload("Checkpoint").Find(&transactions, "route_id = ?", routeId).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
