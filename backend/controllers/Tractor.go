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

// CreateTractor : Create a new tractor
//
// @Summary      Create a new tractor
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        name                body  string  true  "Name"
// @Param        resource_type        body  string  true  "Resource Type"
// @Param        volume               body  float64 true  "Volume"
// @Param        start_checkpoint_id   body  string  true  "Start Checkpoint Id"
// @Param        end_checkpoint_id     body  string  true  "End Checkpoint Id"
// @Param        owner_id             body  string  true  "Owner Id"
// @Param        current_checkpoint_id body  string  false "Current Checkpoint Id"
// @Param        state                body  string  true  "State"
// @Param        min_price_by_km      body  float64 true  "Min Price By Km"
// @Success      201  {object}  models.Tractor
// @Failure      400  "Invalid request payload"
// @Failure      500  "Unable to create tractor"
// @Router       /tractors [post]
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

// GoToNextCheckpoint : Update the current checkpoint of the tractors
//
// @Summary      Update the current checkpoint of the tractors
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Tractor
// @Failure      500  "Unable to fetch tractors"
// @Router       /tractors/next_checkpoint [put]
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

// ListTractorsByOwner : List all tractors by owner id
//
// @Summary      List all tractors by owner id
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        owner_id  path  string  true  "Owner Id"
// @Success      200  {array}  models.Tractor
// @Failure      400  "Invalid owner_id"
// @Failure      500  "Unable to retrieve tractors"
// @Router       /tractors/owner/{owner_id} [get]
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

// ListTractorsByTrafficManagerId : List all tractors by traffic manager id
//
// @Summary      List all tractors by traffic manager id
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        traffic_manager_id  path  string  true  "Traffic Manager Id"
// @Success      200  {array}  models.Tractor
// @Failure      400  "Invalid owner_id"
// @Failure      500  "Unable to retrieve tractors"
// @Router       /tractors/traffic_manager/{traffic_manager_id} [get]
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

// ListTractorsByState : List all tractors by state
//
// @Summary      List all tractors by state
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        state  path  string  true  "State"
// @Success      200  {array}  models.Tractor
// @Failure      400  "Invalid state"
// @Failure      500  "Unable to retrieve tractors"
// @Router       /tractors/state/{state} [get]
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

// ListTractorsByRouteId : List all tractors by route id
//
// @Summary      List all tractors by route id
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        route_id  path  string  true  "Route Id"
// @Success      200  {array}  models.Tractor
// @Failure      400  "Invalid route_id"
// @Failure      500  "Unable to retrieve tractors"
// @Router       /tractors/route/{route_id} [get]
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

// AssociateToTrafficManager : Associate a tractor to a traffic manager
//
// @Summary      Associate a tractor to a traffic manager
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        tractor_id  body  string  true  "Tractor Id"
// @Param        traffic_manager_id  body  string  true  "Traffic Manager Id"
// @Success      200  {object}  models.Tractor
// @Failure      400  "Invalid tractor_id"
// @Failure      400  "Invalid traffic_manager_id"
// @Failure      404  "Tractor not found"
// @Failure      500  "Unable to update traffic_manager"
// @Failure      500  "Unable to update state"
// @Router       /tractors/associate [put]
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

// UpdateTractorState : Update the state of a tractor
//
// @Summary      Update the state of a tractor
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        id     body  string  true  "Id"
// @Param        state  body  string  true  "State"
// @Success      200  {object}  models.Tractor
// @Failure      400  "Invalid request payload"
// @Failure      404  "Tractor not found"
// @Failure      500  "Unable to update tractor"
// @Router       /tractors/state [put]
func (TractorController *TractorController) UpdateTractorState(c *gin.Context) {
	var requestBody struct {
		Id    string       `json:"id" binding:"required"`
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
		lot.UpdateStateByTractorId(TractorController.Db, tractorIdUUID, models.StateInTransit)
	} else if requestBody.State == models.StatePending {
		var lot models.Lot
		lot.UpdateStateByTractorId(TractorController.Db, tractorIdUUID, models.StatePending)
	}
	c.JSON(http.StatusOK, tractor)
}

// BindRoute : Bind a route to a tractor
//
// @Summary      Bind a route to a tractor
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        tractor_id  body  string  true  "Tractor Id"
// @Param        route_id    body  string  true  "Route Id"
// @Success      200  {object}  models.Tractor
// @Failure      400  "Invalid request payload"
// @Failure      404  "Tractor not found"
// @Failure      500  "Unable to update tractor"
// @Router       /tractors/bind_route [put]
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

// UnbindRoute : Unbind a route from a tractor
//
// @Summary      Unbind a route from a tractor
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        tractor_id  body  string  true  "Tractor Id"
// @Success      200  {object}  models.Tractor
// @Failure      400  "Invalid request payload"
// @Failure      404  "Tractor not found"
// @Failure      500  "Unable to update tractor"
// @Router       /tractors/unbind_route [put]
func (TractorController *TractorController) UnbindRoute(c *gin.Context) {
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

// DeleteTractorByOwnerId : Delete a tractor
// @Summary      Delete a tractor with the tractor id
// @Tags         tractors
// @Accept       json
// @Produce      json
// @Param        tractor_id  path  string  true  "Tractor Id"
// @Success      200  "Tractor deleted successfully"
// @Failure      400  "Invalid tractor_id"
// @Failure      404  "Tractor not found"
// @Failure      500  "Unable to delete tractor"
// @Router       /tractors/{tractor_id} [delete]
func (TractorController *TractorController) DeleteTractor(c *gin.Context) {
	tractorId := c.Param("tractor_id")
	tractorIdUUID, errIdUUID := uuid.Parse(tractorId)

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

	if err := TractorController.Db.Delete(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tractor \"" + tractorId + "\" deleted successfully"})
}
