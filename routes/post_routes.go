package routes

import (
	"social-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	router.GET("/posts", controllers.GetPosts)
	router.POST("/posts", controllers.CreatePost)
	router.PUT("/posts", controllers.UpdatePost)
	router.DELETE("/post/:id", controllers.DeletePost)
}