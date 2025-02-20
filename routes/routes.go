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
		commentGroup.DELETE("/:id", controllers.DeleteComment)
	}

	likeGroup := router.Group("/likes")
	{
		likeGroup.POST("/", controllers.CreateLike)
		likeGroup.GET("/user/:user_id", controllers.GetLikesByUser)
		likeGroup.DELETE("/:id", controllers.DeleteLike)
	}
}