package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceType string

const (
	ResourceTypeBulk   ResourceType = "bulk"
	ResourceTypeSolid  ResourceType = "solid"
	ResourceTypeLiquid ResourceType = "liquid"
)

type Tractor struct {
	Id                  uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey"`
	Name                string       `json:"name" gorm:"not null"`
	ResourceType        ResourceType `json:"resource_type" gorm:"type:varchar(10)" binding:"required"`
	MaxVolume           float64      `json:"max_units" gorm:"not null"`
	CurrentVolume       float64      `json:"current_units" gorm:"not null"`
	StartCheckpointId   *uuid.UUID   `json:"start_checkpoint_id" gorm:""` // Changed to pointer to allow null values
	StartCheckpoint     *Checkpoint  `json:"start_checkpoint" gorm:"foreignKey:StartCheckpointId"`
	EndCheckpointId     *uuid.UUID   `json:"end_checkpoint_id" gorm:""` // Changed to pointer to allow null values
	EndCheckpoint       *Checkpoint  `json:"end_checkpoint" gorm:"foreignKey:EndCheckpointId"`
	CurrentCheckpointId *uuid.UUID   `json:"current_checkpoint_id" gorm:"type:uuid"` // Foreign key for Checkpoint
	CurrentCheckpoint   *Checkpoint  `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
	State               State        `json:"state" gorm:"not null"`
	CreatedAt           time.Time    `json:"created_at" gorm:""`
	OwnerId             uuid.UUID    `json:"owner_id" gorm:"not null"` // Foreign key for User
	Owner               User         `json:"owner" gorm:"foreignKey:OwnerId"`
	MinPriceByKm        float64      `json:"min_price_by_km" gorm:"not null"`
	TrafficManagerId    *uuid.UUID   `json:"traffic_manager_id" gorm:"type:uuid"` // Foreign key for User
	TrafficManager      *User        `json:"traffic_manager" gorm:"foreignKey:TrafficManagerId"`
	TraderId            *uuid.UUID   `json:"trader_id" gorm:"type:uuid"` // Foreign key for User
	Trader              *User        `json:"trader" gorm:"foreignKey:TraderId"`
	RouteId             *uuid.UUID   `json:"route_id" gorm:"type:uuid"` // Foreign key for Route
	Route               *Route       `json:"route" gorm:"foreignKey:RouteId"`
	OfferId             *uuid.UUID   `json:"offer_id" gorm:""`
	Offer               *Offer       `json:"offer" gorm:"foreignKey:OfferId"`
	CurrentPrice        float64      `json:"current_price" gorm:"-"`
}

func (tractor *Tractor) BeforeCreate(tx *gorm.DB) (err error) {
	validTypes := map[ResourceType]bool{
		ResourceTypeBulk:   true,
		ResourceTypeSolid:  true,
		ResourceTypeLiquid: true,
	}
	if !validTypes[tractor.ResourceType] {
		return errors.New("invalid resource type")
	}

	validState := map[State]bool{
		StateAvailable: true,
		StateArchive:   true,
		StateAtTrader:  true,
		StateInTransit: true,
		StateOnMarket:  true,
		StatePending:   true,
	}

	if !validState[tractor.State] {
		return errors.New("invalid valid state")
	}

	if tractor.Id == uuid.Nil {
		tractor.Id = uuid.New()
	}
	return
}

func (tractor *Tractor) Save(db *gorm.DB) error {
	return db.Preload("EndCheckpoint").Preload("StartCheckpoint").Save(tractor).Error
}

func (tractor *Tractor) Update(db *gorm.DB) error {
	return db.Model(&tractor).Updates(tractor).Error
}

func (tractor *Tractor) GetAllTractors(db *gorm.DB) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) FindById(db *gorm.DB, tractorId uuid.UUID) (Tractor, error) {
	var foundTractor Tractor
	if err := db.First(&foundTractor, "id = ?", tractorId).Error; err != nil {
		return Tractor{}, err
	}
	return foundTractor, nil
}

func (tractor *Tractor) GetByOwnerId(db *gorm.DB, ownerId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("CurrentCheckpoint").Preload("TrafficManager").Where("owner_id = ?", ownerId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetTractorsByTrader(db *gorm.DB, traderId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Preload("EndCheckpoint").Preload("StartCheckpoint").Where("trader_id = ?", traderId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetByTrafficManagerId(db *gorm.DB, trafficManagerId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Preload("Route").Preload("StartCheckpoint").Preload("EndCheckpoint").Preload("CurrentCheckpoint").Where("traffic_manager_id = ?", trafficManagerId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetByState(db *gorm.DB, state string) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Where("state = ?", state).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}

func (tractor *Tractor) GetByRouteId(db *gorm.DB, routeId uuid.UUID) ([]Tractor, error) {
	var tractors []Tractor
	if err := db.Where("route_id = ?", routeId).Find(&tractors).Error; err != nil {
		return nil, err
	}
	return tractors, nil
}
func (tractor *Tractor) AssociateTraficManager(db *gorm.DB, trafficManagerId uuid.UUID) error {
	return db.Model(&tractor).Update("traffic_manager_id", trafficManagerId).Error
}

func (tractor *Tractor) UpdateState(db *gorm.DB, state State) error {
	return db.Model(&tractor).Update("state", state).Error
}

func (tractor *Tractor) UpdateNextCheckpoint(db *gorm.DB) error {
	var routeCheckpoint RouteCheckpoint
	if err := routeCheckpoint.GetRouteCheckpoint(db, *tractor.RouteId, *tractor.CurrentCheckpointId); err != nil {
		return err
	}
	var position uint = routeCheckpoint.Position
	var nextCheckpoint RouteCheckpoint
	if err := nextCheckpoint.GetNextCheckpoint(db, *tractor.RouteId, position); err != nil {
		tractor.State = StateArchive
		return nil
	}

	//tractor.CurrentCheckpointId = nextCheckpoint.CheckpointId
	if nextCheckpoint.IsNextCheckpoint(db, *tractor.Route) {
		tractor.State = StateInTransit
	} else {
		tractor.State = StateArchive
	}
	return nil
}

func (tractor *Tractor) GetVolumeAtCheckpoint(db *gorm.DB, checkpointId uuid.UUID) (float64, error) {
	var transactionModel Transaction
	if tractor.RouteId == nil {
		return 0, errors.New("Tractor has no route")
	}
	var routeModel Route
	// je récupère le checkpoint actuel
	var currentCheckpointId uuid.UUID = checkpointId
	var currentRouteCheckpoint RouteCheckpoint
	var allRouteCheckpoints []RouteCheckpoint
	// je récupère tous les checkpoints de la route du tracteur
	allRouteCheckpoints, err := routeModel.GetRouteCheckpoint(db, *tractor.RouteId)
	if err != nil {
		return 0, err
	}
	// je récupère le route checkpoint actuel
	if err := currentRouteCheckpoint.GetRouteCheckpoint(db, *tractor.RouteId, currentCheckpointId); err != nil {
		return 0, err
	}
	// je filtre les checkpoints pour ne garder que ceux avant le checkpoint actuel
	var filteredRouteCheckpoints []RouteCheckpoint
	for _, checkpoint := range allRouteCheckpoints {
		if checkpoint.Position <= currentRouteCheckpoint.Position {
			filteredRouteCheckpoints = append(filteredRouteCheckpoints, checkpoint)
		}
	}
	allRouteCheckpoints = filteredRouteCheckpoints
	var result float64 = 0
	// je parcours les transactions pour calculer le volume du tracteur
	for _, checkpoint := range allRouteCheckpoints {
		var transaction []Transaction
		// je récupère les transactions du tracteur pour le checkpoint actuel
		transaction, err = transactionModel.FindByRouteIdAndCheckpointIdAndTractorId(db, *tractor.RouteId, checkpoint.CheckpointId, tractor.Id);
		if err != nil{
			return 0, err;
		}
		// je parcours les transactions pour calculer le volume du tracteur
		for _, transaction := range transaction {
			// je vérifie si la transaction est une entrée ou une sortie
			if transaction.TransactionType == TransactionState(TransactionStateIn) {
				// si c'est une entrée, j'ajoute le volume
				result += transaction.Lot.Volume
			} else {
				// si c'est une sortie, je soustrait le volume
				result -= transaction.Lot.Volume
			}
		}
	}
	// si le checkpoint n'a pas été trouvé, je retourne une erreur
	return result, nil
}

func (tractor *Tractor) ExecTransaction(db *gorm.DB) error {
	var transactionModel Transaction
	var transactions []Transaction
	transactions, err := transactionModel.FindByRouteId(db, *tractor.RouteId)
	if err != nil {
		return err
	}
	var routeCheckpoint RouteCheckpoint
	err = routeCheckpoint.GetRouteCheckpoint(db, *tractor.RouteId, *tractor.CurrentCheckpointId)
	if err != nil {
		return err
	}
	// for each transaction
	for _, transaction := range transactions {
		if transaction.CheckpointId != tractor.CurrentCheckpointId {
			continue
		}

	}
	return nil
}
