package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
	"tms-backend/models"
)

type StockExchangeController struct {
	Db *gorm.DB
}

// CreateLotOffer creates an offer for a lot
//
// @Summary Create an offer for a lot
// @Tags Stock Exchange
// @Accept json
// @Produce json
// @Param       LimitDate body time.Time true "Limit date"
// @Param       LotID body uuid.UUID true "Lot ID"
// @Success 201 {object} models.Offer
// @Failure 400 "Invalid request body"
// @Failure 404 "Lot not found"
// @Failure 500 "Unable to create offer"
// @Router /stock_exchange/lot_offers [post]
func (sec *StockExchangeController) CreateLotOffer(c *gin.Context) {
	var requestBody struct {
		LimitDate time.Time `json:"limit_date" binding:"required"`
		LotId     uuid.UUID `json:"lot_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the lot exists
	var lot models.Lot
	if result := sec.Db.First(&lot, "id = ?", requestBody.LotId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lot not found"})
		return
	}

	// Update the state of the lot
	lot.State = models.StateOnMarket
	if err := sec.Db.Save(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update lot state"})
		return
	}

	var simulation models.Simulation
	if err := sec.Db.First(&simulation); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to fetch simulation date"})
		return
	}
	// Create the offer
	offer := models.Offer{
		LimitDate: requestBody.LimitDate,
		LotId:     &requestBody.LotId,
		CreatedAt: simulation.SimulationDate,
	}

	if err := sec.Db.Create(&offer); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

// CreateTractorOffer creates an offer for a tractor
//
// @Summary Create an offer for a tractor
// @Tags Stock Exchange
// @Accept json
// @Produce json
// @Param       LimitDate body time.Time true "Limit date"
// @Param       TractorId body uuid.UUID true "Tractor ID"
// @Success 201 {object} models.Offer
// @Failure 400 "Invalid request body"
// @Failure 404 "Tractor not found"
// @Failure 500 "Unable to create offer"
// @Router /stock_exchange/tractor_offers [post]
func (sec *StockExchangeController) CreateTractorOffer(c *gin.Context) {
	var requestBody struct {
		LimitDate time.Time `json:"limit_date" binding:"required"`
		TractorId uuid.UUID `json:"tractor_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the lot exists
	var tractor models.Tractor
	if result := sec.Db.First(&tractor, "id = ?", requestBody.TractorId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found"})
		return
	}

	tractor.State = models.StateOnMarket
	if err := sec.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update tractor state"})
		return
	}

	var simulation models.Simulation
	if err := sec.Db.First(&simulation); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to fetch simulation date"})
		return
	}
	// Create the offer
	offer := models.Offer{
		LimitDate: requestBody.LimitDate,
		TractorId: &requestBody.TractorId,
		CreatedAt: simulation.SimulationDate,
	}

	if err := sec.Db.Create(&offer); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create offer"})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

// GetAllTractorOnMarket returns all tractor offers
//
// @Summary Get all tractor offers
// @Tags Stock Exchange
// @Produce json
// @Success 200 {array} models.Offer
// @Failure 500 "Unable to fetch offers"
// @Router /stock_exchange/tractor_offers [get]
func (sec *StockExchangeController) GetAllTractorOnMarket(c *gin.Context) {
	var offers []models.Offer
	if err := sec.Db.Preload("Tractor").Where("tractor_id IS NOT NULL").Find(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch offers"})
		return
	}

	c.JSON(http.StatusOK, offers)
}

// GetAllLotsOnMarket returns all lot offers
//
// @Summary Get all lot offers
// @Tags Stock Exchange
// @Produce json
// @Success 200 {array} models.Offer
// @Failure 500 "Unable to fetch offers"
// @Router /stock_exchange/lot_offers [get]
func (sec *StockExchangeController) GetAllLotsOnMarket(c *gin.Context) {
	var offers []models.Offer
	if err := sec.Db.Preload("Lot").Where("lot_id IS NOT NULL").Find(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch offers"})
		return
	}

	c.JSON(http.StatusOK, offers)
}
