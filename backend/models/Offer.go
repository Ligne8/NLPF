package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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
	offer.Id = uuid.New()
	return
}
