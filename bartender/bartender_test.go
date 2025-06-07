package bartender

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/models"
	"github.com/stretchr/testify/assert"
)

// Mock Notifier
type MockNotifier struct {
	Sent bool
}

func (m *MockNotifier) SendMessage(title, message string) error {
	m.Sent = true
	return nil
}
func (m *MockNotifier) SendFile(fileLoc string) error { return nil }

func TestGetRandomDrinkFromApi(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	testCache := cache.NewCache[models.DrinkResponse](10)

	// Mock GetDrinkFunc
	mockGetDrink := func(liquor string, c *gin.Context) (models.GetRandomDrinkAPI, error) {
		var apiResp models.GetRandomDrinkAPI
		jsonStr := `{
			"drinks": [{
				"idDrink": "123",
				"strDrink": "Test Drink",
				"strCategory": "Test Category",
				"strGlass": "Test Glass",
				"strAlcoholic": "Alcoholic",
				"strInstructions": "Test Instructions"
			}]
		}`
		_ = json.Unmarshal([]byte(jsonStr), &apiResp)
		return apiResp, nil
	}

	// Mock GatherIngredientsFunc
	mockGatherIngredients := func(drink models.GetRandomDrinkAPI) []string {
		return []string{"Ingredient1", "Ingredient2"}
	}

	// Mock Notifier
	mockNotifier := &MockNotifier{}

	// Service with mocks
	service := &DrinkService{
		GetDrinkFunc:          mockGetDrink,
		GatherIngredientsFunc: mockGatherIngredients,
		Notifier:              mockNotifier,
	}

	// Act
	service.GetRandomDrinkFromApi("", c, testCache)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	cached, found := testCache.Get("Test Drink")
	assert.True(t, found, "Drink should be cached")
	assert.Equal(t, "Test Drink", cached.Name)
	testCache.Clear()
	assert.Equal(t, 0, len(testCache.GetAll()))

	expected := models.DrinkResponse{
		Message:      "Drink of the Day",
		ExternalId:   "123",
		Name:         "Test Drink",
		Category:     "Test Category",
		Glass:        "Test Glass",
		Ingredients:  "Ingredient1, Ingredient2",
		Instructions: "Test Instructions",
	}
	var actual models.DrinkResponse
	_ = json.Unmarshal(w.Body.Bytes(), &actual)
	assert.Equal(t, expected, actual)

	for k := range testCache.GetAll() {
		t.Logf("Cache key: %s", k)
	}
}
