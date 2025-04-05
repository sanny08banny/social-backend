package routes

import (
	"social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	commentGroup := router.Group("/comments")
	{
		commentGroup.POST("/", controllers.CreateComment)
		commentGroup.GET("/user/:user_id", controllers.GetCommentsByUser)
		commentGroup.GET("/posts/:post_id", controllers.GetCommentsByPost)
		commentGroup.DELETE("/:id", controllers.DeleteComment)
	}

	likeGroup := router.Group("/likes")
	{
		likeGroup.POST("/", controllers.CreateLike)
		likeGroup.GET("/user/:user_id", controllers.GetLikesByUser)
		likeGroup.DELETE("/", controllers.DeleteLike)
	}
	bookmarkGroup := router.Group("/bookmarks")
	{
		bookmarkGroup.POST("", controllers.CreateBookmark)
		bookmarkGroup.GET("/user/:user_id", controllers.GetBookmarksByUser)
		bookmarkGroup.DELETE("/", controllers.DeleteBookmark)
	}
	repostGroup := router.Group("/reposts")
	{
		repostGroup.POST("/", controllers.CreateRepost)
		repostGroup.GET("/user/:user_id", controllers.GetRepostsByUser)
		repostGroup.DELETE("/", controllers.DeleteRepost)
	}

	categoriesGroup := router.Group("")
	{
		categoriesGroup.GET("/categories", controllers.GetCategoriesWithPosts)
		categoriesGroup.GET("/categories/:category/posts", controllers.GetPostsByCategory)
	}
	router.GET("/search", controllers.SearchHandler)
}
