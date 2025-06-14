package help

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelpEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/help", Help)

	req, _ := http.NewRequest("GET", "/help", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Check that some known endpoints are present
	body, ok := resp["body"]
	assert.True(t, ok)
	assert.Contains(t, body, "GET /help")
	assert.Contains(t, body, "GET /healthcheck")
	assert.Contains(t, body, "GET /dinner/random")
}
