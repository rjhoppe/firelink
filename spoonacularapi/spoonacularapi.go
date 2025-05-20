package spoonacularapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	spoonacular "github.com/ddsky/spoonacular-api-clients/go"
)

// This package wraps the official Spoonacular client and adds custom implementations
// for endpoints that have JSON parsing issues

// Recipe represents a recipe from the Spoonacular API
type Recipe struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
	// Add other fields you need
}

// RandomRecipesResponse represents the response from the random recipes endpoint
type RandomRecipesResponse struct {
	Recipes []Recipe `json:"recipes"`
}

// Client wraps the official Spoonacular client and adds custom methods
type Client struct {
	apiKey     string
	apiClient  *spoonacular.APIClient
	httpClient *http.Client
	baseURL    string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for API requests
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient creates a new Spoonacular client wrapper
func NewClient(apiKey string, options ...ClientOption) *Client {
	// Create official client for endpoints that work correctly
	configuration := spoonacular.NewConfiguration()
	configuration.AddDefaultHeader("x-api-key", apiKey)
	apiClient := spoonacular.NewAPIClient(configuration)

	c := &Client{
		apiKey:     apiKey,
		apiClient:  apiClient,
		httpClient: &http.Client{},
		baseURL:    "https://api.spoonacular.com",
	}

	// Apply any custom options
	for _, option := range options {
		option(c)
	}

	return c
}

// SetHTTPClient allows replacing the HTTP client (useful for testing)
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetBaseURL allows replacing the base URL (useful for testing)
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// GetRandomRecipes gets random recipes from the Spoonacular API
// This is a custom implementation that doesn't rely on the official client
func (c *Client) GetRandomRecipes(ctx context.Context, number int) (*RandomRecipesResponse, error) {
	// Create request
	url := fmt.Sprintf("%s/recipes/random?number=%d", c.baseURL, number)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Add("x-api-key", c.apiKey)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read and parse response
	var randomRecipes RandomRecipesResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&randomRecipes); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &randomRecipes, nil
}

// GetRecipeInformation gets information about a specific recipe
// This method uses the official client as it doesn't have JSON parsing issues
func (c *Client) GetRecipeInformation(ctx context.Context, id int32) (*spoonacular.RecipeInformation, error) {
	result, _, err := c.apiClient.RecipesAPI.GetRecipeInformation(ctx, id).Execute()
	return result, err
}
