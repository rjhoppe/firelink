package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthcheck returns a JSON response with the status of the application
func Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Firelink is running"})
}
