package controllers

import (
	"context"
	"net/http"
	"social-backend/config"
	"social-backend/models"
	"social-backend/queries"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), queries.GetUsersQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		rows.Scan(&user.ID, &user.Name, &user.Email)
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := config.DB.QueryRow(
		context.Background(),
		queries.CreateUserQuery,
		user.Name, user.Email,
	).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

