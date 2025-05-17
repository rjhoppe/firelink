package models

import (
	"gorm.io/gorm"
)

type Dinner struct {
	gorm.Model
	Title        string
	ExternalId   string
	Url          string
	Instructions string
	Ingredients  string
}

type Drink struct {
	gorm.Model
	Name         string
	ExternalId   string
	Category     string
	Glass        string
	Ingredients  string
	Instructions string
}

type DrinkResponse struct {
	Message      string `json:"message"`
	ExternalId   string `json:"idDrink"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Glass        string `json:"glass"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
}

type GetRandomDrinkAPI struct {
	Drinks []struct {
		IDDrink         string `json:"idDrink"`
		StrDrink        string `json:"strDrink"`
		StrCategory     string `json:"strCategory"`
		StrGlass        string `json:"strGlass"`
		StrAlcoholic    string `json:"strAlcoholic"`
		StrInstructions string `json:"strInstructions"`
		StrIngredient1  string `json:"strIngredient1"`
		StrIngredient2  string `json:"strIngredient2"`
		StrIngredient3  string `json:"strIngredient3"`
		StrIngredient4  string `json:"strIngredient4"`
		StrIngredient5  string `json:"strIngredient5"`
		StrIngredient6  string `json:"strIngredient6"`
		StrIngredient7  string `json:"strIngredient7"`
		StrIngredient8  string `json:"strIngredient8"`
		StrIngredient9  string `json:"strIngredient9"`
		StrIngredient10 string `json:"strIngredient10"`
		StrIngredient11 string `json:"strIngredient11"`
		StrIngredient12 string `json:"strIngredient12"`
		StrIngredient13 string `json:"strIngredient13"`
		StrIngredient14 string `json:"strIngredient14"`
		StrIngredient15 string `json:"strIngredient15"`
		StrMeasure1     string `json:"strMeasure1"`
		StrMeasure2     string `json:"strMeasure2"`
		StrMeasure3     string `json:"strMeasure3"`
		StrMeasure4     string `json:"strMeasure4"`
		StrMeasure5     string `json:"strMeasure5"`
		StrMeasure6     string `json:"strMeasure6"`
		StrMeasure7     string `json:"strMeasure7"`
		StrMeasure8     string `json:"strMeasure8"`
		StrMeasure9     string `json:"strMeasure9"`
		StrMeasure10    string `json:"strMeasure10"`
		StrMeasure11    string `json:"strMeasure11"`
		StrMeasure12    string `json:"strMeasure12"`
		StrMeasure13    string `json:"strMeasure13"`
		StrMeasure14    string `json:"strMeasure14"`
		StrMeasure15    string `json:"strMeasure15"`
	} `json:"drinks"`
}

type RandomRecipes struct {
	RecipeOne   string `json:"recipe_one"`
	RecipeTwo   string `json:"recipe_two"`
	RecipeThree string `json:"recipe_three"`
}

type RecipeInfo struct {
	Title        string `json:"title"`
	Id           int32  `json:"id"`
	Url          string `json:"url"`
	Instructions string `json:"instructions"`
	Ingredients  string `json:"ingredients"`
}