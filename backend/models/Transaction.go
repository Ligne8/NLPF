package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionState string

const (
	TransactionStateIn  State = "in"
	TransactionStateOut State = "out"
)

type Transaction struct {
	Id                uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	TransactionType   TransactionState `json:"transaction_type" gorm:"not null" binding:"required"`
	CreateAt          time.Time        `json:"create_at" gorm:"autoCreateTime"`
	LotId             *uuid.UUID       `json:"lot_id" gorm:"not null"` // Foreign key for Lot
	Lot               *Lot             `json:"lot" gorm:"foreignKey:LotId"`
	TractorId         *uuid.UUID       `json:"tractor_id" gorm:"not null"` // Foreign key for Tractor
	Tractor           *Tractor         `json:"tractor" gorm:"foreignKey:TractorId;"`
	RouteId           *uuid.UUID       `json:"route_id" gorm:"not null"` // Foreign key for Route
	Route             *Route           `json:"route" gorm:"foreignKey:RouteId"`
	CheckpointId      *uuid.UUID       `json:"checkpoint_id" gorm:"not null"` // Foreign key for Checkpoint
	Checkpoint        *Checkpoint      `json:"checkpoint" gorm:"foreignKey:CheckpointId"`
	TrafficManager    *User            `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TrafficManagerId  *uuid.UUID       `json:"traffic_manager_id" gorm:"not null"` // Foreign key for Traffic Manager (User)
	RouteCheckpoint   *RouteCheckpoint `json:"route_checkpoint" gorm:"foreignKey:RouteCheckpointId"`
	RouteCheckpointId *uuid.UUID       `json:"route_checkpoint_id" gorm:"not null"` // Foreign key for RouteCheckpoint
}

func (transaction *Transaction) Save(db *gorm.DB) error {
	return db.Save(transaction).Error
}

func (transaction *Transaction) Update(db *gorm.DB) error {
	return db.Model(transaction).Updates(transaction).Error
}

func (transaction *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	transaction.Id = uuid.New()
	var simulation Simulation
	if err := tx.First(&simulation).Error; err != nil {
		return err
	}
	transaction.CreateAt = simulation.SimulationDate
	return
}

func (transaction *Transaction) CreateTransaction(db *gorm.DB, transactionType TransactionState, lotId uuid.UUID, tractorId uuid.UUID, routeId uuid.UUID, checkpointID uuid.UUID, traficManagerId uuid.UUID, routeCheckpointId uuid.UUID) error {
	var t = Transaction{
		TransactionType:   transactionType,
		LotId:             &lotId,
		TractorId:         &tractorId,
		RouteId:           &routeId,
		CheckpointId:      &checkpointID,
		TrafficManagerId:  &traficManagerId,
		RouteCheckpointId: &routeCheckpointId,
	}
	return db.Create(&t).Error
}

func (transaction *Transaction) FindByRouteId(db *gorm.DB, routeId uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	if err := db.Preload("Lot").Preload("Tractor").Preload("Route").Preload("Checkpoint").Find(&transactions, "route_id = ?", routeId).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (transaction *Transaction) FindByRouteIdAndCheckpointIdAndTractorId(db *gorm.DB, routeId uuid.UUID, checkpointId uuid.UUID, tractor_id uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	if err := db.Preload("Lot").Preload("Tractor").Preload("Route").Preload("Checkpoint").Find(&transactions, "route_id = ? AND checkpoint_id = ? AND tractor_id = ?", routeId, checkpointId, tractor_id).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (transaction *Transaction) ExecTransaction(db *gorm.DB) error {
	if transaction.TransactionType == TransactionState(TransactionStateIn) {
		transaction.Lot.State = StateInTransit;
		transaction.Lot.InTractor = true;
		transaction.Tractor.CurrentVolume = transaction.Tractor.CurrentVolume - transaction.Lot.Volume;
	} else {
		transaction.Lot.State = StateArchive;
		transaction.Lot.InTractor = false;
		transaction.Tractor.CurrentVolume = transaction.Tractor.CurrentVolume - transaction.Lot.Volume;
	}
	if err := transaction.Tractor.Save(db); err != nil {
		return err;
	}
	if err := transaction.Lot.Save(db); err != nil {
		return err;
	}
	return nil;
}
