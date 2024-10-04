package dinner

import (
	"context"
	"fmt"
	"net/http"
	"os"

	spoonacular "github.com/ddsky/spoonacular-api-clients/go"
	"github.com/gin-gonic/gin"
)

type RandomRecipes struct {
	RecipeOne   string
	RecipeTwo   string
	RecipeThree string
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

	fmtRecipeOne := fmt.Sprintf("1: %v (%v)", result.Recipes[0].Title, result.Recipes[0].Id)
	fmtRecipeTwo := fmt.Sprintf("2: %v (%v)", result.Recipes[1].Title, result.Recipes[1].Id)
	fmtRecipeThree := fmt.Sprintf("3: %v (%v)", result.Recipes[2].Title, result.Recipes[2].Id)

	jsonResp := RandomRecipes{
		RecipeOne:   fmtRecipeOne,
		RecipeTwo:   fmtRecipeTwo,
		RecipeThree: fmtRecipeThree,
	}

	c.JSON(http.StatusOK, jsonResp)
}

func GetReipe(recipeId int32) {
	client := apiInit()
	result, _, err := client.RecipesAPI.GetRecipeInformation(context.Background(), recipeId).
		Execute()

	if err != nil {
		fmt.Printf("Error returning recipe from GetRecipeInformation endpoint from Spoonacular: %v", err)
	}

	// ingredients := result.ExtendedIngredients

}

func SaveMeal() {

}
