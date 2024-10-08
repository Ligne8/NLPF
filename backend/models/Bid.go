package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Bid struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:""`
	Bid       float64   `json:"bid" gorm:"not null"`
	OfferId   uuid.UUID `json:"offer_id" gorm:"not null"`
}

func (bid *Bid) BeforeCreate(tx *gorm.DB) (err error) {
	bid.Id = uuid.New()
	return
}
