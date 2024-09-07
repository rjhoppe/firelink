package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rjhoppe/firelink/bartender"
	"github.com/rjhoppe/firelink/books"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	r := gin.Default()
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

	r.POST("/ebook/dl/:title", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			title := c.Param("title")
			books.GetBook(c, title)
		})

	r.POST("/dinner/", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			return
		})

	r.POST("/bartender/:liquor", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			liquor := c.Param("liquor")
			bartender.MixMeADrink(c, liquor)
		})

	r.POST("/network", gin.BasicAuth(gin.Accounts{
		username: password}),
		func(c *gin.Context) {
			return
		})

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
