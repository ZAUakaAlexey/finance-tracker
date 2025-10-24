package handlers

import (
	"net/http"
	"strconv"

	"github.com/ZAUakaAlexey/backend_go/internal/database"
	"github.com/ZAUakaAlexey/backend_go/internal/models"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(context *gin.Context) {
	userId, exists := context.Get("user_id")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUser(context *gin.Context) {
	id := context.Param("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": user})
}

type UpdateUserInput struct {
	Name  string `json:"name" binding:"omitempty,fullname"`
	Email string `json:"email" binding:"omitempty,email"`
}

func UpdateUser(context *gin.Context) {
	id := context.Param("id")
	paramID, _ := strconv.ParseUint(id, 10, 32)

	loggedUserID, _ := context.Get("user_id")

	if uint(paramID) != loggedUserID.(uint) {
		context.JSON(http.StatusForbidden, gin.H{"error": "You have no persmissions"})
		return
	}

	var input UpdateUserInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := database.DB.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}
