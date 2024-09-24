package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tms-backend/controllers"
)

func CheckpointsRoute(r *gin.Engine, db *gorm.DB) *gin.Engine {
	CheckpointController := controllers.CheckpointController{
		Db: db,
	}

	v1 := r.Group("/api/v1/checkpoints")
	{
		v1.GET("", CheckpointController.GetAllCheckpoints)
		v1.GET("/countries/:country/cities", CheckpointController.GetCitiesByCountry)
		v1.GET("/cities/:city/country", CheckpointController.GetCountryByCity)
	}
	return r
}
