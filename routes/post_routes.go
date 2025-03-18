package routes

import (
	"social-backend/controllers"
	"social-backend/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {

	protected := router.Group("")
	protected.Use(middleware.JWTAuthMiddleware()) // Apply JWT middleware
	{
		protected.GET("/posts/:post_id", controllers.GetPostById)
		protected.GET("/posts/user/:user_id", controllers.GetPostsByUserId)
		protected.GET("/posts", controllers.GetPosts)
		protected.GET("/posts/paginated", controllers.GetPaginatedPosts)
	}

	router.POST("/posts", controllers.CreatePost)
	router.PUT("/posts", controllers.UpdatePost)
	router.DELETE("/post/:id", controllers.DeletePost)
}