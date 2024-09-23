package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tms-backend/controllers"
)

func LotRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	LotController := controllers.LotController{
		Db: db,
	}

	v1 := r.Group("/api/v1/lots")
	{
		v1.POST("", LotController.AddLot)
		//v1.PATCH(":id", LotController.PatchLot)
		//v1.GET("", LotController.ListLots)
	}
	return r
}
