package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rjhoppe/firelink/bartender"
	"github.com/rjhoppe/firelink/books"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/database"
	"github.com/rjhoppe/firelink/models"
	"github.com/rjhoppe/firelink/ntfy"
	"github.com/rjhoppe/firelink/spoonacularapi"

	"github.com/rjhoppe/firelink/dinner"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	r := gin.Default()

	// Initialize dinner client
	apiKey := os.Getenv("SPOONACULAR_API_KEY")
	if apiKey == "" {
		fmt.Println("WARNING - SPOONACULAR_API_KEY is empty!")
	}

	apiClient := spoonacularapi.NewClient(apiKey)

	adapter := &spoonacularapi.SpoonacularAdapter{
		RealClient: apiClient,
	}

	// Initialize drink service
	drinkService := &bartender.DrinkService{}
	drinkService.GetDrinkFunc = drinkService.GetDrink
	drinkService.GatherIngredientsFunc = drinkService.GatherIngredients
	drinkService.Notifier = ntfy.NewNotifier("drink")

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	var DrinkCache *cache.Cache[models.DrinkResponse]
	var DinnerCache *cache.Cache[models.RecipeInfo]

	// Check if drink_cache.json exists and handle
	drinkCachePath := filepath.Join("/app/", "bartender", "cache.json")
	if _, err := os.Stat(drinkCachePath); os.IsNotExist(err) {
		DrinkCache = cache.NewCache[models.DrinkResponse](15)
	} else {
		DrinkCache, err = cache.RestoreCache[models.DrinkResponse](15, "bartender")
		if err != nil {
			DrinkCache = cache.NewCache[models.DrinkResponse](15)
		}
	}

	// Check if dinner_cache.json exists and handle
	dinnerCachePath := filepath.Join("/app/", "dinner", "cache.json")
	if _, err := os.Stat(dinnerCachePath); os.IsNotExist(err) {
		DinnerCache = cache.NewCache[models.RecipeInfo](15)
	} else {
		DinnerCache, err = cache.RestoreCache[models.RecipeInfo](15, "dinner")
		if err != nil {
			DinnerCache = cache.NewCache[models.RecipeInfo](15)
		}
	}

	// Returns a list of endpoints
	r.GET("/help", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"body": "List of endpoints..."})
	})

	// Checks if a book exists in the Gutenberg project
	r.GET("/ebook/find/:title", func(c *gin.Context) {
		title := c.Param("title")
		books.CheckForBook(c, title)
	})

	// r.POST("/ebook/dl/:title", func(c *gin.Context) {
	// 	title := c.Param("title")
	// 	books.GetBook(c, title)
	// })

	// Returns a random recipe
	r.GET("/dinner/random", func(c *gin.Context) {
		dinner.GetRandomRecipes(c, adapter)
	})

	// Returns a specific recipe based on id
	r.GET("/dinner/recipe:id", func(c *gin.Context) {
		id := c.Param("id")
		dinner.GetRecipeFromApi(c, id, DinnerCache, adapter)
	})

	// backup dinner cache
	r.POST("/dinner/cache/backup", func(c *gin.Context) {
		err := DinnerCache.BackupCache("/app/cache", DinnerCache.GetAll())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"body": "Dinner cache backup created successfully"})
		}
	})

	// Returns a random drink
	r.GET("/bartender/random", func(c *gin.Context) {
		drinkService.GetRandomDrinkFromApi("", c, DrinkCache)
	})

	// Saves last drink recipe to DB
	r.POST("/bartender/save", func(c *gin.Context) {
		drinkService.SaveDrinkToDB(c, DrinkCache)
	})

	// Returns the history for last 15 drinks
	r.GET("/bartender/history", func(c *gin.Context) {
		cachedDrinks := DrinkCache.GetAll()
		c.JSON(http.StatusOK, cachedDrinks)
	})

	// WIP: Get a drink of a specific liquor type from the API
	// r.GET("/bartender/:liquor", func(c *gin.Context) {
	// 	liquor := c.Param("liquor")
	// 	bartender.GetRandomDrinkFromApi(liquor, c, DrinkCache)
	// })

	// backup cache data
	r.POST("/bartender/cache/backup", func(c *gin.Context) {
		err := DrinkCache.BackupCache("/app/cache", DrinkCache.GetAll())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"body": "Drink cache backup created successfully"})
		}
	})

	// backup database
	r.POST("/database/backup", func(c *gin.Context) {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		err := database.BackupDB("db_backup_" + timestamp + ".sql")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"body": "Database backup created successfully"})
		}
	})

	r.Run(":8080")
}
