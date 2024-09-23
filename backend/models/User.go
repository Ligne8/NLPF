package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Firstname      string    `json:"firstname" gorm:"not null" binding:"required"`
	Lastname       string    `json:"lastname" gorm:"not null" binding:"required"`
	HashedPassword string    `json:"hashed_password" gorm:"not null"`
	Role           string    `json:"role" gorm:"not null"`
}
