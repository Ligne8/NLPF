package controllers

import (
	"net/http"
	"tms-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB
}

// GetUsers : Get users
//
// @Summary      Get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  "Unable to retrieve users"
// @Router       /users [get]
func (UserController *UserController) GetUsers(c *gin.Context) {
	var user models.User
	users, err := user.GetAllUsers(UserController.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser : Get user
// @Summary      Get user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Success      200  {object}   models.User
// @Failure      500  "Unable to retrieve user"
// @Failure      400  "Invalid request"
// @Router       /user/{id} [get]
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

// CreateUser : Create a user
// @Summary      Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  models.User  true  "User"
// @Success      200  {object}   models.User
// @Failure      500  "Unable to create user"
// @Failure      400  "Invalid request"
// @Router       /user [post]
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

// UpdateUser : Update a user
// @Summary      Update user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Param        user  body  models.User  true  "User"
// @Success      200  {object}   models.User
// @Failure      500  "Unable to update user"
// @Failure      400  "Invalid request"
// @Router       /user/{id} [put]
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

// DeleteUser : Detele a user
// @Summary      Delete user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Success      200  {object}   models.User
// @Failure      500  "Unable to delete user"
// @Failure      400  "Invalid request"
// @Router       /user/{id} [delete]
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

// GetTrafficManagers : Get traffic managers
// @Summary      Get traffic managers
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  "Unable to retrieve traffic managers"
// @Router 	 	 /users/traffic_managers [get]
func (UserController *UserController) GetTrafficManagers(c *gin.Context) {
	var user models.User
	role := models.RoleTrafficManager
	users, err := user.FindByRole(UserController.Db, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
