package dinner

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	spoonacular "github.com/ddsky/spoonacular-api-clients/go"
	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/models"
	"github.com/stretchr/testify/assert"
)

func setupSpoonacularFakeServer(t *testing.T, recipesJSON string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(recipesJSON))
	}))
}

func TestGetRandomRecipes(t *testing.T) {
	// Fake Spoonacular API response
	data, err := os.ReadFile("testdata/recipes.json")
	if err != nil {
		t.Fatalf("Failed to read recipes.json: %v", err)
	}

	recipesJSON := string(data)

	server := setupSpoonacularFakeServer(t, recipesJSON)
	defer server.Close()

	// Set up the real Spoonacular client to use the fake server
	configuration := spoonacular.NewConfiguration()
	configuration.Servers = []spoonacular.ServerConfiguration{
		{URL: server.URL},
	}
	configuration.AddDefaultHeader("x-api-key", "fake-key")
	apiClient = spoonacular.NewAPIClient(configuration).RecipesAPI

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	GetRandomRecipes(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Pasta with Garlic, Scallions, Cauliflower & Breadcrumbs")
	// assert.Contains(t, w.Body.String(), "Test Recipe 2")
	// assert.Contains(t, w.Body.String(), "Test Recipe 3")
}

func TestGetRecipeFromApi_CacheHit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)
	testRecipe := models.RecipeInfo{Title: "Cached Recipe", Id: 123}
	testCache.Set("123", testRecipe, 0)

	GetRecipeFromApi(c, "123", testCache)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Cached Recipe")
}

func TestGetRecipeFromApi_CacheMiss(t *testing.T) {
	// Fake Spoonacular API response for a single recipe
	recipeJSON := `{
		"id": 999,
		"title": "API Recipe",
		"sourceUrl": "http://example.com",
		"instructions": "Step 1. Step 2.",
		"extendedIngredients": [
			{"amount": 1, "unit": "cup", "name": "Sugar"}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(recipeJSON))
	}))
	defer server.Close()

	configuration := spoonacular.NewConfiguration()
	configuration.Servers = []spoonacular.ServerConfiguration{
		{URL: server.URL},
	}
	configuration.AddDefaultHeader("x-api-key", "fake-key")
	apiClient = spoonacular.NewAPIClient(configuration).RecipesAPI

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)

	GetRecipeFromApi(c, "999", testCache)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "API Recipe")
	assert.Contains(t, w.Body.String(), "Sugar")
}

func TestGetRecipeFromApi_ApiError(t *testing.T) {
	// Simulate Spoonacular API error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "API error", http.StatusInternalServerError)
	}))
	defer server.Close()

	configuration := spoonacular.NewConfiguration()
	configuration.Servers = []spoonacular.ServerConfiguration{
		{URL: server.URL},
	}
	configuration.AddDefaultHeader("x-api-key", "fake-key")
	apiClient = spoonacular.NewAPIClient(configuration).RecipesAPI

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)

	GetRecipeFromApi(c, "999", testCache)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching recipe")
}

// You can add similar tests for GetRecipeFromDB and SaveRecipe, using a test database or a mock DB if needed.
