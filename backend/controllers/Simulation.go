package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SimulationController struct {
	Db *gorm.DB
}

// GetSimulationDate : Handler to get the simulation date
//
// @Summary      Get simulation date
// @Description  get the current simulation date
// @Tags         simulation
// @Accept       json
// @Produce      json
// @Success      200  {object}   string
// @Failure      500  "Unable to fetch simulation date"
// @Router       /simulation/date [get]
func (SimulationController *SimulationController) GetSimulationDate(c *gin.Context) {
	var simulation models.Simulation

	// Get simulation date from database
	if err := SimulationController.Db.First(&simulation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch simulation date"})
		return
	}

	// Format day/month/year
	c.JSON(http.StatusOK, gin.H{"simulation_date": simulation.SimulationDate.Format("2006-01-02")})
}

// UpdateSimulationDate : Handler to increment the simulation date by 1 day
//
// @Summary      Update simulation date
// @Description  increment the simulation date by 1 day
// @Tags         simulation
// @Accept       json
// @Produce      json
// @Success      200  {object}   string
// @Failure      500  "Unable to update simulation date"
// @Router       /simulation/date [put]
func (SimulationController *SimulationController) UpdateSimulationDate(c *gin.Context) {
	var simulation models.Simulation

	// Get the current simulation date from the database
	if err := SimulationController.Db.First(&simulation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch simulation date"})
		return
	}

	// Increase simulation date by one day
	newDate := simulation.SimulationDate.AddDate(0, 0, 1)

	// Update simulation date in database with ID-based WHERE condition
	if err := SimulationController.Db.Model(&simulation).Where("id = ?", simulation.ID).Update("simulation_date", newDate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update simulation date"})
		return
	}

	// Update the state of the offers
	stockExchangeController := StockExchangeController{Db: SimulationController.Db}
	stockExchangeController.ChangeStateToReturnFromMarket(c)

	// Return the new updated date
	c.JSON(http.StatusOK, gin.H{
		"message":         "Simulation date updated successfully",
		"simulation_date": newDate.Format("2006-01-02"),
	})
}
