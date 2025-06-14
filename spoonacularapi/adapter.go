package spoonacularapi

import (
	"context"
)

type SpoonacularAdapter struct {
	RealClient *Client
}

func (a *SpoonacularAdapter) GetRecipeInformation(ctx context.Context, id int32) (*RecipeInformationOverride, error) {
	resp, err := a.RealClient.GetRecipeInformation(ctx, id)
	if err != nil {
		return nil, err
	}
	return ConvertToOverride(resp), nil
}

func (a *SpoonacularAdapter) GetRandomRecipes(ctx context.Context, count int) (*RandomRecipesResponse, error) {
	return a.RealClient.GetRandomRecipes(ctx, count)
}

func ConvertToOverride(resp *RecipeInformationResponse) *RecipeInformationOverride {
	if resp == nil {
		return nil
	}
	return &resp.Recipe
}
