package dinner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type JsonRandomRecipesResp struct {
	RecipeOne   string
	RecipeTwo   string
	RecipeThree string
}

type CustomRecipe struct {
	ID                   *int64   `json:"id,omitempty"`
	Title                *string  `json:"title,omitempty"`
	Image                *string  `json:"image,omitempty"`
	ImageType            *string  `json:"imageType,omitempty"`
	Servings             *int32   `json:"servings,omitempty"`
	ReadyInMinutes       *int32   `json:"readyInMinutes,omitempty"`
	License              *string  `json:"license,omitempty"`
	SourceName           *string  `json:"sourceName,omitempty"`
	SourceUrl            *string  `json:"sourceUrl,omitempty"`
	SpoonacularSourceUrl *string  `json:"spoonacularSourceUrl,omitempty"`
	AggregateLikes       *int32   `json:"aggregateLikes,omitempty"`
	HealthScore          *float32 `json:"healthScore,omitempty"`
	SpoonacularScore     *float32 `json:"spoonacularScore,omitempty"`
	PricePerServing      *float32 `json:"pricePerServing,omitempty"`
	NameClean            *string  `json:"nameClean,omitempty"`
	// Add any other fields that might be in the response
}

type CustomRandomRecipes struct {
	Recipes []CustomRecipe `json:"recipes"`
}

type CacheJson struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Servings       int32  `json:"servings"`
	ReadyInMinutes int32  `json:"readyInMinutes"`
	SourceName     string `json:"sourceName"`
	SourceUrl      string `json:"sourceUrl"`
}

func GetRandomRecipes(c *gin.Context, cache *Cache) {
	var result CustomRandomRecipes
	apiKey := os.Getenv("SPOONACULAR_API_KEY")
	client := &http.Client{}
	url := "https://api.spoonacular.com/recipes/random?number=3"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("x-api-key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	for i, recipe := range result.Recipes {
		fmt.Printf("Recipe %d:\n", i+1)
		if recipe.ID != nil {
			fmt.Printf("  ID: %d\n", *recipe.ID)
		}
		if recipe.Title != nil {
			fmt.Printf("  Title: %s\n", *recipe.Title)
		}
		if recipe.SourceName != nil {
			fmt.Printf("  Source Name: %s\n", *recipe.SourceName)
		}
		if recipe.SourceUrl != nil {
			fmt.Printf("  Source Url: %s\n", *recipe.SourceUrl)
		}
		if recipe.NameClean != nil {
			fmt.Printf("  Name Clean: %s\n", *recipe.NameClean)
		}
		if recipe.ReadyInMinutes != nil {
			fmt.Printf("  Ready in minutes: %d\n", *recipe.ReadyInMinutes)
		}
		if recipe.Servings != nil {
			fmt.Printf("  Servings: %d\n", *recipe.Servings)
		}

		ttl := time.Hour * 24 * 15
		recipeJson := CacheJson{
			ID:             *recipe.ID,
			Title:          *recipe.Title,
			Servings:       *recipe.Servings,
			ReadyInMinutes: *recipe.ReadyInMinutes,
			SourceName:     *recipe.SourceName,
			SourceUrl:      *recipe.SourceUrl,
		}

		strId := strconv.Itoa(int(*recipe.ID))
		fmt.Printf("%v, %v", strId, recipeJson)
		cache.Set(strId, recipeJson, ttl)
		fmt.Println()
	}

	fmtRecipeOne := fmt.Sprintf("%v: %v", strconv.Itoa(int(*result.Recipes[0].ID)), *result.Recipes[0].Title)
	fmtRecipeTwo := fmt.Sprintf("%v: %v", strconv.Itoa(int(*result.Recipes[1].ID)), *result.Recipes[1].Title)
	fmtRecipeThree := fmt.Sprintf("%v: %v", strconv.Itoa(int(*result.Recipes[2].ID)), *result.Recipes[2].Title)

	jsonResp := JsonRandomRecipesResp{
		RecipeOne:   fmtRecipeOne,
		RecipeTwo:   fmtRecipeTwo,
		RecipeThree: fmtRecipeThree,
	}

	c.JSON(http.StatusOK, jsonResp)
}

// Choose recipe from cache
func GetRecipe(c *gin.Context, id string, cache *Cache) {
	recipe, wasFound := cache.Get(id)
	if !wasFound {
		c.JSON(http.StatusNotFound, "Recipe not found in cache")
	}
	c.JSON(http.StatusOK, recipe)
}

func GetAllCachedRecipes(c *gin.Context, cache *Cache) {
	cachedRecipes := cache.GetAll()
	fmt.Println(cachedRecipes)
	c.JSON(http.StatusOK, cachedRecipes)
}

// TODO
func SaveRecipe() {

}
