package help

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Help returns a list of endpoints
func Help(c *gin.Context) {
	endpoints := map[string]string{
		"GET /help":              "List of endpoints...",
		"GET /healthcheck":       "Healthcheck endpoint for monitoring tools",
		"GET /ebook/find/:title": "Check if a book exists in the Gutenberg project",
		// "/ebook/dl/:title": "Download a book from the Gutenberg project",
		"GET /dinner/random":           "Get three random dinner recipes",
		"GET /dinner/recipe/:id":       "Get a specific recipe based on id",
		"POST /dinner/cache/backup":    "Backup the dinner cache to a file",
		"GET /bartender/random":        "Get a random cocktail recipe",
		"POST /bartender/save":         "Save a cocktail recipe to the database",
		"GET /bartender/history":       "Get the history of cocktails received",
		"POST /bartender/cache/backup": "Backup the cocktail cache to a file",
		"POST /database/backup":        "Backup the database to a file",
	}

	c.JSON(http.StatusOK, gin.H{"body": endpoints})
}
