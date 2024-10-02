package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tms-backend/controllers"
)

func SimulationRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	SimulationController := controllers.SimulationController{
		Db: db,
	}

	// Grouper les routes sous /api/v1/simulations
	v1 := r.Group("/api/v1/simulations")
	{
		v1.GET("/date", SimulationController.GetSimulationDate)

		v1.PATCH("/date", SimulationController.UpdateSimulationDate)

	}

	return r
}
