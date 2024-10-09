package controllers

import (
	"log"
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
		LimitDate string    `json:"limit_date" binding:"required"`
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
	var offer models.Offer
	parsedDate, err := time.Parse(time.RFC3339, requestBody.LimitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var offerId uuid.UUID
	offerId, err = offer.CreateOfferLot(sec.Db, parsedDate, lot.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lot.OfferId = &offerId
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
		LimitDate string    `json:"limit_date" binding:"required"`
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
	var offer models.Offer
	parsedDate, err := time.Parse(time.RFC3339, requestBody.LimitDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	var offerId uuid.UUID
	offerId, err = offer.CreateOfferTractor(sec.Db, parsedDate, tractor.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tractor.OfferId = &offerId
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
		LimitDate    time.Time           `json:"limit_date"`
		TractorId    uuid.UUID           `json:"tractor_id"`
		ResourceType models.ResourceType `json:"resource_type"`
		CurrentUnits float64             `json:"current_units"`
		MaxUnits     float64             `json:"max_units"`
		MinPriceByKm float64             `json:"min_price_by_km"`
		CurrentPrice float64             `json:"current_price"`
		OfferId      uuid.UUID           `json:"offer_id"`
	}

	query := `
		SELECT o.id as offer_id, o.limit_date, t.id as tractor_id, t.resource_type, t.current_volume as current_units, t.max_volume as max_units, t.min_price_by_km, MAX(b.bid) as current_price
		FROM tractors t
		JOIN offers o ON t.id = o.tractor_id
		LEFT JOIN bids b ON o.id = b.offer_id
		WHERE o.limit_date > (SELECT simulation_date FROM simulations LIMIT 1)
		GROUP BY o.id, o.limit_date, t.id, t.resource_type, t.current_volume, t.max_volume, t.min_price_by_km
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
		LimitDate    time.Time           `json:"limit_date"`
		LotId        uuid.UUID           `json:"lot_id"`
		ResourceType models.ResourceType `json:"resource_type"`
		Volume       float64             `json:"volume"`
		MaxPriceByKm float64             `json:"max_price_by_km"`
		CurrentPrice float64             `json:"current_price"`
		OfferId      uuid.UUID           `json:"offer_id"`
	}

	query := `
		SELECT o.id as offer_id, o.limit_date, l.id as lot_id, l.resource_type, l.volume, l.max_price_by_km, MIN(b.bid) as current_price
		FROM lots l
		JOIN offers o ON l.id = o.lot_id
		LEFT JOIN bids b ON o.id = b.offer_id
		WHERE o.limit_date > (SELECT simulation_date FROM simulations LIMIT 1)
		GROUP BY o.id, o.limit_date, l.id, l.resource_type, l.volume, l.max_price_by_km
		ORDER BY o.limit_date
	`

	if err := sec.Db.Raw(query).Scan(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch offers"})
		return
	}

	c.JSON(http.StatusOK, offers)
}

func (StockExchangeController *StockExchangeController) CreateBidLot(c *gin.Context) {
	var requestBody struct {
		Bid     float64   `json:"bid" binding:"required"`
		OfferId uuid.UUID `json:"offer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the UUID of the offer ID
	offerUUID, err := uuid.Parse(requestBody.OfferId.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID format"})
		return
	}

	var bid models.Bid
	bid.Bid = requestBody.Bid
	bid.OfferId = offerUUID
	bid.State = "in_progress"

	if err := StockExchangeController.Db.Create(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorrr": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

func (StockExchangeController *StockExchangeController) CreateBidTractor(c *gin.Context) {
	var requestBody struct {
		Bid     float64   `json:"bid" binding:"required"`
		OfferId uuid.UUID `json:"offer_id" binding:"required"`
		Volume  float64   `json:"volume" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the UUID of the offer ID
	offerUUID, err := uuid.Parse(requestBody.OfferId.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID format"})
		return
	}

	var bid models.Bid
	bid.Bid = requestBody.Bid
	bid.OfferId = offerUUID
	bid.State = "in_progress"
	bid.Volume = requestBody.Volume

	if err := StockExchangeController.Db.Create(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorrr": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

// ReturnFromMarket : Return a tractor or a lot from the market Deso on utilise plus la fonction daniel
// @Summary After getting all the offers on morket with the limit date passed, the state of the tractor/lot is changed to return_from_market
// @Tags Stock Exchange
// @Produce json
// @Success 200 {string} string "Tractors and lots returned from market"
// @Failure 500 "Unable to fetch simulation date"
// @Failure 500 "Unable to fetch offers"
// @Router /stock_exchange/return_from_market [put]
func (sec *StockExchangeController) ChangeStateToReturnFromMarket(c *gin.Context) {
	// get offers with the state "on_market" and the LimitDate higher than the current date
	var simulation models.Simulation
	if err := sec.Db.First(&simulation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch simulation date"})
		return
	}
	log.Println("similationDate:", simulation.SimulationDate)
	var offers []models.Offer
	if err := sec.Db.Joins("LEFT JOIN tractors ON offers.tractor_id = tractors.id").
		Joins("LEFT JOIN lots ON offers.lot_id = lots.id").
		Where("(tractors.state = ? OR lots.state = ?) AND offers.limit_date <= ?", models.StateOnMarket, models.StateOnMarket, simulation.SimulationDate).
		Find(&offers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch offers"})
		return
	}

	// for each offer change the state of the tractor/lot to "return_from_market"
	for _, offer := range offers {
		if offer.TractorId != nil {
			var tractor models.Tractor
			if err := sec.Db.First(&tractor, "id = ?", offer.TractorId).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch tractor"})
				return
			}
			tractor.State = models.StateReturnFromMarket
			if err := tractor.Save(sec.Db); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update tractor state"})
				return
			}
		} else if offer.LotId != nil {
			var lot models.Lot
			if err := sec.Db.First(&lot, "id = ?", offer.LotId).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch lot"})
				return
			}
			lot.State = models.StateReturnFromMarket
			if err := lot.Save(sec.Db); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update lot state"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tractors and lots returned from market"})
}

// ReturnFromMarket : Return a tractor or a lot from the market
// @Summary After getting all the offers on morket with the limit date passed, the state of the tractor/lot is changed to return_from_market
// @Tags Stock Exchange
// @Produce json
// @Success 200 {string} string "Tractors and lots returned from market"
// @Failure 500 "Unable to fetch simulation date"
// @Failure 500 "Unable to fetch offers"
// @Router /stock_exchange/return_from_market [put]
func (sec *StockExchangeController) ChangeStateToReturnFromMarket2(c *gin.Context) {

	if err := sec.UpdateLotsBids(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := sec.UpdateTractorsBids(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := sec.updateLotsOffers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := sec.updateTractorsOffers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (sec *StockExchangeController) UpdateLotsBids() error {
	var offers []models.Offer
	query := `
		SELECT offers.*
		FROM offers
		JOIN lots l ON offers.id = l.offer_id
		WHERE offers.limit_date < (SELECT simulation_date FROM simulations LIMIT 1)
		AND offers.lot_id IS NOT NULL AND offers.tractor_id IS NULL
		AND l.state = 'on_market'
	`

	if err := sec.Db.Raw(query).Scan(&offers).Error; err != nil {
		return err
	}

	for _, offer := range offers {
		var bids []models.Bid
		if err := sec.Db.Where("offer_id = ?", offer.Id).Order("bid asc").Find(&bids).Error; err != nil {
			return err
		}
		for i, bid := range bids {
			if i == 0 {
				bid.State = "accepted"
			} else {
				bid.State = "rejected"
			}
			if err := sec.Db.Save(&bid).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (sec *StockExchangeController) UpdateTractorsBids() error {
	var offers []models.Offer
	query := `
		SELECT offers.*
		FROM offers
		JOIN tractors t ON offers.id = t.offer_id
		WHERE offers.limit_date < (SELECT simulation_date FROM simulations LIMIT 1)
		AND offers.tractor_id IS NOT NULL AND offers.lot_id IS NULL
		AND t.state = 'on_market'
	`
	if err := sec.Db.Raw(query).Scan(&offers).Error; err != nil {
		return err
	}
	for _, offer := range offers {
		var bids []models.Bid
		if err := sec.Db.Where("offer_id = ?", offer.Id).Order("bid desc").Find(&bids).Error; err != nil {
			return err
		}
		for _, bid := range bids {
			if bid.State != "in_progress" {
				continue
			}
			var tractor models.Tractor
			tractor, err := tractor.FindById(sec.Db, *offer.TractorId)
			if err != nil {
				return err
			}
			if tractor.MaxVolume-tractor.CurrentVolume < bid.Volume {
				bid.State = "rejected"
				if err := sec.Db.Save(&bid).Error; err != nil {
					return err
				}
				continue
			}
			bid.State = "accepted"
			if err := sec.Db.Save(&bid).Error; err != nil {
				return err
			}
			tractor.CurrentVolume += bid.Volume
			if err := sec.Db.Save(&tractor).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (sec *StockExchangeController) updateLotsOffers() error {
	query := `
		WITH random_trafic_manager AS (
			SELECT id
			FROM users
			WHERE role = 'traffic_manager'
			ORDER BY RANDOM()
			LIMIT 1
		)
		UPDATE lots
		SET
			state = 'return_from_market'
		FROM offers
		WHERE offers.lot_id = lots.id AND offers.limit_date <= (SELECT simulation_date FROM simulations LIMIT 1) AND lots.state = 'on_market'
	`
	if err := sec.Db.Exec(query).Error; err != nil {
		return err
	}
	return nil
}
func (sec *StockExchangeController) updateTractorsOffers() error {
	query := `
		WITH random_trafic_manager AS (
			SELECT id
			FROM users
			WHERE role = 'traffic_manager'
			ORDER BY RANDOM()
			LIMIT 1
		)
		UPDATE tractors
		SET
			state = 'return_from_market'
		FROM offers
		WHERE offers.tractor_id = tractors.id AND offers.limit_date <= (SELECT simulation_date FROM simulations LIMIT 1) AND tractors.state = 'on_market'
	`
	if err := sec.Db.Exec(query).Error; err != nil {
		return err
	}
	return nil
}
