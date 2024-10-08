package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bid struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:""`
	Bid       float64   `json:"bid" gorm:"not null"`
	OfferId   uuid.UUID `json:"offer_id" gorm:"type:uuid;not null"`
	Offer		 Offer     `json:"offer" gorm:"foreignKey:OfferId;references:Id"`
	State 	 string    `json:"state" gorm:"not null"`
	Volume    float64   `json:"volume" gorm:""`
}

func (bid *Bid) BeforeCreate(tx *gorm.DB) (err error) {
	if bid.Id == uuid.Nil {
		bid.Id = uuid.New()
	}
	return
}
