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
		//v1.PATCH("/updateState", TractorController.UpdateTractorState)
		//v1.PATCH(":id", LotController.PatchLot)
		//v1.GET("", LotController.ListLots)
	}
	return r
}
