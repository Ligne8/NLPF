package models

import (
	"github.com/google/uuid"
)

type Checkpoint struct {
	Id      uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name    string    `json:"name" gorm:"not null" binding:"required"`
	Country string    `json:"country" gorm:"not null" binding:"required"`
}
