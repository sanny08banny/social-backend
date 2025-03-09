package routes

import (
	"social-backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", controllers.GetUsers)
	router.POST("/users", controllers.CreateUser)
	router.PUT("/users",controllers.UpdateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)
	router.POST("/login", controllers.LoginUser)
	router.POST("users/login",controllers.Login)
}