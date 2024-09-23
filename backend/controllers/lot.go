package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"tms-backend/models"
)

type LotController struct {
	Db *gorm.DB
}

func (LotController *LotController) AddLot(c *gin.Context) {
	var requestBody struct {
		Type              string    `json:"type" binding:"required"`
		Units             uint      `json:"units" binding:"required"`
		StartCheckpointId uuid.UUID `json:"start_checkpoint_id" binding:"required"`
		EndCheckpointId   uuid.UUID `json:"end_checkpoint_id" binding:"required"`
		OwnerId           uuid.UUID `json:"owner_id" binding:"required"`
		State             string    `json:"state" binding:"required"`
		MaxPriceByKm      float64   `json:"max_price_by_km" binding:"required"`
		TrafficManagerId  uuid.UUID `json:"traffic_manager_id" binding:"required"`
		TraderId          uuid.UUID `json:"trader_id" binding:"required"`
	}

	// Lier les données du JSON à requestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Créer une instance du modèle Lot
	LotModel := models.Lot{
		Type:              requestBody.Type,
		Units:             requestBody.Units,
		StartCheckpointId: requestBody.StartCheckpointId,
		EndCheckpointId:   requestBody.EndCheckpointId,
		OwnerId:           requestBody.OwnerId,
		State:             requestBody.State,
		MaxPriceByKm:      requestBody.MaxPriceByKm,
		TrafficManagerId:  requestBody.TrafficManagerId,
		TraderId:          requestBody.TraderId,
	}

	// Sauvegarder dans la base de données
	if err := LotController.Db.Create(&LotModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retourner la réponse JSON avec le lot créé
	c.JSON(http.StatusCreated, LotModel)
}
