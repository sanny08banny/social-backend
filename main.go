package main

import (
	"context"
	"social-backend/config"
	"social-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	defer config.DB.Close(context.Background())

	r := gin.Default()
	routes.UserRoutes(r)
	// routes.PostRoutes(r)

	r.Run(":8080")
}
