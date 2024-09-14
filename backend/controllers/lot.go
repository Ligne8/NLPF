package controllers

import (
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
func (LotController * LotController) AddLot(c *gin.Context) {
	var LotModel models.Lot

	type RequestBody struct {
		LotType string `json:"type"`
		UnitVolume uint `json:"unit_volume"`
		Weight uint `json:"weight"`
		StartCheckpointId uint `json:"start_checkpoint_id"`
		EndCheckpointId uint `json:"end_checkpoint_id"`
	}

	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		Err500(c, err)
		return
	}

	LotModel.Type = requestBody.LotType
	LotModel.UnitVolume = requestBody.UnitVolume
	LotModel.Weight = requestBody.Weight
	LotModel.StartCheckpointId = requestBody.StartCheckpointId
	LotModel.EndCheckpointId = requestBody.EndCheckpointId
	LotModel.TractorId = nil


	if err := LotModel.Save(LotController.Db); err != nil {
		Err500(c, err)
		return
	}
	c.Status(http.StatusCreated)
}