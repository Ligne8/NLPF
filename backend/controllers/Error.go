package controllers

import "github.com/gin-gonic/gin"

func Err500(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}