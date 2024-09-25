package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tms-backend/controllers"
)

func TractorRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	TractorController := controllers.TractorController{
		Db: db,
	}
	v1 := r.Group("/api/v1/tractors")
	{
		v1.POST("/tractorId/trafficManagerId", TractorController.AddTrafficManager)
		// Get tractors by OwnerID
		v1.GET("/owner/:ownerId", TractorController.GetTractorsByOwnerId)
		// Get tractors by TrafficManagerId
		v1.GET("/trafficManager/:trafficManagerId", TractorController.GetTractorsByTrafficManagerId)
		// Get tractors by State
		v1.GET("/state/:state", TractorController.GetTractorsByState)
		// Get tractors by RouteId
		v1.GET("/route/:routeId", TractorController.GetTractorsByRouteId)
	}
	return r
}
