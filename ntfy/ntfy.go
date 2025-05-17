package ntfy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rjhoppe/firelink/models"
)

// Notifier interface for sending notifications
type Notifier interface {
	Send(title, message string) error
}

// NtfyNotifier implements Notifier for ntfy.sh
type NtfyNotifier struct {
	Topic string
}

// Send sends a notification to the configured ntfy.sh topic
func (n *NtfyNotifier) Send(title, message string) error {
	payload := map[string]string{
		"topic": n.Topic,
		"title": title,
		"body":  message,
	}

	requestUrl := fmt.Sprintf("https://ntfy.rjhoppe.dev/%s", n.Topic)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return err
	}

	request, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return err
	}
	defer response.Body.Close()

	log.Printf("Response status: %v", response.Status)
	return nil
}

// NewNotifier is a factory function for Notifier
func NewNotifier(topic string) Notifier {
	return &NtfyNotifier{Topic: topic}
}

// NtfyDrinkOfTheDay sends a drink notification using the Notifier interface
func NtfyDrinkOfTheDay(drink models.DrinkResponse) {
	notifier := NewNotifier("drink")
	msg := fmt.Sprintf(`
Drink of the Day: %v
Category: %v
Glass: %v
Ingredients: %v
Instructions: %v`, drink.Name, drink.Category, drink.Glass, drink.Ingredients, drink.Instructions)
	err := notifier.Send("Drink of the Day", msg)
	if err != nil {
		log.Printf("Failed to send drink notification: %v", err)
	}
}

// NtfyRandomRecipes sends a random dinner notification using the Notifier interface
func NtfyRandomRecipes(recipeId int32, recipeName string) {
	notifier := NewNotifier("dinner")
	msg := fmt.Sprintf(`
Dinner %v: %v`, recipeId, recipeName)
	err := notifier.Send("Recipe", msg)
	if err != nil {
		log.Printf("Failed to send dinner notification: %v", err)
	}
}

// NtfyRecipe sends a dinner recipe notification using the Notifier interface
func NtfyRecipe(recipe *models.RecipeInfo) {
	notifier := NewNotifier("dinner")
	msg := fmt.Sprintf(`
Dinner %v: %v
Ingredients: %v
Instructions: %v
Url: %v`, recipe.Title, recipe.Id, recipe.Ingredients, recipe.Instructions, recipe.Url)
	err := notifier.Send("Recipe", msg)
	if err != nil {
		log.Printf("Failed to send dinner notification: %v", err)
	}
}
