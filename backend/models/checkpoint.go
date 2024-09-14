package models

import "gorm.io/gorm"

type Checkpoint struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"not null" binding:"required"`
	Country string `json:"country" gorm:"not null" binding:"required"`
}

func (u *Checkpoint) Save(db *gorm.DB) error {
	return db.Save(u).Error
}