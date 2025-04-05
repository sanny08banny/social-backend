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
	config.DB.Preload("Post").Preload("Post.User").Where("user_id = ?", userID).Find(&likes)
	c.JSON(http.StatusOK, likes)
}

// Delete a like and update like count
func DeleteLike(c *gin.Context) {
	UserID := c.Query("user_id")
	PostID := c.Query("post_id")

	// Check if the like exists
	var existingLike models.Like
	err := config.DB.Where("user_id = ? AND post_id = ?", UserID, PostID).First(&existingLike).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	// Delete like
	if err := config.DB.Delete(&existingLike).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
		return
	}

	// Decrement like count safely
	if err := config.DB.Model(&models.Post{}).
		Where("post_id = ? AND like_count > 0", PostID).
		Update("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}

