package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create a like (prevents duplicate likes)
func CreateRepost(c *gin.Context) {
	var repost models.Repost
	if err := c.ShouldBindJSON(&repost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if already exists
	var existingRepost models.Repost
	err := config.DB.Where("user_id = ? AND original_post_id = ?", repost.UserID, repost.OriginalPostID).First(&existingRepost).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post already reposted"})
		return
	}

	// Create new like
	config.DB.Create(&repost)

	config.DB.Model(&models.Post{}).Where("post_id = ?", repost.OriginalPostID).Update("repost_count", gorm.Expr("repost_count + 1"))

	c.JSON(http.StatusCreated, gin.H{"message": "Thought reposted successfully"})
}

// Get all likes by a user
func GetRepostsByUser(c *gin.Context) {
	userID := c.Param("user_id")

	var reposts []models.Repost
	config.DB.Where("user_id = ?", userID).Find(&reposts)

	posts := make([]models.Post, len(reposts))


	for i, repost := range reposts {
		var post models.Post
		result := config.DB.Preload("User").First(&post, repost.OriginalPostID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		posts[i] = post 
	}
	c.JSON(http.StatusOK, posts)
}

// Delete a like and update like count
func DeleteRepost(c *gin.Context) {
	var repost models.Repost
	if err := c.ShouldBindJSON(&repost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if exists
	var existingRepost models.Repost
	err := config.DB.Where("user_id = ? AND original_post_id = ?", repost.UserID, repost.OriginalPostID).First(&existingRepost).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post already reposted"})
		return
	}

	// Delete 
	if err := config.DB.Delete(&existingRepost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove repost"})
		return
	}

	// Decrement count safely
	if err := config.DB.Model(&models.Post{}).
		Where("post_id = ? AND repost_count > 0", repost.OriginalPostID).
		Update("repost_count", gorm.Expr("repost_count - 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update repost count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repost removed successfully"})
}

