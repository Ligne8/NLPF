package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ResourceType string

const (
	ResourceTypeBulk   ResourceType = "bulk"
	ResourceTypeSolid  ResourceType = "solid"
	ResourceTypeLiquid ResourceType = "liquid"
)

type Tractor struct {
	Id                  uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	ResourceType        ResourceType `json:"resource_type" gorm:"type:varchar(10)" binding:"required"`
	MaxVolume           float64      `json:"max_units" gorm:"not null"`
	CurrentVolume       float64      `json:"current_units" gorm:"not null"`
	CurrentCheckpointId uuid.UUID    `json:"current_checkpoint_id" gorm:"type:uuid"` // Foreign key for Checkpoint
	CurrentCheckpoint   Checkpoint   `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	State               State        `json:"state" gorm:"not null"`
	CreatedAt           time.Time    `json:"created_at" gorm:"autoCreateTime"`
	OwnerId             uuid.UUID    `json:"owner_id" gorm:"type:uuid"` // Foreign key for User
	Owner               User         `json:"owner" gorm:"foreignKey:OwnerId"`
	MinPriceByKm        uint         `json:"min_price_by_km" gorm:"not null"`
	TrafficManagerId    uuid.UUID    `json:"traffic_manager_id" gorm:"type:uuid"` // Foreign key for User
	TrafficManager      User         `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TraderId            uuid.UUID    `json:"trader_id" gorm:"type:uuid"` // Foreign key for User
	Trader              User         `json:"trader" gorm:"foreignKey:TraderId"`
	RouteId             uuid.UUID    `json:"route_id" gorm:"type:uuid"` // Foreign key for Route
	Route               Route        `json:"route" gorm:"foreignKey:RouteId"`
}

func (tractor *Tractor) BeforeCreate(tx *gorm.DB) (err error) {
	validTypes := map[ResourceType]bool{
		ResourceTypeBulk:   true,
		ResourceTypeSolid:  true,
		ResourceTypeLiquid: true,
	}
	if !validTypes[tractor.ResourceType] {
		return errors.New("invalid resource type")
	}

	validState := map[State]bool{
		StateAvailable: true,
		StateArchive:   true,
		StateAtTrader:  true,
		StateInTransit: true,
		StateOnMarket:  true,
	}

	if !validState[tractor.State] {
		return errors.New("invalid valid state")
	}

	tractor.Id = uuid.New()
	return
}

func (tractor *Tractor) GetByOwnerId(db *gorm.DB, ownerId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Where("owner_id = ?", ownerId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetByTrafficManagerId(db *gorm.DB, trafficManagerId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Where("traffic_manager_id = ?", trafficManagerId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetByState(db *gorm.DB, state string) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Where("state = ?", state).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetByRouteId(db *gorm.DB, routeId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Where("route_id = ?", routeId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}
