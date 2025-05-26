package ntfy

import (
	"testing"
	"time"

	"github.com/rjhoppe/firelink/models"
	"github.com/stretchr/testify/assert"
)

type MockNotifier struct {
	SentTitle   string
	SentMessage string
	SentFile    string
}

func (m *MockNotifier) SendMessage(title, message string) error {
	m.SentTitle = title
	m.SentMessage = message
	return nil
}

func (m *MockNotifier) SendFile(fileLoc string) error {
	m.SentFile = fileLoc
	return nil
}

func TestNtfyDrinkOfTheDay(t *testing.T) {
	drink := models.DrinkResponse{
		Name:         "Drink of the Day",
		Category:     "Test Category",
		Glass:        "Test Glass",
		Ingredients:  "Test Ingredient 1, Test Ingredient 2",
		Instructions: "Test Instructions",
	}

	mockNotifier := &MockNotifier{}
	NtfyDrinkOfTheDay(drink, mockNotifier)

	time.Sleep(1 * time.Second)

	assert.Equal(t, "Drink of the Day", mockNotifier.SentTitle)
	assert.Equal(t, "\nDrink of the Day: Drink of the Day\nCategory: Test Category\nGlass: Test Glass\nIngredients: Test Ingredient 1, Test Ingredient 2\nInstructions: Test Instructions", mockNotifier.SentMessage)
}

func TestNtfyDinner(t *testing.T) {
	recipe := models.RecipeInfo{
		Title:        "Recipe Title",
		Id:           123,
		Url:          "https://www.google.com",
		Instructions: "Test Instructions",
		Ingredients:  "Test Ingredient 1, Test Ingredient 2",
	}

	mockNotifier := &MockNotifier{}
	NtfyRecipe(&recipe, mockNotifier)

	time.Sleep(1 * time.Second)

	assert.Equal(t, "Recipe", mockNotifier.SentTitle)
	assert.Equal(t, "\nDinner Recipe Title: 123\nIngredients: Test Ingredient 1, Test Ingredient 2\nInstructions: Test Instructions\nUrl: https://www.google.com", mockNotifier.SentMessage)
}

func TestNtfyRandomRecipes(t *testing.T) {
	mockNotifier := &MockNotifier{}
	NtfyRandomRecipes(123, "Recipe Title", mockNotifier)

	time.Sleep(1 * time.Second)

	assert.Equal(t, "Recipe", mockNotifier.SentTitle)
	assert.Equal(t, "\nDinner 123: Recipe Title", mockNotifier.SentMessage)
}

func TestNtfyDBBackup(t *testing.T) {
	mockNotifier := &MockNotifier{}
	NtfyDBBackup("test.txt", mockNotifier)

	time.Sleep(1 * time.Second)

	assert.Equal(t, "DB Backup", mockNotifier.SentTitle)
	assert.Equal(t, "DB Backup sent", mockNotifier.SentMessage)
}
