package controllers

import (
	"net/http"
	"time"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		LimitDate string `json:"limit_date" binding:"required"`
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

	// Create the offer
	var offer models.Offer;
	parsedDate, err := time.Parse(time.RFC3339, requestBody.LimitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var offerId uuid.UUID;
	offerId, err = offer.CreateOfferLot(sec.Db, parsedDate, lot.Id);
	if err != nil {	
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lot.OfferId = &offerId;
	if err := sec.Db.Save(&lot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		LimitDate string `json:"limit_date" binding:"required"`
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

	// Create the offer
	var offer models.Offer;
	parsedDate, err := time.Parse(time.RFC3339, requestBody.LimitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var offerId uuid.UUID;
	offerId, err = offer.CreateOfferTractor(sec.Db, parsedDate, tractor.Id);
	if err != nil {	
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tractor.OfferId = &offerId;
	if err := sec.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var offers []struct {
		LimitDate      time.Time           `json:"limit_date"`
		TractorId      uuid.UUID           `json:"tractor_id"`
		ResourceType   models.ResourceType `json:"resource_type"`
		CurrentUnits   float64             `json:"current_units"`
		MaxUnits       float64             `json:"max_units"`
		MinPriceByKm   float64             `json:"min_price_by_km"`
		CurrentPrice   float64             `json:"current_price"`
	}

	query := `
		SELECT o.limit_date, t.id as tractor_id, t.resource_type, t.current_volume as current_units, t.max_volume as max_units, t.min_price_by_km, MAX(b.bid) as current_price
		FROM tractors t
		JOIN offers o ON t.id = o.tractor_id
		LEFT JOIN bids b ON o.id = b.offer_id
		WHERE o.limit_date > (SELECT simulation_date FROM simulations LIMIT 1)
		GROUP BY o.limit_date, t.id, t.resource_type, t.current_volume, t.max_volume, t.min_price_by_km
		ORDER BY o.limit_date
	`

	if err := sec.Db.Raw(query).Scan(&offers).Error; err != nil {
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
	
	var offers []struct {
		LimitDate      time.Time           `json:"limit_date"`
		LotId          uuid.UUID           `json:"lot_id"`
		ResourceType   models.ResourceType `json:"resource_type"`
		Volume         float64             `json:"volume"`
		MaxPriceByKm   float64             `json:"max_price_by_km"`
		CurrentPrice   float64             `json:"current_price"`
	}

	query := `
		SELECT o.limit_date, l.id as lot_id, l.resource_type, l.volume, l.max_price_by_km, MAX(b.bid) as current_price
		FROM lots l
		JOIN offers o ON l.id = o.lot_id
		LEFT JOIN bids b ON o.id = b.offer_id
		WHERE o.limit_date > (SELECT simulation_date FROM simulations LIMIT 1)
		GROUP BY o.limit_date, l.id, l.resource_type, l.volume, l.max_price_by_km
		ORDER BY o.limit_date
	`

	if err := sec.Db.Raw(query).Scan(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch offers"})
		return
	}


	c.JSON(http.StatusOK, offers)
}
