package bartender

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/utils"
)

type DrinkResponse struct {
	Title        string `json:"message"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

type GetRandomDrinkAPI struct {
	Drinks []struct {
		IDDrink                     string `json:"idDrink"`
		StrDrink                    string `json:"strDrink"`
		StrDrinkAlternate           any    `json:"strDrinkAlternate"`
		StrTags                     any    `json:"strTags"`
		StrVideo                    any    `json:"strVideo"`
		StrCategory                 string `json:"strCategory"`
		StrIBA                      any    `json:"strIBA"`
		StrAlcoholic                string `json:"strAlcoholic"`
		StrGlass                    string `json:"strGlass"`
		StrInstructions             string `json:"strInstructions"`
		StrInstructionsES           any    `json:"strInstructionsES"`
		StrInstructionsDE           string `json:"strInstructionsDE"`
		StrInstructionsFR           any    `json:"strInstructionsFR"`
		StrInstructionsIT           string `json:"strInstructionsIT"`
		StrInstructionsZHHANS       any    `json:"strInstructionsZH-HANS"`
		StrInstructionsZHHANT       any    `json:"strInstructionsZH-HANT"`
		StrDrinkThumb               string `json:"strDrinkThumb"`
		StrIngredient1              string `json:"strIngredient1"`
		StrIngredient2              string `json:"strIngredient2"`
		StrIngredient3              string `json:"strIngredient3"`
		StrIngredient4              string `json:"strIngredient4"`
		StrIngredient5              string `json:"strIngredient5"`
		StrIngredient6              string `json:"strIngredient6"`
		StrIngredient7              string `json:"strIngredient7"`
		StrIngredient8              string `json:"strIngredient8"`
		StrIngredient9              string `json:"strIngredient9"`
		StrIngredient10             string `json:"strIngredient10"`
		StrIngredient11             string `json:"strIngredient11"`
		StrIngredient12             string `json:"strIngredient12"`
		StrIngredient13             string `json:"strIngredient13"`
		StrIngredient14             string `json:"strIngredient14"`
		StrIngredient15             string `json:"strIngredient15"`
		StrMeasure1                 string `json:"strMeasure1"`
		StrMeasure2                 string `json:"strMeasure2"`
		StrMeasure3                 string `json:"strMeasure3"`
		StrMeasure4                 string `json:"strMeasure4"`
		StrMeasure5                 string `json:"strMeasure5"`
		StrMeasure6                 string `json:"strMeasure6"`
		StrMeasure7                 string `json:"strMeasure7"`
		StrMeasure8                 string `json:"strMeasure8"`
		StrMeasure9                 string `json:"strMeasure9"`
		StrMeasure10                string `json:"strMeasure10"`
		StrMeasure11                string `json:"strMeasure11"`
		StrMeasure12                string `json:"strMeasure12"`
		StrMeasure13                string `json:"strMeasure13"`
		StrMeasure14                string `json:"strMeasure14"`
		StrMeasure15                string `json:"strMeasure15"`
		StrImageSource              any    `json:"strImageSource"`
		StrImageAttribution         any    `json:"strImageAttribution"`
		StrCreativeCommonsConfirmed string `json:"strCreativeCommonsConfirmed"`
		DateModified                string `json:"dateModified"`
	} `json:"drinks"`
}

type Ingredients struct {
	Ingredient1  string
	Ingredient2  string
	Ingredient3  string
	Ingredient4  string
	Ingredient5  string
	Ingredient6  string
	Ingredient7  string
	Ingredient8  string
	Ingredient9  string
	Ingredient10 string
	Ingredient11 string
	Ingredient12 string
	Ingredient13 string
	Ingredient14 string
	Ingredient15 string
}

func gatherIngredients(o GetRandomDrinkAPI) Ingredients {
	return Ingredients{
		Ingredient1:  o.Drinks[0].StrMeasure1 + " " + o.Drinks[0].StrIngredient1,
		Ingredient2:  o.Drinks[0].StrMeasure2 + " " + o.Drinks[0].StrIngredient2,
		Ingredient3:  o.Drinks[0].StrMeasure3 + " " + o.Drinks[0].StrIngredient3,
		Ingredient4:  o.Drinks[0].StrMeasure4 + " " + o.Drinks[0].StrIngredient4,
		Ingredient5:  o.Drinks[0].StrMeasure5 + " " + o.Drinks[0].StrIngredient5,
		Ingredient6:  o.Drinks[0].StrMeasure6 + " " + o.Drinks[0].StrIngredient6,
		Ingredient7:  o.Drinks[0].StrMeasure7 + " " + o.Drinks[0].StrIngredient7,
		Ingredient8:  o.Drinks[0].StrMeasure8 + " " + o.Drinks[0].StrIngredient8,
		Ingredient9:  o.Drinks[0].StrMeasure9 + " " + o.Drinks[0].StrIngredient9,
		Ingredient10: o.Drinks[0].StrMeasure10 + " " + o.Drinks[0].StrIngredient10,
		Ingredient11: o.Drinks[0].StrMeasure11 + " " + o.Drinks[0].StrIngredient11,
		Ingredient12: o.Drinks[0].StrMeasure12 + " " + o.Drinks[0].StrIngredient12,
		Ingredient13: o.Drinks[0].StrMeasure13 + " " + o.Drinks[0].StrIngredient13,
		Ingredient14: o.Drinks[0].StrMeasure14 + " " + o.Drinks[0].StrIngredient14,
		Ingredient15: o.Drinks[0].StrMeasure15 + " " + o.Drinks[0].StrIngredient15,
	}
}

func getDrink(c *gin.Context) GetRandomDrinkAPI {
	var drink GetRandomDrinkAPI

	url := "https://www.thecocktaildb.com/api/json/v1/1/random.php"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"body": "Error retrieving data from source api"})
	}
	if err := json.NewDecoder(resp.Body).Decode(&drink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"body": "Error decoding api response body"})
	}

	if drink.Drinks[0].StrAlcoholic != "Alcoholic" || drink.Drinks[0].StrCategory == "Beer" {
		fmt.Println("Invalid drink! Rerolling...")
		getDrink(c)
	}
	return drink
}

func GetRandomDrink(c *gin.Context, cache *Cache) {
	drink := getDrink(c)
	// parse values from json response body
	ingredientsData := gatherIngredients(drink)
	ingredientsList := utils.GetStructVals(ingredientsData)
	// "1 oz Tequila"
	jsonResp := DrinkResponse{
		Title:        "Drink of the Day",
		Name:         drink.Drinks[0].StrDrink,
		Category:     drink.Drinks[0].StrCategory,
		Ingredients:  ingredientsList,
		Instructions: drink.Drinks[0].StrInstructions,
	}
	// ttl working
	ttl := time.Hour * 24 * 15
	cache.Set(jsonResp.Name, jsonResp, ttl)
	c.JSON(http.StatusOK, jsonResp)
}
