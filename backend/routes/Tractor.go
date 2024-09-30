package routes

import (
	"tms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TractorRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	TractorController := controllers.TractorController{
		Db: db,
	}
	v1 := r.Group("/api/v1/tractors")
	{
		v1.POST("/tractorId/trafficManagerId", TractorController.AddTrafficManager)
		v1.POST("/", TractorController.AddTractor)
		// Get tractors by OwnerID
		v1.GET("/owner/:ownerId", TractorController.GetTractorsByOwnerId)
		// Get tractors by TrafficManagerId
		v1.GET("/trafficManager/:trafficManagerId", TractorController.GetTractorsByTrafficManagerId)
		// Get tractors by State
		v1.GET("/state/:state", TractorController.GetTractorsByState)
		// Get tractors by RouteId
		v1.GET("/route/:routeId", TractorController.GetTractorsByRouteId)
		v1.GET("/next-route", TractorController.GoToNextCheckpoint)
		v1.PATCH("/updateState", TractorController.UpdateTractorState)
		//v1.PATCH(":id", LotController.PatchLot)
		//v1.GET("", LotController.ListLots)
	}
	return r
}
