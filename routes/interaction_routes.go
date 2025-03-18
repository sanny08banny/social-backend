package routes

import (
	"social-backend/controllers"
	"social-backend/middleware"

	"github.com/gin-gonic/gin"
)

func InteractionRoutes(router *gin.Engine) {

	followGroup := router.Group("/follows").Use(middleware.JWTAuthMiddleware()) // Apply JWT middleware
	{
		
		followGroup.POST("/", controllers.CreateFollow)                // Follow a user
		followGroup.DELETE("/", controllers.DeleteFollow)  
		followGroup.GET("/", controllers.GetFollowsByUser)  
		followGroup.GET("/:user_id", controllers.GetFollowsByUser)          // Get followers of a user
	}
	followerGroup := router.Group("/followers").Use(middleware.JWTAuthMiddleware())
	{
		followerGroup.GET("/:user_id", controllers.GetFollowersByUser)          // Get followers of a user
		followerGroup.GET("/", controllers.GetFollowersPaginated)               // Get paginated followers (optional user_id)
		followerGroup.GET("/:user_id/paginated", controllers.GetFollowersPaginated) // Get paginated followers for a user
	}
}