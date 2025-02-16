package controllers

import (
	"context"
	"net/http"
	"social-backend/config"
	"social-backend/models"
	"social-backend/queries"

	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), queries.GetPostsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.DateCreated, &post.LastUpdated)
		posts = append(posts, post)
	}
	c.JSON(http.StatusOK, posts)
}

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := config.DB.QueryRow(
		context.Background(),
		queries.CreatePostQuery,
		post.UserID, post.Content,
	).Scan(&post.PostID, &post.DateCreated, &post.LastUpdated)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(context.Background(), queries.UpdatePostQuery, post.Content, post.PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	_, err := config.DB.Exec(context.Background(), queries.DeletePostQuery, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}