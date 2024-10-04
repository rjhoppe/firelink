package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rjhoppe/firelink/bartender"
	"github.com/rjhoppe/firelink/books"
	"github.com/rjhoppe/firelink/dinner"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	username := os.Getenv("GIN_USERNAME")
	password := os.Getenv("GIN_PASSWORD")

	r := gin.Default()

	// Check if cache.json exists and handle
	DrinkCache := bartender.NewCache(15) // Set cache capacity

	r.POST("/help", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"body": "List of endpoints..."})
		})

	r.POST("/ebook/find/:title", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			title := c.Param("title")
			books.CheckForBook(c, title)
		})

	// r.POST("/ebook/dl/:title", gin.BasicAuth(gin.Accounts{
	// 	username: password}),
	// 	func(c *gin.Context) {
	// 		title := c.Param("title")
	// 		books.GetBook(c, title)
	// 	})

	r.POST("/dinner/random", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			dinner.GetRandomRecipes(c)
		})

	r.GET("/bartender/random", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			bartender.GetRandomDrink(c, DrinkCache)
		})

	// Saves last drink recipe to DB
	// r.POST("/bartender/save", gin.BasicAuth(gin.Accounts{
	// 	username: password}),
	// 	func(c *gin.Context) {
	// 		drink := c.Param("drink")
	// 		bartender.SaveDrink(c, liquor)
	// 	})

	// Returns the history for last 15 drinks
	r.GET("/bartender/history", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			cachedDrinks := DrinkCache.GetAll()
			c.JSON(http.StatusOK, cachedDrinks)
		})

	// r.POST("/bartender/:liquor", gin.BasicAuth(gin.Accounts{
	// 	username: password}),
	// 	func(c *gin.Context) {
	// 		liquor := c.Param("liquor")
	// 		bartender.MixMeADrink(c, liquor)
	// 	})

	r.POST("/network", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			return
		})

	// When ready to deploy ngrok

	// ctx := context.Background()
	// l, err := ngrok.Listen(ctx,
	// 	config.HTTPEndpoint(),
	// 	ngrok.WithAuthtokenFromEnv(),
	// )
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Printf("public address: %s\n", l.Addr())

	// if err := r.RunListener(l); err != nil {
	// 	log.Fatalln(err)
	// }

	// for testing without ngrok
	r.Run(":8080")
}
