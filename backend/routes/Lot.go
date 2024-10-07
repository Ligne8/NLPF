package routes

import (
	"tms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LotRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	LotController := controllers.LotController{
		Db: db,
	}

	v1 := r.Group("/api/v1/lots")
	{
		v1.POST("", LotController.CreateLot)
		v1.POST("traffic_manager", LotController.AssociateToTrafficManager)
		v1.PATCH("/state", LotController.UpdateLotState)
		//v1.PATCH(":id", LotController.PatchLot)
		v1.GET("owner/:owner_id", LotController.ListLotsByOwner)
		v1.DELETE("/:lot_id", LotController.DeleteLot)

		v1.GET("traffic_manager/:traffic_manager_id", LotController.ListLotsByTrafficManager)
		v1.GET("/tractors/compatible/:traffic_manager_id/:lot_id", LotController.ListCompatibleTractorsForLot)
		v1.PUT("/assign/tractor", LotController.AssignTractorToLot)
	}
	return r
}
