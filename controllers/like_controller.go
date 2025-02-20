package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateLike(c *gin.Context) {
	var like models.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&like)
	c.JSON(http.StatusCreated, like)
}

func GetLikesByUser(c *gin.Context) {
	var likes []models.Like
	userID := c.Param("user_id")
	config.DB.Where("user_id = ?", userID).Find(&likes)
	c.JSON(http.StatusOK, likes)
}

func DeleteLike(c *gin.Context) {
	likeID := c.Param("id")
	config.DB.Delete(&models.Like{}, likeID)
	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
}
