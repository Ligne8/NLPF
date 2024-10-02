package controllers

import (
	"fmt"
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TractorController struct {
	Db *gorm.DB
}

var tractorModel = models.Tractor{}

func (TractorController *TractorController) AddTractor(c *gin.Context) {
	var newTractor models.Tractor
	if err := c.ShouldBindJSON(&newTractor); err != nil {
		fmt.Println(newTractor)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var simulation models.Simulation
	if err := TractorController.Db.First(&simulation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch simulation date"})
		return
	}

	newTractor.CreatedAt = simulation.SimulationDate

	if err := TractorController.Db.Create(&newTractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newTractor)
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

func (TractorController *TractorController) GetTractorsByOwnerId(c *gin.Context) {
	ownerId := c.Param("ownerId")
	var tractors []models.Tractor
	var ownerUUID uuid.UUID

	if parsedUUID, err := uuid.Parse(ownerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ownerId"})
		return
	} else {
		ownerUUID = parsedUUID
	}

	tractors, err := tractorModel.GetByOwnerId(TractorController.Db, ownerUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) GetTractorsByTrafficManagerId(c *gin.Context) {
	trafficManagerId := c.Param("trafficManagerId")
	var tractors []models.Tractor
	var trafficManagerUUID uuid.UUID

	if parsedUUID, err := uuid.Parse(trafficManagerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trafficManagerId"})
		return
	} else {
		trafficManagerUUID = parsedUUID
	}

	tractors, err := tractorModel.GetByTrafficManagerId(TractorController.Db, trafficManagerUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) GetTractorsByState(c *gin.Context) {
	state := c.Param("state")
	var tractors []models.Tractor

	tractors, err := tractorModel.GetByState(TractorController.Db, state)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) GetTractorsByRouteId(c *gin.Context) {
	routeId := c.Param("routeId")
	var tractors []models.Tractor
	var routeUUID uuid.UUID

	if parsedUUID, err := uuid.Parse(routeId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid routeId"})
		return
	} else {
		routeUUID = parsedUUID
	}

	tractors, err := tractorModel.GetByRouteId(TractorController.Db, routeUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) AddTrafficManager(c *gin.Context) {
	var requestBody struct {
		TractorId        uuid.UUID `json:"tractor_id" binding:"required"`
		TrafficManagerId uuid.UUID `json:"traffic_manager_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var trafficManager models.User
	if err := TractorController.Db.First(&trafficManager, "id = ?", requestBody.TrafficManagerId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "traffic manager not found when trying to add it to the tractor"})
		return
	}

	if trafficManager.Role == models.RoleTrafficManager {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The user is not a traffic manager"})
		return
	}

	var tractor models.Tractor
	if err := TractorController.Db.First(&tractor, "id = ?", requestBody.TractorId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found when trying to add a traffic manager to it"})
		return
	}

	tractor.TrafficManagerId = &requestBody.TrafficManagerId
	tractor.State = models.StatePending

	if err := TractorController.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// UpdateTractorState: Update the state of a tractor
func (TractorController *TractorController) UpdateTractorState(c *gin.Context) {
	var requestBody struct {
		Id    uuid.UUID    `json:"id" binding:"required"`
		State models.State `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tractor models.Tractor
	if err := TractorController.Db.First(&tractor, "id = ?", requestBody.Id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found when trying to update its state"})
		return
	}
	tractor.State = requestBody.State

	if err := TractorController.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tractor)
}
