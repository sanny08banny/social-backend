package utils

import (
	"strings"
)

// ClassifyPost determines the category of a post
func ClassifyPost(content string) string {
	contentLower := strings.ToLower(content)

	// Step 1: Score-based keyword matching (Primary method)
	categoryScores := make(map[string]int)

	for category, keywords := range CATEGORY_KEYWORDS {
		for _, keyword := range keywords {
			if strings.Contains(contentLower, strings.ToLower(keyword)) {
				categoryScores[category]++
			}
		}
	}

	// Step 2: Find the best-matching category
	var bestCategory string
	maxScore := 0

	for category, score := range categoryScores {
		if score > maxScore {
			maxScore = score
			bestCategory = category
		}
	}

	// Step 3: If no strong match, use word frequency classification
	if bestCategory == "" {
		mostRelevantCategory := findCategoryByWordFrequency(contentLower)
		if mostRelevantCategory != "" {
			return mostRelevantCategory
		}
	}

	// Step 4: Ensure no post is left uncategorized
	if bestCategory == "" {
		bestCategory = "miscellaneous"
	}

	return bestCategory
}


// findCategoryByWordFrequency attempts to classify based on most common words
func findCategoryByWordFrequency(content string) string {
	wordCount := make(map[string]int)

	// Tokenize words
	words := strings.Fields(content)
	for _, word := range words {
		word = strings.ToLower(strings.Trim(word, ".,!?\"'"))
		wordCount[word]++
	}

	// Check most frequent words against category keywords
	for category, keywords := range CATEGORY_KEYWORDS {
		for _, keyword := range keywords {
			if wordCount[keyword] > 0 { // If the keyword appears frequently, assign category
				return category
			}
		}
	}

	return "" // No category found based on frequency
}


