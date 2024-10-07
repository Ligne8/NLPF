package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type HttpError struct {
	Code    int
	Message string
}

// Error function to handle error
//
// @Summary Error handling
// @Description Handle error
// @Tags error
// @Accept json
// @Produce json
// @Param code path int true "Code"
// @Param message path string true "Message"
// @Success 200 {string} string
// @Router /error/{code}/{message} [get]
// @Failure 500 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
func (e *HttpError) Error(c *gin.Context) {
	codeString := strconv.Itoa(e.Code)
	c.JSON(e.Code, gin.H{codeString: e.Message})
}

func Err500(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}
