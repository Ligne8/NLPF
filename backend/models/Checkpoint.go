package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Checkpoint struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name    string    `json:"name" gorm:"not null" binding:"required"`
	Country string    `json:"country" gorm:"not null" binding:"required"`
}

func (checkpoint *Checkpoint) BeforeCreate(tx *gorm.DB) (err error) {
	checkpoint.Id = uuid.New()
	return
}
