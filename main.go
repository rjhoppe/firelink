package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rjhoppe/firelink/bartender"
	"github.com/rjhoppe/firelink/books"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/database"
	"github.com/rjhoppe/firelink/models"

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
	dinner.InitializeClient()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// Check if cache.json exists and handle
	DrinkCache := cache.NewCache[models.DrinkResponse](15) // Set cache capacity
	DinnerCache := cache.NewCache[models.RecipeInfo](15)   // Set cache capacity

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
		dinner.GetRandomRecipes(c)
	})

	// Returns a specific recipe based on id
	r.GET("/dinner/recipe:id", func(c *gin.Context) {
		id := c.Param("id")
		dinner.GetRecipeFromApi(c, id, DinnerCache)
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
		bartender.GetRandomDrinkFromApi("", c, DrinkCache)
	})

	// Saves last drink recipe to DB
	r.POST("/bartender/save", func(c *gin.Context) {
		bartender.SaveDrinkToDB(c, DrinkCache)
	})

	// Returns the history for last 15 drinks
	r.GET("/bartender/history", func(c *gin.Context) {
		cachedDrinks := DrinkCache.GetAll()
		c.JSON(http.StatusOK, cachedDrinks)
	})

	r.GET("/bartender/:liquor", func(c *gin.Context) {
		liquor := c.Param("liquor")
		bartender.GetRandomDrinkFromApi(liquor, c, DrinkCache)
	})

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
