package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"

	"github.com/gin-gonic/gin"
)

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

    // Optional: Check if Post exists if post_id is provided
    if comment.PostID != nil {
        var post models.Post
        if err := config.DB.First(&post, *comment.PostID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
            return
        }
    }

    // Optional: Check if Parent Comment exists if parent_id is provided
    if comment.ParentID != nil {
        var parentComment models.Comment
        if err := config.DB.First(&parentComment, *comment.ParentID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Parent Comment not found"})
            return
        }
    }

    if err := config.DB.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, comment)
}

func GetCommentsByUser(c *gin.Context) {
	var comments []models.Comment
	userID := c.Param("user_id")
	config.DB.Preload("Replies").Where("user_id = ?", userID).Find(&comments)
	c.JSON(http.StatusOK, comments)
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	config.DB.Delete(&models.Comment{}, commentID)
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
