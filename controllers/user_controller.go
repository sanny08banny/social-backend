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
		rows.Scan(&user.UserID, &user.Username, &user.ProfileName, &user.Email, &user.Bio, &user.PhoneNumber, &user.ProfilePic, &user.OnlineStatus, &user.DateCreated, &user.LastUpdated)
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
		user.Username, user.ProfileName, user.Email, user.Bio, user.PhoneNumber, user.ProfilePic, user.OnlineStatus,
	).Scan(&user.UserID, &user.DateCreated, &user.LastUpdated)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(context.Background(), queries.UpdateUserQuery, user.Username, user.ProfileName, user.Email, user.Bio, user.PhoneNumber, user.ProfilePic, user.OnlineStatus, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	_, err := config.DB.Exec(context.Background(), queries.DeleteUserQuery, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}