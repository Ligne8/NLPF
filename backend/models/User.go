package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Firstname      string    `json:"firstname" gorm:"not null" binding:"required"`
	Lastname       string    `json:"lastname" gorm:"not null" binding:"required"`
	HashedPassword string    `json:"hashed_password" gorm:"not null"`
	Role           string    `json:"role" gorm:"not null"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Id = uuid.New()
	return
}
