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
	Id       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Username string    `json:"username" gorm:"not null" binding:"required"`
	Password string    `json:"password" gorm:"not null"`
	Role     Role      `json:"role" gorm:"not null"`
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

func (user *User) GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (user *User) FindById(db *gorm.DB, userId uuid.UUID) (User, error) {
	var foundUser User
	if err := db.First(&foundUser, "id = ?", userId).Error; err != nil {
		return User{}, err
	}
	return foundUser, nil
}

func (user *User) getRole(db *gorm.DB, userId uuid.UUID) (Role, error) {
	db.First(&user, "id = ?", userId)
	return user.Role, nil
}
