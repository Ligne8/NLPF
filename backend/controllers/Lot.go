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

// CreateLot : Create a new Lot
//
// @Summary      Create a new Lot
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        resource_type  body  string  true  "Resource Type"
// @Param        volume  body  float64  true  "Volume"
// @Param        start_checkpoint_id  body  string  true  "Start Checkpoint Id"
// @Param        end_checkpoint_id  body  string  true  "End Checkpoint Id"
// @Param        owner_id  body  string  true  "Owner Id"
// @Param        current_checkpoint_id  body  string  false  "Current Checkpoint Id"
// @Param        state  body  string  true  "State"
// @Param        max_price_by_km  body  float64  true  "Max Price By Km"
// @Success      201  {object}  models.Lot
// @Failure      400  "Invalid request payload"
// @Failure      500  "Unable to create lot"
// @Router       /lots [post]
func (LotController *LotController) CreateLot(c *gin.Context) {
	var requestBody struct {
		ResourceType        models.ResourceType `json:"resource_type" binding:"required"`
		Volume              float64             `json:"volume" binding:"required"`
		StartCheckpointId   uuid.UUID           `json:"start_checkpoint_id" binding:"required"`
		EndCheckpointId     uuid.UUID           `json:"end_checkpoint_id" binding:"required"`
		OwnerId             uuid.UUID           `json:"owner_id" binding:"required"`
		CurrentCheckpointId uuid.UUID           `json:"current_checkpoint_id"`
		State               models.State        `json:"state" binding:"required"`
		MaxPriceByKm        float64             `json:"max_price_by_km" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var simulation models.Simulation
	if err := LotController.Db.First(&simulation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch simulation date"})
		return
	}

	LotModel := models.Lot{
		ResourceType:        requestBody.ResourceType,
		Volume:              requestBody.Volume,
		StartCheckpointId:   &requestBody.StartCheckpointId,
		EndCheckpointId:     &requestBody.EndCheckpointId,
		CurrentCheckpointId: &requestBody.CurrentCheckpointId,
		CreatedAt:           simulation.SimulationDate,
		OwnerId:             requestBody.OwnerId,
		State:               requestBody.State,
		MaxPriceByKm:        requestBody.MaxPriceByKm,
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

// ListLotsByOwner : List all lots by owner
//
// @Summary      List all lots by owner
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        owner_id  path  string  true  "Owner Id"
// @Success      200  {array}  models.Lot
// @Failure      400  "Invalid owner_id"
// @Failure      500  "Unable to retrieve lots"
// @Router       /lots/owner/{owner_id} [get]
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

// IsCompatible : Check if a lot is compatible with a tractor
//
// @Summary      Check if a lot is compatible with a tractor
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  body  string  true  "Lot Id"
// @Param        tractor_id  body  string  true  "Tractor Id"
// @Success      200  "Lot is compatible with the tractor"
// @Failure      400  "Lot exceeds tractor's capacity"
// @Failure      400  "Lot is not the same resource type as the tractor"
// @Failure      404  "Lot not found"
// @Failure      404  "Tractor not found"
// @Router       /lots/compatible [post]
func (LotController *LotController) IsCompatible(c *gin.Context) {
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

// UpdateLotState : Update the state of a lot
//
// @Summary      Update the state of a lot
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  body  string  true  "Lot Id"
// @Param        state  body  string  true  "State"
// @Success      200  {object}  models.Lot
// @Failure      400  "Invalid request payload"
// @Failure      404  "Lot not found"
// @Failure      500  "Unable to update lot state"
// @Router       /lots/state [put]
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

// AssociateToTrafficManager : Associate a lot to a traffic manager
//
// @Summary      Associate a lot to a traffic manager
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  body  string  true  "Lot Id"
// @Param        traffic_manager_id  body  string  true  "Traffic Manager Id"
// @Success      200  {object}  models.Lot
// @Failure      400  "Invalid lot_id"
// @Failure      400  "Invalid traffic_manager_id"
// @Failure      404  "Lot not found"
// @Failure      500  "Unable to update traffic_manager"
// @Failure      500  "Unable to update state"
// @Router       /lots/associate [put]
func (LotController *LotController) AssociateToTrafficManager(c *gin.Context) {
	var requestBody struct {
		LotId            string `json:"lot_id" binding:"required"`
		TrafficManagerId string `json:"traffic_manager_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var lotIdUUID, errIdUUID = uuid.Parse(requestBody.LotId)
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lot_id"})
		return
	}
	var trafficManagerIdUUID, errTrafficManagerIdUUID = uuid.Parse(requestBody.TrafficManagerId)
	if errTrafficManagerIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traffic_manager_id"})
		return
	}

	var lot models.Lot
	lot, err := lot.FindById(LotController.Db, lotIdUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
		return
	}
	if err := lot.AssociateTraficManager(LotController.Db, trafficManagerIdUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error updating traffic_manager": err.Error()})
		return

	}

	if err := lot.UpdateState(LotController.Db, models.StatePending); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error updating state": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lot)
}

