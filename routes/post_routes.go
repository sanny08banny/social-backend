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
		protected.GET("/posts/paginated", controllers.GetPaginatedPosts)
	}


	router.GET("/posts", controllers.GetPosts)
	router.POST("/posts", controllers.CreatePost)
	router.PUT("/posts", controllers.UpdatePost)
	router.DELETE("/post/:id", controllers.DeletePost)
	router.GET("/posts/:post_id", controllers.GetPostById)
}