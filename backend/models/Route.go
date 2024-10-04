package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	Id               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name             string    `json:"name" gorm:"not null"`
	TrafficManagerId uuid.UUID `json:"traffic_manager_id" gorm:"type:uuid;not null"` // Foreign key for Traffic Manager (User)
	TrafficManager   User      `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
}

func (route *Route) BeforeCreate(tx *gorm.DB) (err error) {
	route.Id = uuid.New()
	return
}

type RouteCheckpoint struct {
	Id           uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	RouteId      uuid.UUID  `json:"route_id" gorm:"not null"` // Foreign key for Route
	Route        Route      `json:"route" gorm:"foreignKey:RouteId"`
	CheckpointId uuid.UUID  `json:"checkpoint_id" gorm:"not null"` // Foreign key for Checkpoint
	Checkpoint   Checkpoint `json:"checkpoint" gorm:"foreignKey:CheckpointId"`
	Position     uint        `json:"position" gorm:"not null"`
}

func (route *Route) GetRoutesByTrafficManagerId(db *gorm.DB, trafficManagerId uuid.UUID) ([]Route, error) {
	var routes []Route
	if err := db.Find(&routes, "traffic_manager_id = ?", trafficManagerId).Error; err != nil {
		return nil, err
	}
	return routes, nil
}

func (route *Route) GetRouteString(db *gorm.DB) string {
	var routeName string;
	db.Raw("select STRING_AGG(c.name, ' - ' order by rc.position) from route_checkpoints rc join checkpoints c on c.id = rc.checkpoint_id where rc.route_id = ?", route.Id).Scan(&routeName);
	return routeName;
}

func (routeCheckpoint *RouteCheckpoint) GetRouteCheckpointsByRouteId(db *gorm.DB, routeId uuid.UUID) ([]RouteCheckpoint, error) {
	var routeCheckpoints []RouteCheckpoint
	if err := db.Find(&routeCheckpoints, "route_id = ?", routeId).Error; err != nil {
		return nil, err
	}
	return routeCheckpoints, nil
}

func (route *Route) GetAllRoutes(db *gorm.DB) ([]Route, error) {
	var routes []Route
	if err := db.Find(&routes).Error; err != nil {
		return nil, err
	}
	return routes, nil
}

func (routeCheckpoint *RouteCheckpoint) SaveRouteCheckpoint(db *gorm.DB) error {
	return db.Save(routeCheckpoint).Error
} 

func (Route *Route) CreateRoute(db *gorm.DB) error {
	return db.Create(Route).Error
}

func (Route *Route) SaveRoute(db *gorm.DB) error {
	return db.Save(Route).Error
}

func (routeCheckpoint *RouteCheckpoint) BeforeCreate(tx *gorm.DB) (err error) {
	routeCheckpoint.Id = uuid.New()
	return
}

func (route *Route) GetById(db *gorm.DB, id uuid.UUID) error {
	return db.First(route, "id = ?", id).Error
}

func (routeCheckpoint *RouteCheckpoint) GetById(db *gorm.DB, id uuid.UUID) error {
	return db.First(routeCheckpoint, "id = ?", id).Error
}

func (routeCheckpoint *RouteCheckpoint) GetNextCheckpoint(db *gorm.DB, routeId uuid.UUID, position uint) error {
	return db.First(routeCheckpoint, "route_id = ? AND position = ?", routeId, position+1).Error
}

func (routeCheckpoint *RouteCheckpoint) GetRouteCheckpoint(db *gorm.DB, routeId uuid.UUID, checkpointId uuid.UUID) error {
	return db.First(routeCheckpoint, "route_id = ? AND checkpoint_id = ?", routeId, checkpointId).Error
}

func (routeCheckpoint *RouteCheckpoint) IsNextCheckpoint(db *gorm.DB, route Route) bool {
	var nextCheckpoint RouteCheckpoint
	if err := db.First(&nextCheckpoint, "route_id = ? AND position = ?", route.Id, routeCheckpoint.Position).Error; err != nil {
		return false
	}
	return true
}
