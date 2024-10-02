package main

import (
	"time"
	"tms-backend/database"
	docs "tms-backend/docs"
	"tms-backend/models"
	"tms-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swagger embed files

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowHeaders:     []string{"Strict-Transport-Security", "strict-origin-when-cross-origin", "Content-Type"},
	}))

	db := database.InitDb()

	models.CreateCheckpoints(db)
	router = routes.CheckpointsRoute(router, db)
	router = routes.LotRoutes(router, db)
	router = routes.TractorRoutes(router, db)
	router = routes.UserRoutes(router, db)
	router = routes.AuthRoutes(router, db)

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(":8080")
}
