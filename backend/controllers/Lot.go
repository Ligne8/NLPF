package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LotController struct {
	Db *gorm.DB
}

// CreateLot Create a new Lot
//
//  @Summary      List accounts
//  @Description  get accounts
//  @Tags         accounts
//  @Accept       json
//  @Produce      json
//  @Param        q    query     string  false  "name search by q"  Format(email)
//  @Success      200  {array}   model.Account
//  @Failure      500  {object}  httputil.HTTPError
//  @Router       /accounts [get]
func (LotController *LotController) CreateLot(c *gin.Context) {
	var requestBody struct {
		ResourceType      models.ResourceType `json:"resource_type" binding:"required"`
		Volume            float64             `json:"volume" binding:"required"`
		StartCheckpointId uuid.UUID           `json:"start_checkpoint_id" binding:"required"`
		EndCheckpointId   uuid.UUID           `json:"end_checkpoint_id" binding:"required"`
		OwnerId           uuid.UUID           `json:"owner_id" binding:"required"`
		CurrentCheckpointId uuid.UUID         `json:"current_checkpoint_id"`
		State             models.State        `json:"state" binding:"required"`
		MaxPriceByKm      float64             `json:"max_price_by_km" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	LotModel := models.Lot{
		ResourceType:      requestBody.ResourceType,
		Volume:            requestBody.Volume,
		StartCheckpointId: &requestBody.StartCheckpointId,
		EndCheckpointId:   &requestBody.EndCheckpointId,
		OwnerId:           requestBody.OwnerId,
		State:             requestBody.State,
		MaxPriceByKm:      requestBody.MaxPriceByKm,
		TrafficManagerId:  &requestBody.TrafficManagerId,
		TraderId:          &requestBody.TraderId,
	}

	if err := LotController.Db.Create(&LotModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := LotController.Db.Preload("StartCheckpoint").Preload("EndCheckpoint").Preload("Tractor").Preload("Owner").Preload("TrafficManager").Preload("Trader").First(&LotModel, "id = ?", LotModel.Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, LotModel)
}

func (LotController *LotController) ListLotsByOwner(c *gin.Context) {
	var lots []models.Lot
	var lotModel models.Lot
	ownerId := c.Param("owner_id")
	ownerIdUUID, errIdUUID := uuid.Parse(ownerId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
		return
	}

	lots, err := lotModel.GetLotsByOwner(LotController.Db, ownerIdUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lots)
}

func (LotController *LotController) isCompatible(c *gin.Context) {
	var requestBody struct {
		LotId     uuid.UUID `json:"lot_id" binding:"required"`
		TractorId uuid.UUID `json:"tractor_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the Lot
	var lot models.Lot
	if err := LotController.Db.First(&lot, "id = ?", requestBody.LotId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
		return
	}

	// Get the tractor
	var tractor models.Tractor
	if err := LotController.Db.First(&tractor, "id = ?", requestBody.TractorId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found"})
		return
	}

	// Check if there is enough space
	if lot.Volume > (tractor.MaxVolume - tractor.CurrentVolume) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lot exceeds tractor's capacity"})
		return
	}

	if lot.ResourceType != tractor.ResourceType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lot is not the same resource type as the tractor"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lot is compatible with the tractor"})
}

// UpdateLotState: Update the state of a lot
func (LotController *LotController) UpdateLotState(c *gin.Context) {
	var requestBody struct {
		LotId uuid.UUID    `json:"lot_id" binding:"required"`
		State models.State `json:"state" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the Lot
	var lot models.Lot
	if err := LotController.Db.First(&lot, "id = ?", requestBody.LotId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
		return
	}

	lot.State = requestBody.State
	// Change the state of the Lot
	if err := LotController.Db.Save(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lot)
}
