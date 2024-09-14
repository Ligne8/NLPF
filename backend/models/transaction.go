package models

type Transaction struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TransactionType string `json:"transactionType" gorm:"not null" binding:"required"`
	Date string `json:"date" gorm:"not null" binding:"required"`
	Lot Lot `json:"lot"`
	LotId uint `json:"lotId" gorm:"not null" binding:"required"`
	Tractor Tractor `json:"tractor"`
	TractorId uint `json:"tractorId" gorm:"not null" binding:"required"`
	Trip Trip `json:"trip"`
	TripId uint `json:"tripId" gorm:"not null" binding:"required"`
	Checkpoint Checkpoint `json:"checkpoint"`
	CheckpointId uint `json:"checkpointId" gorm:"not null" binding:"required"`
}