package main

import (
	"tms-backend/database"
	"tms-backend/routes"

	docs "tms-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swagger embed files

func main() {
	router := gin.Default()

	db := database.InitDb()

	/*router = routes.CheckpointsRoute(router, db)*/
	router = routes.LotRoutes(router, db)

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(":8080")
}
