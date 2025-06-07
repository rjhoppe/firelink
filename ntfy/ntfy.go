package ntfy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rjhoppe/firelink/models"
)

// Notifier interface for sending notifications
type Notifier interface {
	SendMessage(title, message string) error
	SendFile(fileLoc string) error
}

// NtfyNotifier implements Notifier for ntfy.sh
type NtfyNotifier struct {
	Topic string
}

// Send sends a notification to the configured ntfy.sh topic
func (n *NtfyNotifier) SendMessage(title, message string) error {
	requestUrl := fmt.Sprintf("https://ntfy.rjhoppe.dev/%s", n.Topic)

	request, err := http.NewRequest("POST", requestUrl, strings.NewReader(message))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Title", title)

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

func (n *NtfyNotifier) SendFile(fileLoc string) error {
	file, err := os.Open(fileLoc)
	if err != nil {
		return err
	}
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", filepath.Base(fileLoc))
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", "https://ntfy.rjhoppe.dev/"+n.Topic, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ntfy returned status: %s", resp.Status)
	}

	return nil
}

// formatDrinkContent formats comma-separated ingredients into bullet points
func formatDrinkContent(content string) string {
	// Split by comma and create bullet points
	items := strings.Split(content, ",")
	var formatted []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			formatted = append(formatted, "• "+trimmed)
		}
	}
	return strings.Join(formatted, "\n")
}

// NtfyDrinkOfTheDay sends a drink notification using the Notifier interface
func NtfyDrinkOfTheDay(drink models.DrinkResponse, notifier Notifier) {
	// Format ingredients as bullet points
	formattedIngredients := formatDrinkContent(drink.Ingredients)

	msg := fmt.Sprintf(`%s

👀 Category: %v

🍸 Glass: %v

🛒 Ingredients:
%s

📝 Instructions:
%v`,
		drink.Name,
		drink.Category,
		drink.Glass,
		formattedIngredients,
		drink.Instructions)

	err := notifier.SendMessage("🍹 Drink of the Day", msg)
	if err != nil {
		log.Printf("Failed to send drink notification: %v", err)
	}
}

// NtfyRandomRecipes sends a random dinner notification using the Notifier interface
func NtfyRandomRecipes(recipeId int32, recipeName string, notifier Notifier) {
	msg := fmt.Sprintf(`
🍽️ Dinner %v: %v`, recipeId, recipeName)
	err := notifier.SendMessage("Recipe", msg)
	if err != nil {
		log.Printf("Failed to send dinner notification: %v", err)
	}
}

// NtfyRecipe sends a dinner recipe notification using the Notifier interface
func NtfyRecipe(recipe *models.RecipeInfo, notifier Notifier) {
	// Create a nicely formatted message with all information
	msg := fmt.Sprintf(`%s

📋 Recipe ID: %d

🛒 Ingredients:
%s

📝 Instructions:
%s

🌐 Source: %s`,
		recipe.Title,
		recipe.Id,
		recipe.Ingredients,
		recipe.Instructions,
		recipe.Url)

	err := notifier.SendMessage("Recipe Info! 🍽️", msg)
	if err != nil {
		log.Printf("Failed to send dinner notification: %v", err)
	}
}

func NtfyAllCacheDrinks(drinks []models.DrinkResponse, notifier Notifier) {
	var msg strings.Builder

	for i, drink := range drinks {
		msg.WriteString(fmt.Sprintf(
			"Drink %d: %v\nCategory: %v\nIngredients: %v\nInstructions: %v\n\n",
			i+1, drink.Name, drink.Category, drink.Ingredients, drink.Instructions,
		))
	}

	err := notifier.SendMessage("All Cached Drinks", msg.String())
	if err != nil {
		log.Printf("Failed to send drinks notification: %v", err)
	}
}

func NtfyDBBackup(fileLoc string, notifier Notifier) {
	err := notifier.SendFile(fileLoc)
	if err != nil {
		log.Printf("Failed to send db backup notification: %v", err)
	}

	_ = notifier.SendMessage("DB Backup", "DB Backup sent")
}
