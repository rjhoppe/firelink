package dinner

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/models"
	"github.com/rjhoppe/firelink/spoonacularapi"
	"github.com/stretchr/testify/assert"
)

// MockSpoonacularClient is a mockable implementation of your custom client
type MockSpoonacularClient struct {
	*spoonacularapi.Client
	MockServer *httptest.Server
}

// NewMockSpoonacularClient creates a new mock client for testing
func NewMockSpoonacularClient(t *testing.T, responseBody string) *MockSpoonacularClient {
	// Create a test server that returns the specified response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))

	// Create a custom HTTP client that points to our test server
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				// Rewrite the URL to point to our test server
				return url.Parse(server.URL)
			},
		},
	}

	// Create the client with a fake API key
	client := spoonacularapi.NewClient("fake-api-key")

	// Replace the HTTP client with our custom one
	client.SetHTTPClient(httpClient)

	return &MockSpoonacularClient{
		Client:     client,
		MockServer: server,
	}
}

func TestGetRandomRecipes(t *testing.T) {
	// Fake Spoonacular API response

	recipesJSON := `{
		"recipes": [
			{ "id": 1, "title": "Recipe 1" },
			{ "id": 2, "title": "Recipe 2" },
			{ "id": 3, "title": "Recipe 3" }
		]
	}`

	// Create the mock client with our test data
	mockClient := NewMockSpoonacularClient(t, recipesJSON)
	mockClient.Client.SetBaseURL(mockClient.MockServer.URL)
	defer mockClient.MockServer.Close()

	apiClient = mockClient.Client

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	GetRandomRecipes(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body and make assertions
	var response models.RandomRecipes
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert that the response contains three recipes
	assert.NotEmpty(t, response.RecipeOne)
	assert.NotEmpty(t, response.RecipeTwo)
	assert.NotEmpty(t, response.RecipeThree)
}

func TestGetRecipeFromApi_CacheHit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)
	testRecipe := models.RecipeInfo{Title: "Brownies", Id: 123, Url: "https://www.test.com", Ingredients: "1 cup of sugar, 1 cup of flour, 1 cup of chocolate chips", Instructions: "Bake in oven at 350 degrees for 20 minutes"}
	ttl := 5 * time.Minute
	testCache.Set("123", testRecipe, ttl)

	GetRecipeFromApi(c, "123", testCache)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Brownies")
	assert.Contains(t, w.Body.String(), "https://www.test.com")
	assert.Contains(t, w.Body.String(), "1 cup of sugar, 1 cup of flour, 1 cup of chocolate chips")
	assert.Contains(t, w.Body.String(), "Bake in oven at 350 degrees for 20 minutes")
}

func TestGetRecipeFromApi_CacheMiss(t *testing.T) {
	// Fake Spoonacular API response for a single recipe
	data, err := os.ReadFile("testdata/recipe.json")
	if err != nil {
		t.Fatalf("Failed to read recipes.json: %v", err)
	}

	recipeJSON := string(data)

	// Create the mock client with our test data
	mockClient := NewMockSpoonacularClient(t, recipeJSON)
	mockClient.Client.SetBaseURL(mockClient.MockServer.URL)
	defer mockClient.MockServer.Close()

	apiClient = mockClient.Client

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)

	GetRecipeFromApi(c, "999", testCache)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "API Recipe")
	assert.Contains(t, w.Body.String(), "Sugar")
}

// func TestGetRecipeFromApi_ApiError(t *testing.T) {
// 	// Simulate Spoonacular API error
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		http.Error(w, "API error", http.StatusInternalServerError)
// 	}))
// 	defer server.Close()

// 	configuration := spoonacular.NewConfiguration()
// 	configuration.Servers = []spoonacular.ServerConfiguration{
// 		{URL: server.URL},
// 	}
// 	configuration.AddDefaultHeader("x-api-key", "fake-key")
// 	apiClient = spoonacular.NewAPIClient(configuration).RecipesAPI

// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)

// 	testCache := cache.NewCache[models.RecipeInfo](10)

// 	GetRecipeFromApi(c, "999", testCache)

// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
// 	assert.Contains(t, w.Body.String(), "Error fetching recipe")
// }

// You can add similar tests for GetRecipeFromDB and SaveRecipe, using a test database or a mock DB if needed.
