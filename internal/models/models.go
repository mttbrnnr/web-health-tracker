package models

import "time"

// Food represents a food item from the database.
type Food struct {
	ID                 int64
	Name               string
	Category           string
	CaloriesPerServing float64
	ProteinPerServing  float64
	FatPerServing      float64
	CarbsPerServing    float64
	ServingDescription string
	IsFavorite         bool
}

// DailyLog represents a day's log entry.
type DailyLog struct {
	ID               int64
	Date             string
	DayType          string
	WeightKg         *float64
	Notes            *string
	SyncedToObsidian bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// MealEntry represents a food logged to a specific meal.
type MealEntry struct {
	ID           int64
	DailyLogID   int64
	MealType     string
	FoodID       *int64
	CustomFoodID *int64
	Quantity     float64
	SavedAt      *time.Time
	CreatedAt    time.Time
}

// CustomFood represents a user-added custom food without macros.
type CustomFood struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

// CustomFoodCheckbox represents a custom food with its checked state for the UI.
type CustomFoodCheckbox struct {
	CustomFood CustomFood
	Checked    bool
}

// MealTypes are the valid meal type values.
var MealTypes = []string{"breakfast", "lunch", "dinner", "snacks", "shake"}

// MealSection holds a meal's foods and calculated totals for display.
type MealSection struct {
	MealType       string
	DisplayName    string
	Foods          []FoodCheckbox
	CustomFoods    []CustomFoodCheckbox
	TotalCalories  float64
	TotalProtein   float64
	HasSavedMeals  bool
}

// FoodCheckbox represents a food item with its checked state for the UI.
type FoodCheckbox struct {
	Food    Food
	Checked bool
}

// DayTotals holds the calculated totals for a day.
type DayTotals struct {
	Calories float64
	Protein  float64
}
