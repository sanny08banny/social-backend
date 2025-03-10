package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment - Create a new comment or reply
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if User exists
	var user models.User
	if err := config.DB.First(&user, comment.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Check if Post exists
	if comment.PostID != nil {
		var post models.Post
		if err := config.DB.First(&post, *comment.PostID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
			return
		}
	}

	// Check if Parent Comment exists if it's a reply
	if comment.ParentID != nil {
		var parentComment models.Comment
		if err := config.DB.First(&parentComment, *comment.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Parent Comment not found",
				"parent_id":   *comment.ParentID,
				"suggestion":  "Ensure the parent comment exists before replying.",
			})
			return
		}
	}

	// Create the comment
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&models.Post{}).Where("post_id = ?", comment.PostID).Update("comment_count", gorm.Expr("comment_count + 1"))

	c.JSON(http.StatusCreated, comment)
}


// GetCommentsByPost - Retrieve all comments on a post, including nested replies
func GetCommentsByPost(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []models.Comment

	config.DB.
		Preload("User").
		Preload("Replies.User").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Find(&comments)

	c.JSON(http.StatusOK, comments)
}

// GetCommentsByUser - Retrieve all comments made by a user
func GetCommentsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	var comments []models.Comment

	config.DB.Preload("Replies").Where("user_id = ?", userID).Find(&comments)
	c.JSON(http.StatusOK, comments)
}

// DeleteComment - Deletes a comment
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	config.DB.Delete(&models.Comment{}, commentID)
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
