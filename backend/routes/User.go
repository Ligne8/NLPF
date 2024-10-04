package routes

import (
	"tms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	UserController := controllers.UserController{
		Db: db,
	}
	v1 := r.Group("/api/v1/users")
	{
		v1.GET("/", UserController.GetUsers)
		v1.GET("/:id", UserController.GetUser)
		v1.POST("/", UserController.CreateUser)
		v1.PATCH("/:id", UserController.UpdateUser)
		v1.DELETE("/:id", UserController.DeleteUser)
		v1.GET("/traffic_managers", UserController.GetTrafficManagers)
	}
	return r
}
