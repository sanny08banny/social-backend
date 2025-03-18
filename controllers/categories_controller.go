package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"

	"github.com/gin-gonic/gin"
)

func GetCategoriesWithPosts(c *gin.Context) {
	var categoryCounts []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	// Query to count posts per category
	result := config.DB.Model(&models.Post{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Find(&categoryCounts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, categoryCounts)
}


func GetPostsByCategory(c *gin.Context) {
	category := c.Param("category") // Get category from URL

	var posts []models.Post

	// Fetch posts where category matches
	result := config.DB.Where("category = ?", category).Preload("User").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// If no posts found, return an empty list
	c.JSON(http.StatusOK, posts)
}

