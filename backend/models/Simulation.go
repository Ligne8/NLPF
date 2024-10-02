package models

import (
	"github.com/google/uuid"
	"time"
)

type Simulation struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	SimulationDate time.Time `json:"simulation_date" gorm:"not null"`
}
