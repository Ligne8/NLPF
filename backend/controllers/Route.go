package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteController struct {
	Db *gorm.DB
}

type checkpointPosition struct {
	CheckpointId string `json:"checkpoint_id" binding:"required"`
	Position      uint `json:"position" binding:"required"`
}

func (RouteController *RouteController) GetAllRoutes(c *gin.Context) {
	var routeModel models.Route;
	var routes []models.Route;
	routes, err := routeModel.GetAllRoutes(RouteController.Db);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()});
	}
	c.JSON(http.StatusOK, routes);
}

// CreateRoute Create a new Route
func (RouteController *RouteController) CreateRoute(c *gin.Context) {
	var requestBody struct {
		Name string `json:"name" binding:"required"`
		TrafficManagerId     uuid.UUID `json:"traffic_manager_id" binding:"required"`
		Route []checkpointPosition `json:"route" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var routeModel models.Route;
	routeModel.Id = uuid.New()
	routeModel.Name = requestBody.Name
	routeModel.TrafficManagerId = requestBody.TrafficManagerId
	routeModel.SaveRoute(RouteController.Db)
	for _, checkpoint := range requestBody.Route {
		var routeCheckpointModel models.RouteCheckpoint;
		routeCheckpointModel.RouteId = routeModel.Id;
		var checkpointId uuid.UUID = uuid.MustParse(checkpoint.CheckpointId);
		routeCheckpointModel.CheckpointId = checkpointId
		routeCheckpointModel.Position = checkpoint.Position
		routeCheckpointModel.SaveRouteCheckpoint(RouteController.Db)
	}
	c.Status(http.StatusCreated)
}
