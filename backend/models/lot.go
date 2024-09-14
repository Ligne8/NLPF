package models

import "gorm.io/gorm"

type Lot struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Type string `json:"type" gorm:"not null" binding:"required"`
	UnitVolume uint `json:"unit_volume" gorm:"not null" binding:"required"`
	Weight uint `json:"weight" gorm:"not null" binding:"required"`
	StartCheckpointId uint `json:"start_checkpoint_id" gorm:"not null" binding:"required"`
	EndCheckpointId uint `json:"end_checkpoint_id" gorm:"not null" binding:"required"`
	StartCheckpoint Checkpoint `json:"startCheckpoint"`
	EndCheckpoint Checkpoint `json:"endCheckpoint"`
	Tractor Tractor `json:"tractor"`
	TractorId *uint `json:"tractorId"`
}

func (u *Lot) Save(db *gorm.DB) error {
	return db.Save(u).Error
}

func (u *Lot) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}

func (u* Lot) ListAll (db *gorm.DB) ([]Lot, error) {
	var lots []Lot
	err := db.Find(&lots).Error
	return lots, err
}

func (u *Lot) GetByID(db *gorm.DB, id string) (*Lot, error) {
	var lot Lot
	err := db.First(&lot, id).Error
	return &lot, err
}

func (u *Lot) GetByTractorID(db *gorm.DB, id string) ([]Lot, error) {
	var lots []Lot
	err := db.Where("tractor_id = ?", id).Find(&lots).Error
	return lots, err
}

func (u *Lot) GetByStartCheckpointID(db *gorm.DB, id string) ([]Lot, error) {
	var lots []Lot
	err := db.Where("start_checkpoint_id = ?", id).Find(&lots).Error
	return lots, err
}

func (u *Lot) GetByEndCheckpointID(db *gorm.DB, id string) ([]Lot, error) {
	var lots []Lot
	err := db.Where("end_checkpoint_id = ?", id).Find(&lots).Error
	return lots, err
}