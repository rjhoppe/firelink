package bartender

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/cache"
	"github.com/rjhoppe/firelink/database"
	"github.com/rjhoppe/firelink/models"
	"github.com/rjhoppe/firelink/ntfy"
	"github.com/rjhoppe/firelink/utils"
)

type DrinkService struct {
	GetDrinkFunc          func(string, *gin.Context) (models.GetRandomDrinkAPI, error)
	GatherIngredientsFunc func(models.GetRandomDrinkAPI) []string
	SaveDrinkToDBFunc     func(c *gin.Context, cache *cache.Cache[models.DrinkResponse])
	GetDrinkFromDBFunc    func(drinkName string, c *gin.Context, cache *cache.Cache[models.DrinkResponse])
	GetAllCacheDrinksFunc func(c *gin.Context, cache *cache.Cache[models.DrinkResponse])
	Notifier              ntfy.Notifier
}

// GatherIngredients builds a slice of non-empty ingredient strings.
func (s *DrinkService) GatherIngredients(drink models.GetRandomDrinkAPI) []string {
	if len(drink.Drinks) == 0 {
		return nil
	}
	d := drink.Drinks[0]
	ingredients := []string{}
	for i := 1; i <= 15; i++ {
		measure := utils.GetFieldValue(d, fmt.Sprintf("StrMeasure%d", i))
		ingredient := utils.GetFieldValue(d, fmt.Sprintf("StrIngredient%d", i))
		if ingredient != "" {
			ingredients = append(ingredients, fmt.Sprintf("%s %s", measure, ingredient))
		}
	}
	return ingredients
}

// GetDrink fetches a valid random drink from the API.
func (s *DrinkService) GetDrink(liquor string, c *gin.Context) (models.GetRandomDrinkAPI, error) {
	url := "https://www.thecocktaildb.com/api/json/v1/1/random.php"
	for {
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"body": "Error retrieving data from source api"})
			return models.GetRandomDrinkAPI{}, err
		}
		defer resp.Body.Close()

		var drink models.GetRandomDrinkAPI
		if err := json.NewDecoder(resp.Body).Decode(&drink); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"body": "Error decoding api response body"})
			return models.GetRandomDrinkAPI{}, err
		}

		if len(drink.Drinks) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"body": "No drinks found"})
			return models.GetRandomDrinkAPI{}, err
		}

		if liquor != "" {
			if drink.Drinks[0].StrAlcoholic == "Alcoholic" && drink.Drinks[0].StrCategory == liquor {
				return drink, nil
			}
		} else {
			if drink.Drinks[0].StrAlcoholic == "Alcoholic" && drink.Drinks[0].StrCategory != "Beer" {
				return drink, nil
			}
		}
		// Otherwise, loop again
	}
}

// GetRandomDrink handles the random drink endpoint.
func (s *DrinkService) GetRandomDrinkFromApi(liquor string, c *gin.Context, cache *cache.Cache[models.DrinkResponse]) {
	drink, err := s.GetDrinkFunc(liquor, c)
	if err != nil {
		return // Error already handled in getDrink
	}
	ingredients := s.GatherIngredientsFunc(drink)
	ingredientsList := strings.Join(ingredients, ", ")
	jsonResp := models.DrinkResponse{
		Message:      "Drink of the Day",
		ExternalId:   drink.Drinks[0].IDDrink,
		Name:         drink.Drinks[0].StrDrink,
		Category:     drink.Drinks[0].StrCategory,
		Glass:        drink.Drinks[0].StrGlass,
		Ingredients:  ingredientsList,
		Instructions: drink.Drinks[0].StrInstructions,
	}
	ttl := 15 * 24 * time.Hour
	cache.Set(jsonResp.Name, jsonResp, ttl)
	ntfy.NtfyDrinkOfTheDay(jsonResp, ntfy.NewNotifier("drink"))
	// s.Notifier.SendMessage("Drink of the Day", jsonResp.Name)
	c.JSON(http.StatusOK, jsonResp)
}

// Saves top drink record on the cache to db if not already present
func (s *DrinkService) SaveDrinkToDB(c *gin.Context, cache *cache.Cache[models.DrinkResponse]) {
	drink, found := cache.GetTop()
	if !found {
		c.JSON(404, gin.H{"error": "No drink found in cache"})
		return
	}

	// Check if drink already exists in DB by name
	var existing models.Drink
	db := database.GetDB()
	exists, err := database.CheckRecordExists(db, &existing)
	if err == nil && exists {
		c.JSON(200, gin.H{"message": "Drink already exists in database"})
		return
	}

	msg := fmt.Sprintf("Drink saved to database: %s", drink.Name)

	// Only create if not found
	database.SaveToDB(db, &models.Drink{
		Name:         drink.Name,
		Glass:        drink.Glass,
		Category:     drink.Category,
		Ingredients:  drink.Ingredients,
		Instructions: drink.Instructions,
	})
	c.JSON(201, gin.H{"message": msg})
}

func (s *DrinkService) GetDrinkFromDB(drinkName string, c *gin.Context, cache *cache.Cache[models.DrinkResponse]) {
	db := database.GetDB()
	var drink models.Drink
	db.Where("name = ?", drinkName).First(&drink)
	drinkResponse := models.DrinkResponse{
		Message:      "Drink of the Day",
		ExternalId:   drink.ExternalId,
		Name:         drink.Name,
		Category:     drink.Category,
		Glass:        drink.Glass,
		Ingredients:  drink.Ingredients,
		Instructions: drink.Instructions,
	}
	cache.Set(drink.Name, drinkResponse, 0)
	c.JSON(http.StatusOK, drinkResponse)
}

// GetAllCacheDrinks returns all drinks from the cache
func (s *DrinkService) GetAllCacheDrinks(c *gin.Context, cache *cache.Cache[models.DrinkResponse]) {
	allDrinks := cache.GetAll()
	drinkResponses := []models.DrinkResponse{}
	for _, drink := range allDrinks {
		drinkResponses = append(drinkResponses, drink)
	}
	ntfy.NtfyAllCacheDrinks(drinkResponses, ntfy.NewNotifier("drink"))
	c.JSON(http.StatusOK, drinkResponses)
}
