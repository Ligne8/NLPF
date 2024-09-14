package models

import "gorm.io/gorm"

type Tractor struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	MaxWeight uint `json:"maxWeight" gorm:"not null" binding:"required"`
	MaxVolumeUnit uint `json:"maxVolumeUnit" gorm:"not null" binding:"required"`
	CurrentCheckpointId uint `json:"currentCheckpointId" gorm:"not null" binding:"required"`
	CurrentCheckpoint Checkpoint `json:"currentCheckpoint"`
	Lots []Lot `json:"lots"`
	TripId uint 
	
}


func (u *Tractor) Save(db *gorm.DB) error {
	return db.Save(u).Error
}

func (u *Tractor) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}


