package dinner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	spoonacular "github.com/ddsky/spoonacular-api-clients/go"
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

type MockSpoonacularAdapter struct {
	RecipeJSON string
}

func (m *MockSpoonacularAdapter) GetRecipeInformation(ctx context.Context, id int32) (*spoonacularapi.RecipeInformationOverride, error) {
	var recipe spoonacularapi.RecipeInformationOverride
	err := json.Unmarshal([]byte(m.RecipeJSON), &recipe)
	return &recipe, err
}

func (m *MockSpoonacularAdapter) GetRandomRecipes(ctx context.Context, number int) (*spoonacularapi.RandomRecipesResponse, error) {
	return &spoonacularapi.RandomRecipesResponse{}, nil
}

type ErrorMockSpoonacularAdapter struct{}

func (m *ErrorMockSpoonacularAdapter) GetRecipeInformation(ctx context.Context, id int32) (*spoonacularapi.RecipeInformationOverride, error) {
	return nil, fmt.Errorf("API error")
}

func (m *ErrorMockSpoonacularAdapter) GetRandomRecipes(ctx context.Context, number int) (*spoonacularapi.RandomRecipesResponse, error) {
	return &spoonacularapi.RandomRecipesResponse{}, nil
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

	adapter := &spoonacularapi.SpoonacularAdapter{
		RealClient: mockClient.Client,
	}

	// Set up Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	GetRandomRecipes(c, adapter)

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
	// Create the mock client with our test data
	mockClient := NewMockSpoonacularClient(t, "")
	mockClient.Client.SetBaseURL(mockClient.MockServer.URL)
	defer mockClient.MockServer.Close()

	adapter := &spoonacularapi.SpoonacularAdapter{
		RealClient: mockClient.Client,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)
	testRecipe := models.RecipeInfo{Title: "Brownies", Id: 123, Url: "https://www.test.com", Ingredients: "1 cup of sugar, 1 cup of flour, 1 cup of chocolate chips", Instructions: "Bake in oven at 350 degrees for 20 minutes"}
	ttl := 5 * time.Minute
	testCache.Set("123", testRecipe, ttl)

	GetRecipeFromApi(c, "123", testCache, adapter)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Brownies")
	assert.Contains(t, w.Body.String(), "https://www.test.com")
	assert.Contains(t, w.Body.String(), "1 cup of sugar, 1 cup of flour, 1 cup of chocolate chips")
	assert.Contains(t, w.Body.String(), "Bake in oven at 350 degrees for 20 minutes")
}

func TestGetRecipeFromApi_CacheMiss(t *testing.T) {
	// Mock Spoonacular API response (matches expected structure)
	data, err := os.ReadFile("testdata/recipe.json")
	if err != nil {
		t.Fatalf("Failed to read recipe.json: %v", err)
	}
	adapter := &MockSpoonacularAdapter{RecipeJSON: string(data)}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	testCache := cache.NewCache[models.RecipeInfo](10)

	// Act
	GetRecipeFromApi(c, "716429", testCache, adapter)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var resp models.RecipeInfo
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, int32(716429), resp.Id)
	assert.Equal(t, "Pasta with Garlic, Scallions, Cauliflower & Breadcrumbs", resp.Title)
	assert.Contains(t, resp.Ingredients, "1tbsp of butter")
	assert.Contains(t, resp.Instructions, "Preheat oven to 400 degrees")
	assert.Contains(t, resp.Ingredients, "butter")
	assert.Contains(t, resp.Ingredients, "scallions")
}

func TestGetRecipeFromApi_ApiError(t *testing.T) {
	recipesJSON := `{
		"recipes": [
			{ "id": 1, "title": "Recipe 1" },
			{ "id": 2, "title": "Recipe 2" },
			{ "id": 3, "title": "Recipe 3" }
		]
	}`

	// Simulate Spoonacular API error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "API error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Set up the API client to use the test server
	configuration := spoonacular.NewConfiguration()
	configuration.Servers = []spoonacular.ServerConfiguration{
		{URL: server.URL},
	}
	configuration.AddDefaultHeader("x-api-key", "fake-key")

	mockClient := NewMockSpoonacularClient(t, recipesJSON)
	mockClient.Client.SetBaseURL(mockClient.MockServer.URL)
	defer mockClient.MockServer.Close()

	adapter := &ErrorMockSpoonacularAdapter{}

	// Set up Gin context and cache
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testCache := cache.NewCache[models.RecipeInfo](10)

	// Call the handler
	GetRecipeFromApi(c, "999", testCache, adapter)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error fetching recipe")
}
