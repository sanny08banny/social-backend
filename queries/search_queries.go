package queries

import (
	"social-backend/config"
	"social-backend/models"
	"strings"
)

func SearchDatabase(query string, limit int, offset int) (map[string]interface{}, error) {
	searchQuery := "%" + strings.ToLower(query) + "%" // Ensure case-insensitive search

	var posts []models.Post
	var users []models.User
	var comments []models.Comment

	// Search in Posts (All Text Fields)
	if err := config.DB.Where(
		"LOWER(content) LIKE ? OR LOWER(category) LIKE ?",
		searchQuery, searchQuery).
		Limit(limit).Offset(offset).Preload("User").Find(&posts).Error; err != nil {
		return nil, err
	}

	// Search in Users (All Text Fields)
	if err := config.DB.Where(
		"LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(bio) LIKE ?",
		searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	// Search in Comments (All Text Fields)
	if err := config.DB.Where(
		"LOWER(content) LIKE ?",
		searchQuery).
		Limit(limit).Offset(offset).Find(&comments).Error; err != nil {
		return nil, err
	}

	// Structure response
	results := map[string]interface{}{
		"posts":    posts,
		"users":    users,
		"comments": comments,
	}

	return results, nil
}
