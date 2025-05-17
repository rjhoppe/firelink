package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rjhoppe/firelink/bartender"
	"github.com/rjhoppe/firelink/books"
	"github.com/rjhoppe/firelink/cache"
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

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// Check if cache.json exists and handle
	DrinkCache := cache.NewCache[models.DrinkResponse](15) // Set cache capacity
	DinnerCache := cache.NewCache[models.RecipeInfo](15) // Set cache capacity

	// Returns a list of endpoints
	r.POST("/help", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"body": "List of endpoints..."})
	})

	// Checks if a book exists in the Gutenberg project
	r.POST("/ebook/find/:title", func(c *gin.Context) {
		title := c.Param("title")
		books.CheckForBook(c, title)
	})

	// r.POST("/ebook/dl/:title", func(c *gin.Context) {
	// 	title := c.Param("title")
	// 	books.GetBook(c, title)
	// })

	// Returns a random recipe
	r.POST("/dinner/random", func(c *gin.Context) {
		dinner.GetRandomRecipes(c)
	})

	// Returns a specific recipe based on id
	r.POST("/dinner/recipe:id", func(c *gin.Context) {
		id := c.Param("id")
		dinner.GetRecipe(c, id, DinnerCache)
	})

	// Returns a random drink
	r.GET("/bartender/random", func(c *gin.Context) {
		bartender.GetRandomDrink(c, DrinkCache)
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

	// r.POST("/bartender/:liquor", func(c *gin.Context) {
	// 	liquor := c.Param("liquor")
	// 	bartender.MixMeADrink(c, liquor)
	// })

	r.Run(":8080")
}
