package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TractorController struct {
	Db *gorm.DB
}

//var tractorModel = models.Tractor{}

func (TractorController *TractorController) CreateTractor(c *gin.Context) {
	var requestBody struct {
		Name                string              `json:"name" binding:"required"`
		ResourceType        models.ResourceType `json:"resource_type" binding:"required"`
		MaxVolume           float64             `json:"volume" binding:"required"`
		StartCheckpointId   uuid.UUID           `json:"start_checkpoint_id" binding:"required"`
		EndCheckpointId     uuid.UUID           `json:"end_checkpoint_id" binding:"required"`
		OwnerId             uuid.UUID           `json:"owner_id" binding:"required"`
		CurrentCheckpointId uuid.UUID           `json:"current_checkpoint_id"`
		State               models.State        `json:"state" binding:"required"`
		MinPriceByKm        float64             `json:"min_price_by_km" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var simulation models.Simulation
	if err := TractorController.Db.First(&simulation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch simulation date"})
		return
	}

	TractorModel := models.Tractor{
		Name:                requestBody.Name,
		ResourceType:        requestBody.ResourceType,
		MaxVolume:           requestBody.MaxVolume,
		StartCheckpointId:   &requestBody.StartCheckpointId,
		EndCheckpointId:     &requestBody.EndCheckpointId,
		CurrentCheckpointId: &requestBody.CurrentCheckpointId,
		CreatedAt:           simulation.SimulationDate,
		OwnerId:             requestBody.OwnerId,
		State:               requestBody.State,
		MinPriceByKm:        requestBody.MinPriceByKm,
	}

	if err := TractorController.Db.Create(&TractorModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := TractorController.Db.Preload("StartCheckpoint").Preload("EndCheckpoint").Preload("Owner").Preload("TrafficManager").Preload("Trader").First(&TractorModel, "id = ?", TractorModel.Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TractorModel)
}

func (TractorController *TractorController) GoToNextCheckpoint(c *gin.Context) {
	var tractors []models.Tractor
	if err := TractorController.Db.Find(&tractors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// for each tractor
	for _, tractor := range tractors {

		tractor.UpdateNextCheckpoint(TractorController.Db)

		TractorController.Db.Save(&tractor)
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) ListTractorsByOwner(c *gin.Context) {
	var tractors []models.Tractor
	var tractorModel models.Tractor
	ownerId := c.Param("ownerId")
	ownerIdUUID, errIdUUID := uuid.Parse(ownerId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
		return
	}

	tractors, err := tractorModel.GetByOwnerId(TractorController.Db, ownerIdUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) ListTractorsByTrafficManagerId(c *gin.Context) {
	var tractors []models.Tractor
	var tractorModel models.Tractor
	trafficManagerId := c.Param("trafficManagerId")
	trafficManagerIdUUID, errIdUUID := uuid.Parse(trafficManagerId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
		return
	}

	tractors, err := tractorModel.GetByTrafficManagerId(TractorController.Db, trafficManagerIdUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) ListTractorsByState(c *gin.Context) {
	var tractors []models.Tractor
	var tractorModel models.Tractor
	state := c.Param("state")

	tractors, err := tractorModel.GetByState(TractorController.Db, state)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) ListTractorsByRouteId(c *gin.Context) {
	var tractors []models.Tractor
	var tractorModel models.Tractor
	routeId := c.Param("routeId")
	routeIdUUID, errIdUUID := uuid.Parse(routeId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
		return
	}

	tractors, err := tractorModel.GetByRouteId(TractorController.Db, routeIdUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) AssociateToTrafficManager(c *gin.Context) {
	var requestBody struct {
		TractorId        string `json:"tractor_id" binding:"required"`
		TrafficManagerId string `json:"traffic_manager_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tractorIdUUID, errIdUUID = uuid.Parse(requestBody.TractorId)
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tractor_id"})
		return
	}
	var trafficManagerIdUUID, errTrafficManagerIdUUID = uuid.Parse(requestBody.TrafficManagerId)
	if errTrafficManagerIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traffic_manager_id"})
		return
	}

	var tractor models.Tractor
	tractor, err := tractor.FindById(TractorController.Db, tractorIdUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found"})
		return
	}
	if err := tractor.AssociateTraficManager(TractorController.Db, trafficManagerIdUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error updating traffic_manager": err.Error()})
		return
	}

	if err := tractor.UpdateState(TractorController.Db, models.StatePending); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractor)
}

// UpdateTractorState: Update the state of a tractor
func (TractorController *TractorController) UpdateTractorState(c *gin.Context) {
	var requestBody struct {
		Id    string    `json:"id" binding:"required"`
		State models.State `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tractorIdUUID, errIdUUID = uuid.Parse(requestBody.Id)
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tractor_id"})
		return
	}
	var tractor models.Tractor
	tractor, err := tractor.FindById(TractorController.Db, tractorIdUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found"})
		return
	}
	tractor.State = requestBody.State
	

	if err := TractorController.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if requestBody.State == models.StateInTransit {
		var lot models.Lot
		lot.UpdateStateByTractorId(TractorController.Db, tractorIdUUID, models.StateInTransit);
	} else if requestBody.State == models.StatePending {
		var lot models.Lot
		lot.UpdateStateByTractorId(TractorController.Db, tractorIdUUID, models.StatePending);
	}
	c.JSON(http.StatusOK, tractor)
}

func (TractorController *TractorController) BindRoute(c *gin.Context) {
	var requestBody struct {
		TractorId string `json:"tractor_id" binding:"required"`
		RouteId   string `json:"route_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tractorIdUUID, errIdUUID = uuid.Parse(requestBody.TractorId)
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tractor_id"})
		return
	}
	var routeIdUUID, errRouteIdUUID = uuid.Parse(requestBody.RouteId)
	if errRouteIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route_id"})
		return
	}

	var tractor models.Tractor
	tractor, err := tractor.FindById(TractorController.Db, tractorIdUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found"})
		return
	}

	tractor.RouteId = &routeIdUUID
	if err := tractor.Save(TractorController.Db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractor)
}

func (TractorController *TractorController) UnbindeRoute(c *gin.Context) {
	var requestBody struct {
		TractorId string `json:"tractor_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tractorIdUUID, errIdUUID = uuid.Parse(requestBody.TractorId)
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tractor_id"})
		return
	}

	var tractor models.Tractor
	tractor, err := tractor.FindById(TractorController.Db, tractorIdUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found"})
		return
	}

	tractor.RouteId = nil
	if err := tractor.Save(TractorController.Db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractor)
}