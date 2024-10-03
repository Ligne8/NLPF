package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"tms-backend/models"
)

type SimulationController struct {
	Db *gorm.DB
}

// GetSimulationDate Handler to get simulationDate
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

// UpdateSimulationDate Handler to increment the simulation date by 1 day
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

	// Return the new updated date
	c.JSON(http.StatusOK, gin.H{
		"message":         "Simulation date updated successfully",
		"simulation_date": newDate.Format("2006-01-02"),
	})
}
