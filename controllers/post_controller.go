package controllers

import (
	"log"
	"net/http"
	"social-backend/config"
	"social-backend/models"
	"social-backend/utils"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get all posts with preloaded user data
func GetPosts(c *gin.Context) {
	userID, exists := c.Get("user_id") // Extract logged-in user ID from context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var posts []models.Post
	result := config.DB.Preload("User").Order("date_created DESC").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Get post IDs for batch querying
	postIDs := make([]uint, len(posts))
	for i, post := range posts {
		postIDs[i] = post.PostID
	}

	// Fetch likes and bookmarks in bulk for efficiency
	likedPostIDs := getUserLikedPostIDs(userID.(uint), postIDs)
	bookmarkedPostIDs := getUserBookmarkedPostIDs(userID.(uint), postIDs)

	// Assign `IsLiked` and `IsBookmarked`
	for i := range posts {
		posts[i].IsLiked = likedPostIDs[posts[i].PostID]
		posts[i].IsBookmarked = bookmarkedPostIDs[posts[i].PostID]
	}

	c.JSON(http.StatusOK, posts)
}
func getUserLikedPostIDs(userID uint, postIDs []uint) map[uint]bool {
	var likedPosts []uint
	config.DB.Table("likes").
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Pluck("post_id", &likedPosts)

	likedMap := make(map[uint]bool)
	for _, id := range likedPosts {
		likedMap[id] = true
	}
	return likedMap
}

func getUserBookmarkedPostIDs(userID uint, postIDs []uint) map[uint]bool {
	var bookmarkedPosts []uint
	config.DB.Table("bookmarks").
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Pluck("post_id", &bookmarkedPosts)

	bookmarkedMap := make(map[uint]bool)
	for _, id := range bookmarkedPosts {
		bookmarkedMap[id] = true
	}
	return bookmarkedMap
}

// Get a post by ID and track unique views
func GetPostById(c *gin.Context) {
	var post models.Post
	postID := c.Param("post_id")

	// Retrieve user ID if logged in (from middleware)
	userID, exists := c.Get("user_id")
	var userIDInt uint
	if exists && userID != nil {
		userIDInt = userID.(uint)
	}

	// Get IP address for guest users
	ipAddress := c.ClientIP()

	// Check if this view is unique
	var existingView models.PostView
	err := config.DB.Where("post_id = ? AND (user_id = ? OR ip_address = ?)", postID, userIDInt, ipAddress).First(&existingView).Error

	// If no existing view, add a new entry and increment view count
	if err != nil {
		newView := models.PostView{
			PostID:    strToUint(postID),
			UserID:    userIDInt,
			IPAddress: ipAddress,
		}
		config.DB.Create(&newView)

		// Increment view count in database
		config.DB.Model(&models.Post{}).Where("post_id = ?", postID).Update("view_count", gorm.Expr("view_count + 1"))
	}

	// Retrieve post details
	result := config.DB.Preload("User").First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}
func GetPostsByUserId(c *gin.Context) {
	receivedUserID := c.Param("user_id")
	userID, exists := c.Get("user_id") // Extract logged-in user ID from context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var posts []models.Post
	result := config.DB.Where("user_id = ?", receivedUserID).Preload("User").Order("date_created DESC").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Get post IDs for batch querying
	postIDs := make([]uint, len(posts))
	for i, post := range posts {
		postIDs[i] = post.PostID
	}

	// Fetch likes and bookmarks in bulk for efficiency
	likedPostIDs := getUserLikedPostIDs(userID.(uint), postIDs)
	bookmarkedPostIDs := getUserBookmarkedPostIDs(userID.(uint), postIDs)

	// Assign `IsLiked` and `IsBookmarked`
	for i := range posts {
		posts[i].IsLiked = likedPostIDs[posts[i].PostID]
		posts[i].IsBookmarked = bookmarkedPostIDs[posts[i].PostID]
	}

	c.JSON(http.StatusOK, posts)
}
// Create a new post
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save post to DB first
	result := config.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Respond to client immediately
	c.JSON(http.StatusCreated, post)

	// Perform classification in a Goroutine (asynchronously)
	go classifyAndUpdatePost(post.PostID, post.Content)
}

// classifyAndUpdatePost classifies a post and updates it in the database
func classifyAndUpdatePost(postID uint, content string) {
	category := utils.ClassifyPost(content) // Get category

	// Update the post asynchronously
	config.DB.Model(&models.Post{}).Where("post_id = ?", postID).Update("category", category)
}

