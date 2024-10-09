package controllers

import (
	"net/http"
	"time"
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
// @Router       /lots/traffic_manager [post]
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
	lot.TrafficManagerId = &trafficManagerIdUUID
	lot.State = models.StatePending
	if err := lot.Save(LotController.Db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error updating traffic_manager": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lot)
}

// DeleteLot : Delete a lot
// @Summary      Delete a lot with the lot id
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  path  string  true  "Lot Id"
// @Success      200  "Lot deleted successfully"
// @Failure      400  "Invalid lot_id"
// @Failure      404  "Lot not found"
// @Failure      500  "Unable to delete lot"
// @Router       /lots/{lot_id} [delete]
func (LotController *LotController) DeleteLot(c *gin.Context) {
	lotId := c.Param("lot_id")
	lotIdUUID, errIdUUID := uuid.Parse(lotId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lot_id"})
		return
	}

	var lot models.Lot
	lot, err := lot.FindById(LotController.Db, lotIdUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lot not found"})
		return
	}

	if err := LotController.Db.Delete(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lot '" + lotId + "' deleted successfully"})
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
	if err := LotController.Db.Preload("EndCheckpoint").Preload("Route").Preload("StartCheckpoint").Preload("CurrentCheckpoint").Preload("TrafficManager").Where("traffic_manager_id = ?", trafficManagerIdUUID).Find(&tractors).Error; err != nil {
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
	if tractor.RouteId == nil {
		return false
	}
	if lot.ResourceType != tractor.ResourceType {
		return false
	}
	if tractor.State != models.StatePending {
		return false
	}

	return LotController.checkTractorCheckpointCompatibility(lot, tractor)
}

func (LotController *LotController) checkTractorCheckpointCompatibility(lot models.Lot, tractor models.Tractor) bool {
	var currentRouteCheckpoint models.RouteCheckpoint
	var lotRouteCheckpoint models.RouteCheckpoint
	if err := currentRouteCheckpoint.GetRouteCheckpoint(LotController.Db, *tractor.RouteId, *tractor.CurrentCheckpointId); err != nil {
		return false
	}
	if lot.ResourceType != tractor.ResourceType {
		return false
	}
	if err := lotRouteCheckpoint.GetRouteCheckpoint(LotController.Db, *tractor.RouteId, *lot.CurrentCheckpointId); err != nil {
		return false
	}
	if currentRouteCheckpoint.Position > lotRouteCheckpoint.Position {
		return false
	}
	volumAtCheckpoint, err := tractor.GetVolumeAtCheckpoint(LotController.Db, *lot.StartCheckpointId)
	if err != nil {
		return false
	}
	var remainingVolume = tractor.MaxVolume - volumAtCheckpoint
	return remainingVolume >= lot.Volume
}

// AssignTractorToLot : Assign a tractor to a lot
// @Summary      Assign a tractor to a lot
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  body  string  true  "Lot Id"
// @Param        tractor_id  body  string  true  "Tractor Id"
// @Success      200  {object}  models.Lot
// @Failure      400  "Invalid request payload"
// @Failure      404  "Lot not found"
// @Failure      404  "Tractor not found"
// @Failure      500  "Unable to assign tractor to lot"
// @Router       /lots/assign [put]
func (LotController *LotController) AssignTractorToLot(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tractor not found"})
		return
	}

	if !LotController.checkCompatibility(lot, tractor) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lot is not compatible with the tractor"})
		return
	}

	var transactionIn models.Transaction
	var transactionOut models.Transaction

	var routeCheckpointStart models.RouteCheckpoint
	var routeCheckpointEnd models.RouteCheckpoint
	if err := routeCheckpointStart.GetRouteCheckpoint(LotController.Db, *tractor.RouteId, *lot.StartCheckpointId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := routeCheckpointEnd.GetRouteCheckpoint(LotController.Db, *tractor.RouteId, *lot.EndCheckpointId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := transactionIn.CreateTransaction(LotController.Db, models.TransactionState(models.TransactionStateIn), lot.Id, tractor.Id, *tractor.RouteId, *lot.StartCheckpointId, *lot.TrafficManagerId, routeCheckpointStart.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := transactionOut.CreateTransaction(LotController.Db, models.TransactionState(models.TransactionStateOut), lot.Id, tractor.Id, *tractor.RouteId, *lot.EndCheckpointId, *lot.TrafficManagerId, routeCheckpointEnd.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := lot.AssociateTractor(LotController.Db, tractor.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if lot.StartCheckpointId.String() == tractor.CurrentCheckpointId.String() {
		lot.InTractor = true
		lot.State = models.StateInTransit
		tractor.CurrentVolume += lot.Volume
		lot.Update(LotController.Db)
		tractor.Update(LotController.Db)
	}
	c.JSON(http.StatusOK, lot)
}

func (LotController *LotController) GetAvailableTrader(c *gin.Context) (models.User, error) {
	var user models.User
	traders, err := user.FindByRole(LotController.Db, models.RoleTrader)
	if err != nil {
		return models.User{}, err
	}

	var lots []models.Lot
	if err := LotController.Db.Where("state = ?", models.StateAtTrader).Find(&lots).Error; err != nil {
		return models.User{}, err
	}

	traderCounts := make(map[uuid.UUID]int)
	for _, lot := range lots {
		if lot.TraderId != nil {
			traderCounts[*lot.TraderId]++
		}
	}

	var selectedTrader models.User
	minCount := len(lots)

	for _, trader := range traders {
		count := traderCounts[trader.Id]
		if count <= minCount {
			minCount = count
			selectedTrader = trader
		}
	}

	return selectedTrader, nil
}

// AssignTraderToLot : Assign a trader to a lot
//
// @Summary      Assign a trader to a lot
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        lot_id  body  string  true  "Lot Id"
// @Param        trader_id  body  string  true  "Trader Id"
// @Success      200  {object}  models.Lot
// @Failure      400  "Invalid request payload"
// @Failure      404  "Lot not found"
// @Failure      404  "Trader not found"
// @Failure      500  "Unable to assign trader to lot"
// @Router       /lots/assign/{lot_id}/trader [post]
func (LotController *LotController) AssignTraderToLot(c *gin.Context) {
	lotId := c.Param("lot_id")
	lotIdUUID, errIdUUID := uuid.Parse(lotId)
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lot_id"})
		return
	}

	var lot models.Lot
	if err := LotController.Db.First(&lot, "id = ?", lotIdUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
		return
	}

	trader, err := LotController.GetAvailableTrader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lot.TraderId = &trader.Id
	lot.State = models.StateAtTrader
	if err := LotController.Db.Save(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var requestBody struct {
		Date string `json:"limit_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var offer models.Offer
	parsedDate, err := time.Parse(time.RFC3339, requestBody.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	_, err = offer.CreateOfferLot(LotController.Db, parsedDate, lot.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lot)
}

// GetAllLotTraderId : Get all lots with trader id
//
// @Summary      Get all lots with trader id
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        trader_id  path  string  true  "Trader Id"
// @Success      200  {array}  models.Lot
// @Failure      400  "Invalid trader_id"
// @Failure      404  "Trader not found"
// @Failure      500  "Unable to retrieve lots"
// @Router       /lots/trader/{trader_id} [get]
func (LotController *LotController) GetAllLotTraderId(c *gin.Context) {
	var lots []models.Lot
	traderId := c.Param("trader_id")
	traderIdUUID, errIdUUID := uuid.Parse(traderId)

	// Validate trader_id
	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trader_id"})
		return
	}

	// Check if the trader exists
	var trader models.User
	if err := LotController.Db.First(&trader, "id = ? AND role = ?", traderIdUUID, "trader").Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trader not found"})
		return
	}

	// Retrieve lots for the trader
	if err := LotController.Db.Preload("StartCheckpoint").
		Preload("EndCheckpoint").
		Where("trader_id = ?", traderIdUUID).
		Find(&lots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve lots", "details": err.Error()})
		return
	}

	// Enrich the lots with current prices and limit_date from associated offers
	for i := range lots {
		var maxBid float64
		var offer models.Offer
		if err := LotController.Db.First(&offer, "lot_id = ?", lots[i].Id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve offer", "details": err.Error()})
			return
		}
		// Get the maximum bid for the offer
		if err := LotController.Db.Raw("SELECT COALESCE(MAX(bid), 0) FROM bids WHERE offer_id = ?", offer.Id).Scan(&maxBid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve max bid", "details": err.Error()})
			return
		}
		lots[i].CurrentPrice = maxBid
		lots[i].LimitDate = offer.LimitDate // Set the limit_date from the offer
	}

	// Return the enriched lots in the JSON response
	c.JSON(http.StatusOK, lots)
}

// GetLotBidByOwnerId : Get all bids for a lot by owner id
//
// @Summary      Get all bids for a lot by owner id
// @Tags         lots
// @Accept       json
// @Produce      json
// @Param        client_id  path  string  true  "Client Id"
// @Success      200  {array}  models.Bid
// @Failure      400  "Invalid client_id"
// @Failure      500  "Unable to retrieve bids"
// @Router       /lots/bids/{client_id} [get]
func (LotController *LotController) GetLotBidByOwnerId(c *gin.Context) {
	var bids []models.Bid
	clientId := c.Param("owner_id")
	clientIdUUID, errIdUUID := uuid.Parse(clientId)

	if errIdUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner_id"})
		return
	}

	if err := LotController.Db.Joins("JOIN offers ON offers.id = bids.offer_id").Where("bids.owner_id = ? AND offers.lot_id IS NOT NULL", clientIdUUID).Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}
