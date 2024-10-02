package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tms-backend/controllers/Authentification"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {
	AuthController := controllers.AuthController{
		Db: db,
	}
	v1 := r.Group("/api/v1/auth")
	{
		v1.POST("/register", AuthController.CreateUser)
		v1.POST("/login", AuthController.LoginUser)
	}
	return r
}
