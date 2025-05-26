package spoonacularapi

type RecipeInformationOverride struct {
	ID                       int                   `json:"id"`
	Title                    string                `json:"title"`
	Image                    string                `json:"image"`
	ImageType                string                `json:"imageType"`
	Servings                 int                   `json:"servings"`
	ReadyInMinutes           int                   `json:"readyInMinutes"`
	CookingMinutes           int                   `json:"cookingMinutes"`
	PreparationMinutes       int                   `json:"preparationMinutes"`
	License                  string                `json:"license"`
	SourceName               string                `json:"sourceName"`
	SourceURL                string                `json:"sourceUrl"`
	SpoonacularSourceURL     string                `json:"spoonacularSourceUrl"`
	HealthScore              float64               `json:"healthScore"`
	SpoonacularScore         float64               `json:"spoonacularScore"`
	PricePerServing          float64               `json:"pricePerServing"`
	AnalyzedInstructions     []AnalyzedInstruction `json:"analyzedInstructions"`
	Cheap                    bool                  `json:"cheap"`
	CreditsText              string                `json:"creditsText"`
	Cuisines                 []string              `json:"cuisines"`
	DairyFree                bool                  `json:"dairyFree"`
	Diets                    []string              `json:"diets"`
	Gaps                     string                `json:"gaps"`
	GlutenFree               bool                  `json:"glutenFree"`
	Instructions             string                `json:"instructions"`
	Ketogenic                bool                  `json:"ketogenic"`
	LowFodmap                bool                  `json:"lowFodmap"`
	Occasions                []string              `json:"occasions"`
	Sustainable              bool                  `json:"sustainable"`
	Vegan                    bool                  `json:"vegan"`
	Vegetarian               bool                  `json:"vegetarian"`
	VeryHealthy              bool                  `json:"veryHealthy"`
	VeryPopular              bool                  `json:"veryPopular"`
	Whole30                  bool                  `json:"whole30"`
	WeightWatcherSmartPoints int                   `json:"weightWatcherSmartPoints"`
	DishTypes                []string              `json:"dishTypes"`
	ExtendedIngredients      []ExtendedIngredient  `json:"extendedIngredients"`
	Summary                  string                `json:"summary"`
	WinePairing              *WinePairing          `json:"winePairing,omitempty"`
}

type ExtendedIngredient struct {
	Aisle        string             `json:"aisle"`
	Amount       float64            `json:"amount"`
	Consistency  string             `json:"consistency"`
	ID           int                `json:"id"`
	Image        string             `json:"image"`
	Measures     IngredientMeasures `json:"measures"`
	Meta         []string           `json:"meta"`
	Name         string             `json:"name"`
	NameClean    string             `json:"nameClean"`
	Original     string             `json:"original"`
	OriginalName string             `json:"originalName"`
	Unit         string             `json:"unit"`
}

type IngredientMeasures struct {
	Metric Measure `json:"metric"`
	US     Measure `json:"us"`
}

type Measure struct {
	Amount    float64 `json:"amount"`
	UnitLong  string  `json:"unitLong"`
	UnitShort string  `json:"unitShort"`
}

type AnalyzedInstruction struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Number int    `json:"number"`
	Step   string `json:"step"`
	// Optionally: Ingredients, Equipment, etc.
}

type WinePairing struct {
	PairedWines    []string      `json:"pairedWines"`
	PairingText    string        `json:"pairingText"`
	ProductMatches []WineProduct `json:"productMatches"`
}

type WineProduct struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Price         string  `json:"price"`
	ImageURL      string  `json:"imageUrl"`
	AverageRating float64 `json:"averageRating"`
	RatingCount   float64 `json:"ratingCount"`
	Score         float64 `json:"score"`
	Link          string  `json:"link"`
}
