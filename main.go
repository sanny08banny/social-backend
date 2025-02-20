package main

import (
	"social-backend/config"
	"social-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()


	r := gin.Default()
	routes.UserRoutes(r)
	routes.PostRoutes(r)
	routes.ImageRoutes(r)
	routes.RegisterRoutes(r)
	
	// routes.PostRoutes(r)


	r.Run("0.0.0.0:8080")
}
