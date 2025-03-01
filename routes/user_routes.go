package routes

import (
	"social-backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", controllers.GetUsers)
	router.POST("/users", controllers.CreateUser)
	router.PUT("/users",controllers.UpdateUser)
}