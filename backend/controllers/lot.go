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

func (LotController * LotController) ListLots(c *gin.Context) {
	var LotModel models.Lot
	lots, err := LotModel.ListAll(LotController.Db)
	if err != nil {
		Err500(c, err)
		return
	}
	c.JSON(http.StatusOK, lots)
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

// associateTractorToLot gère l'association d'un tracteur à un lot
// @version 1
// @Summary Associate Tractor to Lot
// @Tags Lots
// @Accept json
// @Param   lot_id 	 path   uint     true        "ID of the lot" example(1)
// @Param   tractor_id 	body    uint     true        "ID of the tractor" example(1)
// @Produce json
// @Failure 500 "Error from gin"
// @Success 200 "No content"
// @Router /lots/{lot_id}/ [patch]
func (LotController * LotController) AssociateTractorToLot(c *gin.Context) {
	var LotModel models.Lot

	lotId := c.Param("lot_id")
	tractorId := c.PostForm("tractor_id")



	if err := LotModel.AssociateTractorToLot(LotController.Db, lotId, tractorId); err != nil {
		Err500(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (LotController * LotController) PatchLot(c *gin.Context) {

	var LotModel *models.Lot

	LotModel, err := LotModel.GetByID(LotController.Db, c.Param("id"))

	if err != nil {
		Err500(c, err)
		return
	}

	type LotRequestPatch struct {
		Type *string `json:"type"`
		UnitVolume *uint `json:"unit_volume"`
		Weight *uint `json:"weight"`
		StartCheckpointId *uint `json:"start_checkpoint_id"`
		EndCheckpointId *uint `json:"end_checkpoint_id"`
		TractorId *uint `json:"tractor_id"`
	}

	var RequestBody LotRequestPatch

	if err := c.ShouldBindJSON(&RequestBody); err != nil {
		Err500(c, err)
		return
	}

	if RequestBody.Type != nil {
		LotModel.Type = *RequestBody.Type
	}
	if RequestBody.UnitVolume != nil {
		LotModel.UnitVolume = *RequestBody.UnitVolume
	}
	if RequestBody.Weight != nil {
		LotModel.Weight = *RequestBody.Weight
	}
	if RequestBody.StartCheckpointId != nil {
		LotModel.StartCheckpointId = *RequestBody.StartCheckpointId
	}
	if RequestBody.EndCheckpointId != nil {
		LotModel.EndCheckpointId = *RequestBody.EndCheckpointId
	}
	if RequestBody.TractorId != nil {
		LotModel.TractorId = RequestBody.TractorId
	}
	LotModel.Save(LotController.Db)

	c.JSON(http.StatusOK, LotModel)
}



