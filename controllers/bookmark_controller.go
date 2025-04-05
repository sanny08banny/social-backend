package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"
	"social-backend/queries"

	"github.com/gin-gonic/gin"
)

// Create a bookmark and update count
func CreateBookmark(c *gin.Context) {
	var bookmark models.BookMark
	if err := c.ShouldBindJSON(&bookmark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent duplicate bookmarks
	var existingBookmark models.BookMark
	err := config.DB.Where("user_id = ? AND post_id = ?", bookmark.UserID, bookmark.PostID).First(&existingBookmark).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post already bookmarked"})
		return
	}


	// Insert bookmark and update count
	config.DB.Exec(queries.CreateBookmarkQuery, bookmark.UserID, bookmark.PostID)
	config.DB.Exec(queries.UpdateBookmarkCountQuery, bookmark.PostID)

	c.JSON(http.StatusCreated, gin.H{"message": "Bookmark added successfully"})
}

// Get all bookmarks by a user
func GetBookmarksByUser(c *gin.Context) {
	var bookmarks []models.BookMark
	userID := c.Param("user_id")
	result := config.DB.Where("user_id = ?", userID).Preload("Post").Preload("Post.User").Find(&bookmarks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, bookmarks)
}

// Delete bookmark and update count
func DeleteBookmark(c *gin.Context) {
	PostID := c.Query("post_id")
	UserID := c.Query("user_id")

	var existingBookmark models.BookMark
	err := config.DB.Where("user_id = ? AND post_id = ?", UserID, PostID).First(&existingBookmark).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bookmark doesn't exist"})
		return
	}

	// Delete bookmark and update count
	config.DB.Exec(queries.DeleteBookmarkQuery, existingBookmark.BookMarkID)
	config.DB.Exec(queries.UpdateBookmarkCountQuery, PostID)

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed successfully"})
}
