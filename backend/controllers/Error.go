package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type HttpError struct {
	Code    int
	Message string
}

func (e *HttpError) Error(c *gin.Context) {
	codeString := strconv.Itoa(e.Code)
	c.JSON(e.Code, gin.H{codeString: e.Message})
}

func Err500(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}
