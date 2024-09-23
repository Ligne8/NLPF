package models

import "github.com/google/uuid"

type Route struct {
	Id               uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name             string    `json:"name" gorm:"not null"`
	TrafficManagerId uint      `json:"traffic_manager_id" gorm:"not null"`
}

type RouteCheckpoint struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RouteId      Route      `json:"route_id"`
	CheckpointId Checkpoint `json:"checkpoint_id"`
	Position     uint       `json:"position"`
}
