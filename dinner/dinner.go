package dinner

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/database"
	"github.com/rjhoppe/firelink/models"
	"github.com/rjhoppe/firelink/ntfy"
	"github.com/rjhoppe/firelink/spoonacularapi"
)

type SpoonacularClient interface {
	GetRandomRecipes(ctx context.Context, count int) (*spoonacularapi.RandomRecipesResponse, error)
	GetRecipeInformation(ctx context.Context, id int32) (*spoonacularapi.RecipeInformationOverride, error)
}

func GetRandomRecipes(c *gin.Context, apiClient SpoonacularClient) {
	result, err := apiClient.GetRandomRecipes(context.Background(), 3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching random recipes: %v", err)})
		return
	}

	recipes := result.Recipes
	if len(recipes) < 3 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not enough recipes returned"})
		return
	}

	for _, recipe := range recipes {
		ntfy.NtfyRandomRecipes(recipe.Id, recipe.Title, ntfy.NewNotifier("dinner"))
	}

	jsonResp := models.RandomRecipes{
		RecipeOne:   fmt.Sprintf("%v: %v", recipes[0].Id, recipes[0].Title),
		RecipeTwo:   fmt.Sprintf("%v: %v", recipes[1].Id, recipes[1].Title),
		RecipeThree: fmt.Sprintf("%v: %v", recipes[2].Id, recipes[2].Title),
	}

	c.JSON(http.StatusOK, jsonResp)
}

func GetRecipeFromApi(c *gin.Context, recipeId string, cache *cache.Cache[models.RecipeInfo], apiClient SpoonacularClient) {
	cacheRecipe, found := cache.Get(recipeId)
	if found {
		c.JSON(http.StatusOK, cacheRecipe)
		return
	}

	recipeIdInt64, err := strconv.ParseInt(recipeId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	result, err := apiClient.GetRecipeInformation(context.Background(), int32(recipeIdInt64))
	if err != nil {
		// Try to print the raw response if possible (if your client exposes it)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching recipe: %v", err)})
		return
	}

	// Ingredients
	var ingredients []string
	for _, ingredient := range result.ExtendedIngredients {
		ingredients = append(ingredients, fmt.Sprintf("%v%v of %v", ingredient.Amount, ingredient.Unit, ingredient.Name))
	}
	ingredientsStr := strings.Join(ingredients, ", ")

	ttl := 15 * 24 * time.Hour
	data := models.RecipeInfo{
		Title:        result.Title,
		Id:           int32(result.ID),
		Url:          result.SourceName,
		Instructions: result.Instructions,
		Ingredients:  ingredientsStr,
	}
	cache.Set(recipeId, data, ttl)
	ntfy.NtfyRecipe(&data, ntfy.NewNotifier("dinner"))
	c.JSON(http.StatusOK, data)
}

func GetRecipeFromDB(c *gin.Context, recipeId string, cache *cache.Cache[models.RecipeInfo]) {
	cacheRecipe, found := cache.Get(recipeId)
	if found {
		c.JSON(http.StatusOK, cacheRecipe)
		return
	}

	db := database.GetDB()
	var recipe models.Dinner
	db.Where("external_id = ?", recipeId).First(&recipe)
	recipeIdInt, err := strconv.Atoi(recipe.ExternalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error converting recipe ID: %v", err)})
		return
	}

	recipeInfo := models.RecipeInfo{
		Title:        recipe.Title,
		Id:           int32(recipeIdInt),
		Url:          recipe.Url,
		Instructions: recipe.Instructions,
		Ingredients:  recipe.Ingredients,
	}
	cache.Set(recipeId, recipeInfo, 0)
	c.JSON(http.StatusOK, recipeInfo)
}

func SaveRecipe(c *gin.Context, cache *cache.Cache[models.RecipeInfo], recipe *models.RecipeInfo) {
	// Check if drink already exists in DB by name
	var existing models.Dinner
	db := database.GetDB()
	exists, err := database.CheckRecordExists(db, &existing)
	if err == nil && exists {
		c.JSON(200, gin.H{"message": "Dinner recipe already exists in database"})
		return
	}

	database.SaveToDB(db, &models.Dinner{
		Title:        recipe.Title,
		ExternalId:   strconv.Itoa(int(recipe.Id)),
		Url:          recipe.Url,
		Instructions: recipe.Instructions,
		Ingredients:  recipe.Ingredients,
	})
}
