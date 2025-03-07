package controllers

import (
	"net/http"
	"social-backend/config"
	"social-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	var posts []models.Post
	result := config.DB.Preload("User").Preload("Comments").Preload("Likes").Preload("BookMarks").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetPostById(c *gin.Context) {
	var post []models.Post
	postID := c.Param("post_id")
	config.DB.Where("post_id = ?", postID).Preload("User").Preload("Comments").Preload("Likes").Preload("BookMarks").Find(&post)
	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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

	result := config.DB.Model(&post).Where("post_id = ?", post.PostID).Updates(models.Post{
		Content: post.Content,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	result := config.DB.Delete(&models.Post{}, postID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
func GetPaginatedPosts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	var posts []models.Post
	var totalPosts int64

	config.DB.Model(&models.Post{}).Count(&totalPosts)

	// Convert page and pageSize to integers
	pageInt := strToInt(page)
	pageSizeInt := strToInt(pageSize)

	offset := (pageInt - 1) * pageSizeInt

	result := config.DB.Preload("User").Limit(pageSizeInt).Offset(offset).Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var postDTOs []models.PostDTO
	for _, post := range posts {
		var commentCount int64
		var likeCount int64
		var bookmarkCount int64

		config.DB.Model(&models.Comment{}).Where("post_id = ?", post.PostID).Count(&commentCount)
		config.DB.Model(&models.Like{}).Where("post_id = ?", post.PostID).Count(&likeCount)
		config.DB.Model(&models.BookMark{}).Where("post_id = ?", post.PostID).Count(&bookmarkCount)


		postDTOs = append(postDTOs, models.PostDTO{
			PostID:      post.PostID,
			Content:     post.Content,
			User:        models.UserDTO{UserID: post.User.UserID, Username: post.User.Username, ProfileName: post.User.ProfileName, ProfilePic: post.User.ProfilePic},
			DateCreated: post.DateCreated.Format("2006-01-02 15:04:05"),
			LastUpdated: post.LastUpdated.Format("2006-01-02 15:04:05"),
			CommentCount: commentCount,
			LikeCount:    likeCount,
			BookMarkCount: bookmarkCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       postDTOs,
		"total_posts": totalPosts,
		"page":        pageInt,
		"page_size":   pageSizeInt,
	})
}

func strToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 1 
	}
	return val
}
