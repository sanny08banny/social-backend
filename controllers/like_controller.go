package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create a like (prevents duplicate likes)
func CreateLike(c *gin.Context) {
	var like models.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the like already exists
	var existingLike models.Like
	err := config.DB.Where("user_id = ? AND post_id = ?", like.UserID, like.PostID).First(&existingLike).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post already liked"})
		return
	}

	// Create new like
	config.DB.Create(&like)

	// Increment like count
	config.DB.Model(&models.Post{}).Where("post_id = ?", like.PostID).Update("like_count", gorm.Expr("like_count + 1"))

	c.JSON(http.StatusCreated, gin.H{"message": "Like added successfully"})
}

// Get all likes by a user
func GetLikesByUser(c *gin.Context) {
	var likes []models.Like
	userID := c.Param("user_id")
	config.DB.Preload("Post").Where("user_id = ?", userID).Find(&likes)
	c.JSON(http.StatusOK, likes)
}

// Delete a like and update like count
func DeleteLike(c *gin.Context) {
	likeID := c.Param("id")

	// Retrieve like to get post_id before deleting
	var like models.Like
	if err := config.DB.First(&like, likeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	// Delete like
	config.DB.Delete(&like)

	// Decrement like count
	config.DB.Model(&models.Post{}).Where("post_id = ?", like.PostID).Update("like_count", gorm.Expr("like_count - 1"))

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}
