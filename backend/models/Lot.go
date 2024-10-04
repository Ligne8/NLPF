package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type State string

const (
	StateAvailable State = "available"
	StatePending   State = "pending"
	StateInTransit State = "in_transit"
	StateArchive   State = "archive"
	StateOnMarket  State = "on_market"
	StateAtTrader  State = "at_trader"
)

type Lot struct {
	Id                  uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	ResourceType        ResourceType `json:"resource_type" gorm:"not null" binding:"required"`
	Volume              float64      `json:"volume" gorm:"not null" binding:"required"`
	StartCheckpointId   *uuid.UUID   `json:"start_checkpoint_id" gorm:""` // Changed to pointer to allow null values
	StartCheckpoint     *Checkpoint  `json:"start_checkpoint" gorm:"foreignKey:StartCheckpointId"`
	EndCheckpointId     *uuid.UUID   `json:"end_checkpoint_id" gorm:""` // Changed to pointer to allow null values
	EndCheckpoint       *Checkpoint  `json:"end_checkpoint" gorm:"foreignKey:EndCheckpointId"`
	TractorId           *uuid.UUID   `json:"tractor_id" gorm:""` // Changed to pointer to allow null values
	Tractor             *Tractor     `json:"tractor" gorm:"foreignKey:TractorId"`
	CreatedAt           time.Time    `json:"created_at" gorm:""`
	CurrentCheckpointId *uuid.UUID   `json:"current_checkpoint_id" gorm:""` // Changed to pointer to allow null values
	CurrentCheckpoint   *Checkpoint  `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	OwnerId             uuid.UUID    `json:"owner_id" gorm:"not null"` // Changed from User to UUID
	Owner               User         `json:"owner" gorm:"foreignKey:OwnerId"`
	State               State        `json:"state" gorm:"not null"`
	MaxPriceByKm        float64      `json:"max_price_by_km" gorm:"not null"`
	TrafficManagerId    *uuid.UUID   `json:"traffic_manager_id" gorm:""` // Added TrafficManagerId as pointer to allow null values
	TrafficManager      *User        `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TraderId            *uuid.UUID   `json:"trader_id" gorm:""` // Changed to pointer to allow null values
	Trader              *User        `json:"trader" gorm:"foreignKey:TraderId"`
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

func (lot *Lot) GetAllLots(db *gorm.DB) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) FindById(db *gorm.DB, lotId uuid.UUID) (Lot, error) {
	var foundLot Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").First(&foundLot, "id = ?", lotId).Error; err != nil {
		return Lot{}, err
	}
	return foundLot, nil
}

func (lot *Lot) GetLotsByState(db *gorm.DB, state State) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("state = ?", state).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) GetLotsByTrafficManager(db *gorm.DB, trafficManagerId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("traffic_manager_id = ?", trafficManagerId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) GetLotsByTrader(db *gorm.DB, traderId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("trader_id = ?", traderId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) GetLotsByOwner(db *gorm.DB, ownerId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Preload("CurrentCheckpoint").Preload("TrafficManager").Where("owner_id = ?", ownerId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) AssociateTraficManager(db *gorm.DB, trafficManagerId uuid.UUID) error {
	return db.Model(&lot).Update("traffic_manager_id", trafficManagerId).Error
}

func (lot *Lot) GetLotsByTractor(db *gorm.DB, tractorId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("tractor_id = ?", tractorId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) GetLotsByCurrentCheckpoint(db *gorm.DB, currentCheckpointId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("current_checkpoint_id = ?", currentCheckpointId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) GetLotsByStartCheckpoint(db *gorm.DB, startCheckpointId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("start_checkpoint_id = ?", startCheckpointId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) GetLotsByEndCheckpoint(db *gorm.DB, endCheckpointId uuid.UUID) ([]Lot, error) {
	var lots []Lot
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Where("end_checkpoint_id = ?", endCheckpointId).Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (lot *Lot) UpdateState(db *gorm.DB, state State) error {
	return db.Model(&lot).Update("state", state).Error
}

func (lot *Lot) UpdateStateByTractorId(db *gorm.DB, tractorId uuid.UUID, state State) error {
	return db.Model(&lot).Where("tractor_id = ?", tractorId).Update("state", state).Error
}
