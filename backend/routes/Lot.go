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
		v1.PATCH("/state", LotController.UpdateLotState) 
		//v1.PATCH(":id", LotController.PatchLot)
		v1.GET("owner/:owner_id", LotController.ListLotsByOwner)

	}
	return r
}
