package controllers

import (
	"net/http"
	"social-backend/queries"
	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	query := c.Query("q") // Get search term from query parameter
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	page := strToInt(c.DefaultQuery("page", "1"))
	pageSize := strToInt(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	results, err := queries.SearchDatabase(query, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":       query,
		"results":     results,
		"page":        page,
		"page_size":   pageSize,
	})
}