// Update a post
func UpdatePost(c *gin.Context) {
	postID := c.Param("post_id")
	var post models.Post

	// Check if the post exists
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&post).Updates(models.Post{
		Content: post.Content,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// Delete a post
func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	result := config.DB.Delete(&models.Post{}, postID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// Get paginated posts
func GetPaginatedPosts(c *gin.Context) {
	userID, exists := c.Get("user_id") // Extract logged-in user ID from context

	page := strToInt(c.DefaultQuery("page", "1"))
	pageSize := strToInt(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	var posts []models.Post
	var totalPosts int64

	// Get total post count for pagination
	config.DB.Model(&models.Post{}).Count(&totalPosts)

	// Fetch paginated posts
	result := config.DB.Preload("User").Order("date_created DESC").Limit(pageSize).Offset(offset).Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// If user is logged in, fetch additional data (likes & bookmarks)
	if exists && userID != nil {
		log.Println("Authenticated user:", userID)

		// Get post IDs for batch querying
		postIDs := make([]uint, len(posts))
		for i, post := range posts {
			postIDs[i] = post.PostID
		}

		// Fetch likes and bookmarks in bulk for efficiency
		likedPostIDs := getUserLikedPostIDs(userID.(uint), postIDs)
		bookmarkedPostIDs := getUserBookmarkedPostIDs(userID.(uint), postIDs)

		// Assign `IsLiked` and `IsBookmarked`
		for i := range posts {
			posts[i].IsLiked = likedPostIDs[posts[i].PostID]
			posts[i].IsBookmarked = bookmarkedPostIDs[posts[i].PostID]
		}
	}

	// Return paginated response
	c.JSON(http.StatusOK, gin.H{
		"data":        posts,
		"total_posts": totalPosts,
		"page":        page,
		"page_size":   pageSize,
	})
}


// Handle reposting
func Repost(c *gin.Context) {
	var repost models.Repost
	if err := c.ShouldBindJSON(&repost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user has already reposted this post
	var existingRepost models.Repost
	if err := config.DB.Where("original_post_id = ? AND user_id = ?", repost.OriginalPostID, repost.UserID).First(&existingRepost).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reposted this post"})
		return
	}

	// Create repost
	result := config.DB.Create(&repost)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Increment repost count
	config.DB.Model(&models.Post{}).Where("post_id = ?", repost.OriginalPostID).Update("repost_count", gorm.Expr("repost_count + 1"))

	c.JSON(http.StatusCreated, gin.H{"message": "Post reposted successfully"})
}

// Convert string to uint
func strToUint(s string) uint {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return uint(val)
}

// Convert string to int
func strToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 1
	}
	return val
}

func GetFriendsPosts(c *gin.Context) {
	var allPosts []models.Post
	userID, exists := c.Get("user_id") // Extract logged-in user ID from context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	log.Printf("Received request user: %s", userID)
	var follows []models.Follow

	if err := config.DB.Where("user_id = ?", userID).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uniqueIDs := uniqueSortedList(follows)

	for i := range uniqueIDs{

	var posts []models.Post
	result := config.DB.Where("user_id = ?", uniqueIDs[i]).Preload("User").Order("date_created DESC").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	allPosts = append(allPosts, posts...)
}

	// Get post IDs for batch querying
	postIDs := make([]uint, len(allPosts))
	for i, post := range allPosts {
		postIDs[i] = post.PostID
	}

	// Fetch likes and bookmarks in bulk for efficiency
	likedPostIDs := getUserLikedPostIDs(userID.(uint), postIDs)
	bookmarkedPostIDs := getUserBookmarkedPostIDs(userID.(uint), postIDs)

	// Assign `IsLiked` and `IsBookmarked`
	for i := range allPosts {
		allPosts[i].IsLiked = likedPostIDs[allPosts[i].PostID]
		allPosts[i].IsBookmarked = bookmarkedPostIDs[allPosts[i].PostID]
	}

	c.JSON(http.StatusOK, allPosts)
}
func uniqueSortedList(follows []models.Follow) []uint {
	// Step 1: Collect all user_id and owner_id in a slice
	var allIDs []uint
	for _, follow := range follows {
		allIDs = append(allIDs, follow.UserID, follow.OwnerID)
	}

	// Step 2: Sort the slice
	sort.Slice(allIDs, func(i, j int) bool { return allIDs[i] < allIDs[j] })

	// Step 3: Remove duplicates in O(n) time
	var uniqueIDs []uint
	for i, id := range allIDs {
		if i == 0 || id != allIDs[i-1] { // Add only if not duplicate
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	return uniqueIDs
}
func GetPostsAndReposts(c *gin.Context) {
    var posts []models.Post
    var reposts []models.Repost
    
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    offset := (page - 1) * pageSize
    
	var totalCount int64

	config.DB.Raw("SELECT (SELECT COUNT(*) FROM posts) + (SELECT COUNT(*) FROM reposts)").Scan(&totalCount)	
    
    var wg sync.WaitGroup
    var postErr, repostErr error
    
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        postErr = config.DB.Preload("User").Order("date_created DESC").Limit(pageSize).Offset(offset).Find(&posts).Error
    }()
    
    go func() {
        defer wg.Done()
        repostErr = config.DB.Preload("User").Preload("Post").Preload("Post.User").Order("reposted_at DESC").Limit(pageSize).Offset(offset).Find(&reposts).Error
    }()
    
    wg.Wait()
    
    if postErr != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
        return
    }
    
    if repostErr != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reposts"})
        return
    }
    
	combined := make([]interface{}, 0, len(posts)+len(reposts))
    
    for _, post := range posts {
        combined = append(combined, post)
    }
    
    for _, repost := range reposts {
        combined = append(combined, repost)
    }
    
    // Sort combined list by date descending
    sort.Slice(combined, func(i, j int) bool {
        return getDate(combined[i]).After(getDate(combined[j]))
    })
    
	c.JSON(http.StatusOK, gin.H{
		"data":        combined,
		"total_posts": totalCount,
		"page":        page,
		"page_size":   pageSize,
	})
}

func getDate(item interface{}) time.Time {
    switch v := item.(type) {
    case models.Post:
        return v.DateCreated
    case models.Repost:
        return v.RepostedAt
    }
    return time.Time{}
}
