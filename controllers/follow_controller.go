package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create a like (prevents duplicate likes)
func CreateFollow(c *gin.Context) {
	var follow models.Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if already exists
	var existingFollow models.Follow
	err := config.DB.Where("user_id = ? AND owner_id = ?", follow.UserID, follow.OwnerID).First(&existingFollow).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already follow this user"})
		return
	}

	// Create new like
	config.DB.Create(&follow)

	config.DB.Model(&models.User{}).Where("user_id = ?", follow.OwnerID).Update("following", gorm.Expr("following + 1"))

	config.DB.Model(&models.User{}).Where("user_id = ?", follow.UserID).Update("followers", gorm.Expr("followers + 1"))

	c.JSON(http.StatusCreated, gin.H{"message": "User followed successfully"})
}

// Get all likes by a user
func GetFollowersByUser(c *gin.Context) {
	userID := c.Param("user_id")

	var followers []models.Follow
	config.DB.Where("user_id = ?", userID).Preload("OwnerUser").Preload("User").Find(&followers)

	c.JSON(http.StatusOK, followers)
}

func GetFollowsByUser(c *gin.Context) {
	userID := c.Param("user_id")

	var following []models.Follow
	config.DB.Where("owner_id = ?", userID).Preload("User").Find(&following)

	c.JSON(http.StatusOK, following)
}

func GetUsersPaginated(c *gin.Context) {
	userID := c.Param("user_id") // Optional
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var followers []models.Follow
	var total int64

	query := config.DB.Model(&models.Follow{})

	// Apply filter if user_id is provided
	if userID != "" {
		query = query.Where("owner_id = ?", userID)
	}

	// Get total count
	query.Count(&total)

	// Fetch paginated followers
	query.Preload("User").
		Limit(limit).
		Offset(offset).
		Find(&followers)

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"limit":     limit,
		"total":     total,
		"followers": followers,
	})
}

// Get paginated list of followers (works with or without user_id)
func GetFollowersPaginated(c *gin.Context) {
	userID := c.Param("user_id") // Optional
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var followers []models.Follow
	var total int64

	query := config.DB.Model(&models.Follow{})

	// Apply filter if user_id is provided
	if userID != "" {
		query = query.Where("owner_id = ?", userID)
	}

	// Get total count
	query.Count(&total)

	// Fetch paginated followers
	query.Preload("User").
		Limit(limit).
		Offset(offset).
		Find(&followers)

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"limit":     limit,
		"total":     total,
		"followers": followers,
	})
}


// Delete a like and update like count
func DeleteFollow(c *gin.Context) {
	var follow models.Follow
	if err := c.ShouldBindJSON(follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if already exists
	var existingFollow models.Follow
	err := config.DB.Where("user_id = ? AND owner_id = ?", follow.UserID, follow.OwnerID).First(&existingFollow).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already follow this user"})
		return
	}

	// Delete 
	if err := config.DB.Delete(&existingFollow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	// Decrement count safely
	if err := config.DB.Model(&models.User{}).
		Where("user_id = ? AND followers > 0", follow.UserID).
		Update("followers", gorm.Expr("followers - 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update followers"})
		return
	}

	if err := config.DB.Model(&models.User{}).
	Where("user_id = ? AND following > 0", follow.OwnerID).
	Update("following", gorm.Expr("following - 1")).Error; err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update followers"})
	return
}

	c.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully"})
}
