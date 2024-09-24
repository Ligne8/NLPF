package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleTrader         Role = "trader"
	RoleTrafficManager Role = "traffic_manager"
	RoleClient         Role = "client"
	RoleAdmin          Role = "admin"
)

type User struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Firstname      string    `json:"firstname" gorm:"not null" binding:"required"`
	Lastname       string    `json:"lastname" gorm:"not null" binding:"required"`
	HashedPassword string    `json:"hashed_password" gorm:"not null"`
	Role           Role      `json:"role" gorm:"not null"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	validUser := map[Role]bool{
		RoleTrader:         true,
		RoleTrafficManager: true,
		RoleClient:         true,
		RoleAdmin:          true,
	}
	if !validUser[user.Role] {
		return errors.New("invalid user role")
	}

	user.Id = uuid.New()
	return
}
