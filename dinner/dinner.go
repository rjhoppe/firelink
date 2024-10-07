package dinner

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	spoonacular "github.com/ddsky/spoonacular-api-clients/go"
	"github.com/gin-gonic/gin"
)

type RandomRecipes struct {
	RecipeOne   string
	RecipeTwo   string
	RecipeThree string
}

type GetRecipeInfo struct {
	Title        string
	Id           int32
	Url          string
	Instructions string
	Ingredients  string
}

func apiInit() *spoonacular.APIClient {
	configuration := spoonacular.NewConfiguration()
	configuration.AddDefaultHeader("x-api-key", os.Getenv("SPOONACULAR_API_KEY"))
	client := spoonacular.NewAPIClient(configuration)
	return client
}

func GetRandomRecipes(c *gin.Context) {
	client := apiInit()
	result, _, err := client.RecipesAPI.GetRandomRecipes(context.Background()).
		Number(3).
		Execute()

	if err != nil {
		fmt.Printf("Error returning recipes from GetRandomRecipes endpoint from Spoonacular: %v", err)
	}

	fmtRecipeOne := fmt.Sprintf("%v: %v", result.Recipes[0].Id, result.Recipes[0].Title)
	fmtRecipeTwo := fmt.Sprintf("%v: %v", result.Recipes[1].Id, result.Recipes[1].Title)
	fmtRecipeThree := fmt.Sprintf("%v: %v", result.Recipes[2].Id, result.Recipes[2].Title)

	jsonResp := RandomRecipes{
		RecipeOne:   fmtRecipeOne,
		RecipeTwo:   fmtRecipeTwo,
		RecipeThree: fmtRecipeThree,
	}

	c.JSON(http.StatusOK, jsonResp)
}

func GetReipe(c *gin.Context, recipeId string) {
	var instructions string

	recipeIdInt64, err := strconv.ParseInt(recipeId, 10, 32)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}

	recipeIdInt32 := int32(recipeIdInt64)

	client := apiInit()
	result, _, err := client.RecipesAPI.GetRecipeInformation(context.Background(), recipeIdInt32).
		Execute()

	if err != nil {
		fmt.Printf("Error returning recipe from GetRecipeInformation endpoint from Spoonacular: %v", err)
	}

	if result.Instructions.IsSet() {
		instructions = ""
		instructionsPtr := result.Instructions.Get()
		instructionsStr := *instructionsPtr
		instructionsSlice := strings.Split(instructionsStr, ".")
		for i, step := range instructionsSlice {
			if trimmedStep := strings.TrimSpace(step); trimmedStep != "" {
				instructions = fmt.Sprintf("%v, %v", instructions, trimmedStep)
				i++
			}
		}
	}

	ingredients := ""
	ingredientList := result.GetExtendedIngredients()
	for i, ingredient := range ingredientList {
		ingredients = fmt.Sprintf("%v, %v%v of %v", ingredients, ingredient.Amount, ingredient.Unit, ingredient.Name)
		i++
	}

	jsonResp := GetRecipeInfo{
		Title:        result.Title,
		Id:           result.Id,
		Url:          result.SourceUrl,
		Instructions: instructions,
		Ingredients:  ingredients,
	}
	c.JSON(http.StatusOK, jsonResp)
}

// TODO
// func SaveRecipe() {

// }
