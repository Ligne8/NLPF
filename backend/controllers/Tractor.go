package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"tms-backend/models"
)

type TractorController struct {
	Db *gorm.DB
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

	tractor.TrafficManagerId = requestBody.TrafficManagerId

	// C le bon state ?
	tractor.State = models.StateInTransit

	//Il faut ajouter le traffic manager au tracteur direct comme ca ? pq pas juste utiliser la cle etrangere ?
	//tractor.TrafficManager = &trafficManager

	if err := TractorController.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
