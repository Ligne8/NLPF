package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type CheckpointController struct {
	Db *gorm.DB
}

// AddCheckpoint gère la création d'un checkpoint
// @version 1
// @Summary Add checkpoint
// @Tags Checkpoint
// @Accept json
// @Param   name 	 body   string     true        "Name of the city" example(Paris)
// @Param   country 	body    string     true        "Name of the country" example(France)
// @Produce json
// @Failure 500 "Error from gin"
// @Success 201 "No content"
// @Router /checkpoints [post]
func (CheckpointController * CheckpointController) AddCheckpoint(c *gin.Context) {
	var checkpointModel models.Checkpoint

	if err := c.ShouldBindJSON(&checkpointModel); err != nil {
		Err500(c, err)
		return
	}
	if err := checkpointModel.Save(CheckpointController.Db); err != nil {
		Err500(c, err)
		return
	}
	c.Status(http.StatusCreated)
}