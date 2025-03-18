package controllers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"social-backend/auth"
	"social-backend/config"
	"social-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Get all users
func GetUsers(c *gin.Context) {
	var users []models.User
	result := config.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Create new user
func CreateUser(c *gin.Context) {
	var input models.NewUser
	// Log raw request body to debug 400 errors
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Restore the request body (since ReadAll consumes it)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Log raw JSON request for debugging
	log.Printf("Received request body: %s", string(body))

	// Log the incoming JSON request
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Invalid JSON input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	log.Printf("Received user creation request: %+v\n", input)

	// Validate required fields
	if input.Username == "" || input.Email == "" || input.Password == "" {
		log.Println("Validation failed: Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}

	// Hash password before saving
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	log.Println("Password hashing failed:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
	// 	return
	// }

	// Default online status if not provided
	if input.OnlineStatus == "" {
		input.OnlineStatus = "offline"
	}

	// Create user struct with hashed password
	user := models.User{
		Username:     input.Username,
		ProfileName:  input.ProfileName,
		Email:        input.Email,
		Password:     input.Password,
		Bio:          input.Bio,
		PhoneNumber:  input.PhoneNumber,
		ProfilePic:   input.ProfilePic,
		OnlineStatus: input.OnlineStatus,
	}

	// Save user in database
	if err := config.DB.Create(&user).Error; err != nil {
		log.Println("Database error:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User created successfully: ID=%d, Username=%s\n", user.UserID, user.Username)

	// Successful response
	c.JSON(http.StatusCreated, gin.H{
		"message":      "User created successfully",
		"user_id":      user.UserID,
		"date_created": user.DateCreated,
	})
}

// Update user info
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&user).Where("user_id = ?", user.UserID).Updates(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
func GetUserById(c *gin.Context) {
	ownerID, exists := c.Get("user_id") // Extract logged-in user ID from context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user models.User
	userID := c.Param("user_id")

	// Retrieve details
	result := config.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var existingFollow models.Follow
	err := config.DB.Where("user_id = ? AND owner_id = ?", userID, ownerID).First(&existingFollow).Error
	followExists := err == nil

	c.JSON(http.StatusOK, gin.H{
		"user":        user,
		"isFollowing": followExists,
	})
}

// Delete user
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := config.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// User login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func RegisterSession(c *gin.Context) {
	var loginData struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func ValidateUserJWT(c *gin.Context) {
	userID, exists := c.Get("user_id") // Extract logged-in user ID from context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Println("User found")

	if userID != nil {
		c.JSON(http.StatusOK, gin.H{"status": "Authorized"})
		return
	}
}
