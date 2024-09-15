package models

import "gorm.io/gorm"

type Lot struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Type string `json:"type" gorm:"not null" binding:"required"`
	UnitVolume uint `json:"unit_volume" gorm:"not null" binding:"required"`
	Weight uint `json:"weight" gorm:"not null" binding:"required"`
	StartCheckpointId uint `json:"start_checkpoint_id" gorm:"not null" binding:"required"`
	EndCheckpointId uint `json:"end_checkpoint_id" gorm:"not null" binding:"required"`
	StartCheckpoint Checkpoint `json:"startCheckpoint" gorm:"foreignKey:start_checkpoint_id"`
	EndCheckpoint Checkpoint `json:"endCheckpoint" gorm:"foreignKey:end_checkpoint_id"`
	Tractor Tractor `json:"tractor"`
	TractorId *uint `json:"tractorId"`
}

func (u *Lot) Save(db *gorm.DB) error {
	return db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Save(u).Error
}

func (u *Lot) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}

func (u* Lot) ListAll (db *gorm.DB) ([]Lot, error) {
	var lots []Lot
	err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").Find(&lots).Error
	return lots, err
}

func (u *Lot) GetByID(db *gorm.DB, id string) (*Lot, error) {
	var lot Lot
	err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("Tractor").First(&lot, id).Error
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

func (u *Lot) AssociateTractorToLot(db *gorm.DB, lotId string, tractorId string) error {
	return db.Model(u).Where("id = ?", lotId).Update("tractor_id", tractorId).Error
}