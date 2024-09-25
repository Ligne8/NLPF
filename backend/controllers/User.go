package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"tms-backend/models"
)

type UserController struct {
	Db *gorm.DB
}

func (UserController *UserController) GetUsers(c *gin.Context) {
	var user models.User
	users, err := user.GetAllUsers(UserController.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (UserController *UserController) GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	userId, err := uuid.Parse(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	user, err = user.FindById(UserController.Db, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (UserController *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := UserController.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (UserController *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}
	user, err = user.FindById(UserController.Db, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateUser models.User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := UserController.Db.Model(&user).Updates(updateUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (UserController *UserController) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}
	user, err = user.FindById(UserController.Db, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := UserController.Db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
