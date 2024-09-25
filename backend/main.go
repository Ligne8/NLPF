package main

import (
	"tms-backend/database"
	docs "tms-backend/docs"
	"tms-backend/models"
	"tms-backend/routes"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swagger embed files

func main() {
	router := gin.Default()

	db := database.InitDb()

	models.CreateCheckpoints(db)
	router = routes.CheckpointsRoute(router, db)
	router = routes.LotRoutes(router, db)
	router = routes.TractorRoutes(router, db)
	router = routes.UserRoutes(router, db)

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(":8080")
}
