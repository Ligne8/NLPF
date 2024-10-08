package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Return the new updated date
	c.JSON(http.StatusOK, gin.H{
		"message":         "Simulation date updated successfully",
		"simulation_date": newDate.Format("2006-01-02"),
	})
}

func (SimulationController *SimulationController) MoveTractorForward(c *gin.Context){
	var tractorModel models.Tractor;
	var tractors []models.Tractor
	var err error;
	tractors, err = tractorModel.GetByState(SimulationController.Db, "in_transit")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch tractors"})
		return
	}
  for _, tractor := range tractors {

    // get les trucs a get
    var currentRouteCheckpoint models.RouteCheckpoint;
    if err := currentRouteCheckpoint.GetRouteCheckpoint(SimulationController.Db, *tractor.RouteId, *tractor.CurrentCheckpointId); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch route checkpoint"});
      return;
    }
    var nextRouteCheckpoint models.RouteCheckpoint;
    if err := nextRouteCheckpoint.GetNextCheckpoint(SimulationController.Db, *tractor.RouteId, currentRouteCheckpoint.Position); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch next route checkpoint"});
      return;
    }

    UpdateTractorCheckpoint(SimulationController.Db, c, tractor, currentRouteCheckpoint, nextRouteCheckpoint);
    UpdateLotCheckpoint(SimulationController.Db, c, tractor.Id, nextRouteCheckpoint.CheckpointId);
    ExecAllTransactions(SimulationController.Db, nextRouteCheckpoint.CheckpointId, tractor.Id, *tractor.RouteId, c);

    // update lot checkpoint
    
  }
}

func UpdateTractorCheckpoint(db *gorm.DB, c *gin.Context, tractor models.Tractor, currentRouteCheckpoint models.RouteCheckpoint, nextRouteCheckpoint models.RouteCheckpoint){
  tractor.CurrentCheckpointId = &nextRouteCheckpoint.CheckpointId;
  var lastCheckpointPosition uint; 
    db.Raw("select max(position) from route_checkpoints where route_id = ?", tractor.RouteId).Scan(&lastCheckpointPosition);
    if nextRouteCheckpoint.Position == lastCheckpointPosition {
        tractor.State = models.StateArchive
    }
    if err := tractor.Update(db); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save tractor"});
      return;
    }
}

func UpdateLotCheckpoint(db *gorm.DB, c *gin.Context, tractorId uuid.UUID, newCheckpointId uuid.UUID){
  var lot models.Lot;
  var lots []models.Lot;
  var err error;
  lots, err = lot.GetAllInTractorByTracorId(db, tractorId);
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch lots"});
    return;
  }
  for _, lot := range lots {
    if lot.InTractor{
      lot.CurrentCheckpointId = &newCheckpointId;
      if err := lot.Update(db); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save lot"});
        return;
      }
    }
  }
}

func ExecAllTransactions(db *gorm.DB, checkpointId uuid.UUID, tractorId uuid.UUID, routeId uuid.UUID, c *gin.Context){
  var transactionModel models.Transaction;
  var transactions []models.Transaction;
  var err error;
  transactions, err = transactionModel.FindByRouteIdAndCheckpointIdAndTractorId(db, routeId, checkpointId, tractorId);
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch transactions"});
  }
  for _, transaction := range transactions {
    if err := transaction.ExecTransaction(db); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to execute transaction"});
    }
  } 
}