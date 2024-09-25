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

// GetUsers Retrieve users
//
//		@Summary      List users
//		@Description  get all users
//		@Tags         users
//		@Accept       json
//		@Produce      json
//		@Success      200  {array}   models.User
//	 	@Failure   	  500 "Unable to retrieve users"
//		@Router       /users [get]
func (UserController *UserController) GetUsers(c *gin.Context) {
	var user models.User
	users, err := user.GetAllUsers(UserController.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser Retrieve user
//
//		@Summary      user
//		@Description  get a user by id
//		@Tags         users
//		@Accept       json
//		@Produce      json
//		@Success      200 {object} models.User
//	 	@Failure   	  500 "Unable to retrieve user"
//		@Failure   	  400 "Invalid userId"
//		@Router       /user/:id [get]
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

// CreateUser Retrieve user
//
//		@Summary      user
//		@Description  create a user
//		@Tags         users
//		@Accept       json
//		@Produce      json
//		@Param        user  body      models.User  true  "User"
//		@Success      200 {object} models.User
//	 	@Failure   	  500 "Unable to retrieve user"
//		@Failure   	  400 "Invalid request"
//		@Router       /user [post]
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

// UpdateUser Retrieve user
//
//		@Summary      user
//		@Description  updates a user
//		@Tags         users
//		@Accept       json
//		@Produce      json
//		@Param        user  body models.User  true  "User"
//		@Success      200 {object} models.User
//	 	@Failure   	  500 "Unable to retrieve user"
//		@Failure   	  400 "Invalid request"
//		@Router       /user/:id [patch]
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

// DeleteUser Retrieve user
//
//		@Summary      user
//		@Description  deletes a user
//		@Tags         users
//		@Accept       json
//		@Produce      json
//		@Success      200 {object} models.User
//	 	@Failure   	  500 "Unable to retrieve user or unable to delete user"
//		@Failure   	  400 "Invalid request"
//		@Router       /user/:id [delete]
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
