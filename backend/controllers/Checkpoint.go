package controllers

import (
	"gorm.io/gorm"
)

type CheckpointController struct {
	Db *gorm.DB
}