// ListLotsByTrafficManager Get all lots associated with a traffic manager
//
//	@Summary      Get all lots associated with a traffic manager
//	@Tags         lots
//	@Accept       json
//	@Produce      json
//	@Param        traffic_manager_id  path  string  true  "Traffic Manager ID"
//	@Success      200  {array}  models.Lot
//	@Failure      400  {object}  error
//	@Failure      404  {object}  error
//	@Failure      500  {object}  error
//	@Router       /lots/traffic_manager/{traffic_manager_id} [get]
func (LotController *LotController) ListLotsByTrafficManager(c *gin.Context) {
	var lots []models.Lot
	var lotModel models.Lot
	trafficManagerId := c.Param("traffic_manager_id")
	ownerIdUUID, errIdUUID := uuid.Parse(trafficManagerId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traffic_manager_id"})
		return
	}

	var trafficManager models.User
	if err := LotController.Db.First(&trafficManager, "id = ? AND role = ?", ownerIdUUID, "traffic_manager").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traffic Manager not found"})
		return
	}

	lots, err := lotModel.GetLotsByTrafficManager(LotController.Db, ownerIdUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lots)
}

// ListCompatibleTractorsForLot : Get all compatible tractors for a lot
//
// @Summary      Get all compatible tractors for a lot
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  path  string  true  "Lot Id"
// @Param        traffic_manager_id  path  string  true  "Traffic Manager Id"
// @Success      200  {array}  models.Tractor
// @Failure      400  "Invalid lot_id"
// @Failure      400  "Invalid traffic_manager_id"
// @Failure      404  "Lot not found"
// @Failure      404  "Traffic Manager not found"
// @Failure      500  "Unable to retrieve tractors"
// @Router       /lots/tractors/compatible/{traffic_manager_id}/{lot_id} [get]
func (LotController *LotController) ListCompatibleTractorsForLot(c *gin.Context) {
	lotId := c.Param("lot_id")
	trafficManagerId := c.Param("traffic_manager_id")

	lotIdUUID, errLotId := uuid.Parse(lotId)
	if errLotId != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lot_id"})
		return
	}

	trafficManagerIdUUID, errTrafficManagerId := uuid.Parse(trafficManagerId)
	if errTrafficManagerId != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid traffic_manager_id"})
		return
	}

	var lot models.Lot
	if err := LotController.Db.First(&lot, "id = ?", lotIdUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
		return
	}

	var trafficManager models.User
	if err := LotController.Db.First(&trafficManager, "id = ? AND role = ?", trafficManagerIdUUID, "traffic_manager").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Traffic Manager not found"})
		return
	}

	var tractors []models.Tractor
	if err := LotController.Db.Preload("EndCheckpoint").Preload("StartCheckpoint").Preload("CurrentCheckpoint").Preload("TrafficManager").Where("traffic_manager_id = ?", trafficManagerIdUUID).Find(&tractors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve tractors"})
		return
	}

	var compatibleTractors []models.Tractor
	for _, tractor := range tractors {
		if LotController.checkCompatibility(lot, tractor) {
			compatibleTractors = append(compatibleTractors, tractor)
		}
	}

	c.JSON(http.StatusOK, compatibleTractors)
}

// checkCompatibility : Check if a lot is compatible with a tractor
func (LotController *LotController) checkCompatibility(lot models.Lot, tractor models.Tractor) bool {
	if lot.Volume > (tractor.MaxVolume - tractor.CurrentVolume) {
		return false
	}

	if lot.ResourceType != tractor.ResourceType {
		return false
	}

	return true
}
