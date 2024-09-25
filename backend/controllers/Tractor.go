package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"tms-backend/models"
)

type TractorController struct {
	Db *gorm.DB
}

func (TractorController *TractorController) AddTractor(c *gin.Context) {
	var requestBody struct {
		ResourceType  models.ResourceType `json:"resource_type" binding:"required"`
		MaxVolume     float64             `json:"max_units"`
		CurrentVolume float64             `json:"current_units"`
		//CurrentCheckpointId uuid.UUID    `json:"current_checkpoint_id" gorm:"type:uuid"` // Foreign key for Checkpoint
		//CurrentCheckpoint   Checkpoint   `json:"current_checkpoint" gorm:"foreignKey:CurrentCheckpointId"`
		State   models.State `json:"state"`
		OwnerId uuid.UUID    `json:"owner_id"` // Foreign key for User
		//Owner            models.User  `json:"owner"`
		MinPriceByKm     uint      `json:"min_price_by_km"`
		TrafficManagerId uuid.UUID `json:"traffic_manager_id"` // Foreign key for User
		//TrafficManager   models.User  `json:"traffic_manager"`
		//TraderId         uuid.UUID    `json:"trader_id"` // Foreign key for User
		//Trader           models.User  `json:"trader" gorm:"foreignKey:TraderId"`
		//RouteId          uuid.UUID    `json:"route_id" gorm:"type:uuid"` // Foreign key for Route
		//Route            models.Route `json:"route" gorm:"foreignKey:RouteId"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	TractorModel := models.Tractor{
		ResourceType:  requestBody.ResourceType,
		MaxVolume:     requestBody.MaxVolume,
		CurrentVolume: requestBody.CurrentVolume,
		//CurrentCheckpointId: uuid.UUID{},
		//CurrentCheckpoint:   models.Checkpoint{},
		State:   requestBody.State,
		OwnerId: requestBody.OwnerId,
		//Owner:               models.User{},
		MinPriceByKm:     requestBody.MinPriceByKm,
		TrafficManagerId: requestBody.TrafficManagerId,
		//TrafficManager:      models.User{},
		//TraderId:            uuid.UUID{},
		//Trader:              models.User{},
		//RouteId:             uuid.UUID{},
		//Route:               models.Route{},
	}

	if err := TractorController.Db.Create(&TractorModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TractorModel)
}

func (TractorController *TractorController) AddTrafficManager(c *gin.Context) {
	var requestBody struct {
		TractorId        uuid.UUID `json:"tractor_id" binding:"required"`
		TrafficManagerId uuid.UUID `json:"traffic_manager_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var trafficManager models.User
	if err := TractorController.Db.First(&trafficManager, "id = ?", requestBody.TrafficManagerId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "traffic manager not found when trying to add it to the tractor"})
		return
	}

	if trafficManager.Role == models.RoleTrafficManager {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The user is not a traffic manager"})
		return
	}

	var tractor models.Tractor
	if err := TractorController.Db.First(&tractor, "id = ?", requestBody.TractorId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tractor not found when trying to add a traffic manager to it"})
		return
	}

	tractor.TrafficManagerId = requestBody.TrafficManagerId

	// C le bon state ?
	tractor.State = models.StateInTransit

	//Il faut ajouter le traffic manager au tracteur direct comme ca ? pq pas juste utiliser la cle etrangere ?
	//tractor.TrafficManager = &trafficManager

	if err := TractorController.Db.Save(&tractor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (TractorController *TractorController) GetTractorsByOwnerId(c *gin.Context) {
	ownerId := c.Param("ownerId")
	var tractors []models.Tractor

	if _, err := uuid.Parse(ownerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ownerId"})
		return
	}

	if err := TractorController.Db.Where("owner_id = ?", ownerId).Find(&tractors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) GetTractorsByTrafficManagerId(c *gin.Context) {
	trafficManagerId := c.Param("trafficManagerId")
	var tractors []models.Tractor

	if _, err := uuid.Parse(trafficManagerId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trafficManagerId"})
		return
	}

	if err := TractorController.Db.Where("traffic_manager_id = ?", trafficManagerId).Find(&tractors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) GetTractorsByState(c *gin.Context) {
	state := c.Param("state")
	var tractors []models.Tractor

	if err := TractorController.Db.Where("state = ?", state).Find(&tractors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tractors)
}

func (TractorController *TractorController) GetTractorsByRouteId(c *gin.Context) {
	routeId := c.Param("routeId")
	var tractors []models.Tractor

	if _, err := uuid.Parse(routeId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid routeId"})
		return
	}

	if err := TractorController.Db.Where("route_id = ?", routeId).Find(&tractors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tractors)
}
