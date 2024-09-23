package controllers

import (
	"github.com/google/uuid"
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LotController struct {
	Db *gorm.DB
}

// AddLot gère la création d'un Lot
// @version 1
// @Summary Add Lot
// @Tags Lots
// @Accept json
// @Param   type 	 body   string     true        "Type of the lot" example(Liquid)
// @Param   unit_volume 	body    uint     true        "Standard volume of the lot" example(10)
// @Param   weight 	body    uint     true        "weight in KG of the lot" example(100)
// @Param   start_checkpoint_id 	body    uint     true        "Starting checkpoint ID of the Lot" example(1)
// @Param   end_checkpoint_id 	body    uint     true        "Ending checkpoint ID of the Lot" example(2)
// @Produce json
// @Failure 500 "Error from gin"
// @Success 201 "No content"
// @Router /lots [post]
func (LotController *LotController) AddLot(c *gin.Context) {
	var LotModel models.Lot

	type RequestBody struct {
		LotType           string    `json:"type"`
		Units             uint      `json:"units"`
		StartCheckpointId uuid.UUID `json:"start_checkpoint_id"`
		EndCheckpointId   uuid.UUID `json:"end_checkpoint_id"`
		OwnerId           uuid.UUID `json:"owner_id"`
		TrafficManagerId  uuid.UUID `json:"traffic_manager_id"`
		TraderId          uuid.UUID `json:"trader_id"`
		MaxPriceByKm      float64   `json:"max_price_by_km"`
	}

	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		Err500(c, err)
		return
	}

	// Populate the Lot model
	LotModel.Type = requestBody.LotType
	LotModel.Units = requestBody.Units
	LotModel.StartCheckpoint = models.Checkpoint{Id: requestBody.StartCheckpointId}
	LotModel.EndCheckpoint = models.Checkpoint{Id: requestBody.EndCheckpointId}
	LotModel.OwnerId = models.User{Id: requestBody.OwnerId}
	LotModel.TrafficManagerId = models.User{Id: requestBody.TrafficManagerId}
	LotModel.TraderId = models.User{Id: requestBody.TraderId}
	LotModel.MaxPriceByKm = requestBody.MaxPriceByKm

	if err := LotModel.Save(LotController.Db); err != nil {
		Err500(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
