package routes

import (
	"tms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StockExchangeRoute(r *gin.Engine, db *gorm.DB) *gin.Engine {
	StockExchangeController := controllers.StockExchangeController{
		Db: db,
	}

	v1 := r.Group("/api/v1/stock_exchange")
	{
		v1.POST("/lot_offers", StockExchangeController.CreateLotOffer)
		v1.GET("/lot_offers", StockExchangeController.GetAllLotsOnMarket)

		v1.POST("/lot/bid", StockExchangeController.CreateBidLot)
		v1.POST("/tractor/bid", StockExchangeController.CreateBidTractor)

		v1.POST("/tractor_offers", StockExchangeController.CreateTractorOffer)
		v1.GET("/tractor_offers", StockExchangeController.GetAllTractorOnMarket)
	}
	return r
}
