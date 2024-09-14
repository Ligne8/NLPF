package routes

import (
	"tms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CheckpointsRoute(r *gin.Engine, db *gorm.DB) *gin.Engine {

	CheckpointController := controllers.CheckpointController{
		Db: db,
	}

	v1 := r.Group("/api/v1/checkpoints")
	{
		v1.POST("", CheckpointController.AddCheckpoint)
	}
	return r
}