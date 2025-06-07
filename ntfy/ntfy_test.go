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
		Name:         "Test Mojito",
		Category:     "Test Category",
		Glass:        "Test Glass",
		Ingredients:  "Test Ingredient 1, Test Ingredient 2",
		Instructions: "Test Instructions",
	}

	mockNotifier := &MockNotifier{}
	NtfyDrinkOfTheDay(drink, mockNotifier)

	time.Sleep(1 * time.Second)

	expectedMessage := `Test Mojito

👀 Category: Test Category

🍸 Glass: Test Glass

🛒 Ingredients:
• Test Ingredient 1
• Test Ingredient 2

📝 Instructions:
Test Instructions`

	assert.Equal(t, "🍹 Drink of the Day", mockNotifier.SentTitle)
	assert.Equal(t, expectedMessage, mockNotifier.SentMessage)
}

func TestNtfyDinner(t *testing.T) {
	recipe := models.RecipeInfo{
		Title:        "Test Recipe Title",
		Id:           123,
		Url:          "https://www.google.com",
		Instructions: "Test Instructions",
		Ingredients:  "Test Ingredient 1, Test Ingredient 2",
	}

	mockNotifier := &MockNotifier{}
	NtfyRecipe(&recipe, mockNotifier)

	time.Sleep(1 * time.Second)

	expectedMessage := `Test Recipe Title

📋 Recipe ID: 123

🛒 Ingredients:
Test Ingredient 1, Test Ingredient 2

📝 Instructions:
Test Instructions

🌐 Source: https://www.google.com`

	assert.Equal(t, "Recipe Info! 🍽️", mockNotifier.SentTitle)
	assert.Equal(t, expectedMessage, mockNotifier.SentMessage)
}

func TestNtfyRandomRecipes(t *testing.T) {
	mockNotifier := &MockNotifier{}
	NtfyRandomRecipes(123, "Recipe Title", mockNotifier)

	time.Sleep(1 * time.Second)

	assert.Equal(t, "Recipe", mockNotifier.SentTitle)
	assert.Equal(t, "\n🍽️ Dinner 123: Recipe Title", mockNotifier.SentMessage)
}

func TestNtfyDBBackup(t *testing.T) {
	mockNotifier := &MockNotifier{}
	NtfyDBBackup("test.txt", mockNotifier)

	time.Sleep(1 * time.Second)

	assert.Equal(t, "DB Backup", mockNotifier.SentTitle)
	assert.Equal(t, "DB Backup sent", mockNotifier.SentMessage)
}
