package routes

import (
	"tms-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesRoute(r *gin.Engine, db *gorm.DB) *gin.Engine {

	RouteController := controllers.RouteController{
		Db: db,
	}

	v1 := r.Group("/api/v1/routes")
	{
		v1.POST("", RouteController.CreateRoute)
		v1.GET("", RouteController.GetAllRoutes)
		v1.GET("/traffic_manager/parsed/:traffic_manager_id", RouteController.GetRouteStringByTrafficManagerId)

	}
	return r
}
