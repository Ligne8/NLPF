package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Offer struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	LimitDate time.Time  `json:"limit_date" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	TractorId *uuid.UUID `json:"tractor_id" gorm:""` // Changed to pointer to allow null values
	Tractor   *Tractor   `json:"tractor" gorm:"foreignKey:TractorId"`
	LotId     *uuid.UUID `json:"lot_id" gorm:""` // Changed to pointer to allow null values
	Lot       *Lot       `json:"lot" gorm:"foreignKey:LotId"`
}

func (offer *Offer) BeforeCreate(tx *gorm.DB) (err error) {
	if offer.Id == uuid.Nil {
		offer.Id = uuid.New()
	}
	var simulation Simulation
	if err := tx.First(&simulation).Error; err != nil {
		return err
	}
	if offer.CreatedAt.IsZero() {
		offer.CreatedAt = simulation.SimulationDate
	}
	return
}

func (offer *Offer) CreateOfferLot(db *gorm.DB, limitDate time.Time, lotId uuid.UUID) (uuid.UUID, error) {
	var o = Offer{
		LimitDate: limitDate,
		LotId:     &lotId,
	}
	if err := db.Create(&o).Error; err != nil {
		return uuid.Nil, err
	}
	return o.Id, nil
}

func (offer *Offer) CreateOfferTractor(db *gorm.DB, limitDate time.Time, tractorId uuid.UUID) (uuid.UUID, error) {
	var o = Offer{
		LimitDate: limitDate,
		TractorId: &tractorId,
	}
	if err := db.Create(&o).Error; err != nil {
		return uuid.Nil, err
	}
	return o.Id, nil
}
