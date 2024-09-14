package models

type Trip struct {
	Id 		 uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TractorId uint
	Tractor Tractor `json:"-"`
	Transactions []Transaction `json:"transactions"`
}