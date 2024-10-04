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
	Position     uint   `json:"position" binding:"required"`
}

// GetAllRoutes : Get all routes
//
// @Summary      Get all routes
// @Tags         routes
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Route
// @Failure      400  "Unable to retrieve routes"
// @Router       /routes [get]
func (RouteController *RouteController) GetAllRoutes(c *gin.Context) {
	var routeModel models.Route
	var routes []models.Route
	routes, err := routeModel.GetAllRoutes(RouteController.Db)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, routes)
}

type routePayload struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Route_path string `json:"route_path"`
}

// GetRouteStringByTrafficManagerId : Get route string by traffic manager id
//
// @Summary      Get route string by traffic manager id
// @Tags         routes
// @Accept       json
// @Produce      json
// @Param        traffic_manager_id  path  string  true  "Traffic Manager Id"
// @Success      200  {array}   routePayload
// @Failure      400  "Unable to retrieve routes"
// @Router       /routes/traffic_manager/{traffic_manager_id} [get]
func (RouteController *RouteController) GetRouteStringByTrafficManagerId(c *gin.Context) {
	var allRoutes []routePayload
	var routeModel models.Route
	var routes []models.Route

	trafficManagerId := c.Param("traffic_manager_id")
	trafficManagerIdUUID := uuid.MustParse(trafficManagerId)
	routes, err := routeModel.GetRoutesByTrafficManagerId(RouteController.Db, trafficManagerIdUUID)

	for _, route := range routes {
		var path_string = route.GetRouteString(RouteController.Db)
		var route_payload = routePayload{
			Id:         route.Id.String(),
			Name:       route.Name,
			Route_path: path_string,
		}
		allRoutes = append(allRoutes, route_payload)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, allRoutes)
}

// GetRouteByTrafficManagerId : Get route by traffic manager id
//
// @Summary      Get route by traffic manager id
// @Tags         routes
// @Accept       json
// @Produce      json
// @Param        traffic_manager_id  path  string  true  "Traffic Manager Id"
// @Success      200  {array}   models.Route
// @Failure      400  "Unable to retrieve routes"
// @Router       /routes/traffic_manager/{traffic_manager_id} [get]
func (RouteController *RouteController) GetRouteByTrafficManagerId(c *gin.Context) {
	var routeModel models.Route
	var routes []models.Route
	trafficManagerId := c.Param("traffic_manager_id")
	trafficManagerIdUUID := uuid.MustParse(trafficManagerId)
	routes, err := routeModel.GetRoutesByTrafficManagerId(RouteController.Db, trafficManagerIdUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, routes)
}

// CreateRoute : Create a new Route
//
// @Summary      Create a new Route
// @Tags         routes
// @Accept       json
// @Produce      json
// @Param        name  body  string  true  "Name"
// @Param        traffic_manager_id  body  string  true  "Traffic Manager Id"
// @Param        route  body  []checkpointPosition  true  "Route"
// @Success      201  "Route created"
// @Failure      400  "Invalid request payload"
// @Router       /routes [post]
func (RouteController *RouteController) CreateRoute(c *gin.Context) {
	var requestBody struct {
		Name             string               `json:"name" binding:"required"`
		TrafficManagerId uuid.UUID            `json:"traffic_manager_id" binding:"required"`
		Route            []checkpointPosition `json:"route" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var routeModel models.Route
	routeModel.Id = uuid.New()
	routeModel.Name = requestBody.Name
	routeModel.TrafficManagerId = requestBody.TrafficManagerId
	routeModel.SaveRoute(RouteController.Db)
	for _, checkpoint := range requestBody.Route {
		var routeCheckpointModel models.RouteCheckpoint
		routeCheckpointModel.RouteId = routeModel.Id
		var checkpointId uuid.UUID = uuid.MustParse(checkpoint.CheckpointId)
		routeCheckpointModel.CheckpointId = checkpointId
		routeCheckpointModel.Position = checkpoint.Position
		routeCheckpointModel.SaveRouteCheckpoint(RouteController.Db)
	}
	c.Status(http.StatusCreated)
}
