package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"web-health-tracker/internal/db"
	"web-health-tracker/internal/models"
)

// MealSectionData holds data for rendering the meal section partial.
type MealSectionData struct {
	models.MealSection
	Date string
}

// MealResponseData holds both the meal section and day totals for OOB swap.
type MealResponseData struct {
	MealSection MealSectionData
	DayTotals   models.DayTotals
}

// SaveMeal handles POST /meals/:meal/save - saves checked food IDs and returns updated partial.
func SaveMeal(w http.ResponseWriter, r *http.Request) {
	mealType := chi.URLParam(r, "meal")
	dateStr := r.URL.Query().Get("date")

	if dateStr == "" {
		http.Error(w, "date parameter required", http.StatusBadRequest)
		return
	}

	// Validate meal type
	if !isValidMealType(mealType) {
		http.Error(w, "invalid meal type", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get food IDs from form
	var foodIDs []int64
	for _, idStr := range r.Form["food_ids"] {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		foodIDs = append(foodIDs, id)
	}

	// Get or create daily log
	dayLog, err := db.GetOrCreateDailyLog(database, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save meal entries
	if err := db.SaveMealEntries(database, dayLog.ID, mealType, foodIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the response
	renderMealResponse(w, dayLog, mealType, dateStr)
}

// SameAsYesterday handles POST /meals/:meal/yesterday - copies previous day's entries.
func SameAsYesterday(w http.ResponseWriter, r *http.Request) {
	mealType := chi.URLParam(r, "meal")
	dateStr := r.URL.Query().Get("date")

	if dateStr == "" {
		http.Error(w, "date parameter required", http.StatusBadRequest)
		return
	}

	// Validate meal type
	if !isValidMealType(mealType) {
		http.Error(w, "invalid meal type", http.StatusBadRequest)
		return
	}

	// Get previous day's food IDs
	foodIDs, err := db.GetPreviousDayMealEntries(database, dateStr, mealType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get or create daily log
	dayLog, err := db.GetOrCreateDailyLog(database, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the copied entries
	if err := db.SaveMealEntries(database, dayLog.ID, mealType, foodIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the response
	renderMealResponse(w, dayLog, mealType, dateStr)
}

// renderMealResponse renders the meal section partial with OOB day totals.
func renderMealResponse(w http.ResponseWriter, dayLog *models.DailyLog, mealType, dateStr string) {
	// Get all foods for this meal's category
	allFoods, err := db.GetAllFoods(database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Group foods by category
	foodsByCategory := make(map[string][]models.Food)
	for _, f := range allFoods {
		foodsByCategory[f.Category] = append(foodsByCategory[f.Category], f)
	}

	// Get checked food IDs for this meal
	checkedFoods, err := db.GetCheckedFoodIDs(database, dayLog.ID, mealType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build meal section
	section := models.MealSection{
		MealType:    mealType,
		DisplayName: mealDisplayNames[mealType],
	}

	categories := mealCategories[mealType]
	for _, category := range categories {
		for _, food := range foodsByCategory[category] {
			isChecked := checkedFoods[food.ID]
			section.Foods = append(section.Foods, models.FoodCheckbox{
				Food:    food,
				Checked: isChecked,
			})

			if isChecked {
				section.TotalCalories += food.CaloriesPerServing
				section.TotalProtein += food.ProteinPerServing
				section.HasSavedMeals = true
			}
		}
	}

	// Get custom foods for this meal
	customFoods, err := db.GetCustomFoodsForMeal(database, dayLog.ID, mealType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, cf := range customFoods {
		section.CustomFoods = append(section.CustomFoods, models.CustomFoodCheckbox{
			CustomFood: cf,
			Checked:    true,
		})
		section.HasSavedMeals = true
	}

	// Get day totals
	dayTotals, err := db.GetDayTotals(database, dayLog.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render meal section partial
	sectionData := MealSectionData{
		MealSection: section,
		Date:        dateStr,
	}

	if err := templates.ExecuteTemplate(w, "meal_section", sectionData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render OOB day totals
	if err := templates.ExecuteTemplate(w, "day_totals", dayTotals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// isValidMealType checks if a meal type is valid.
func isValidMealType(mealType string) bool {
	for _, mt := range models.MealTypes {
		if mt == mealType {
			return true
		}
	}
	return false
}

// CustomFoodCheckboxData holds data for rendering a single custom food checkbox.
type CustomFoodCheckboxData struct {
	CustomFood models.CustomFood
	Checked    bool
}

// AddCustomFood handles POST /foods/custom - creates a custom food and adds it to a meal.
func AddCustomFood(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("custom_name")
	description := r.FormValue("custom_description")
	mealType := r.FormValue("meal_type")
	dateStr := r.FormValue("date")

	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if !isValidMealType(mealType) {
		http.Error(w, "invalid meal type", http.StatusBadRequest)
		return
	}

	if dateStr == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}

	// Create the custom food
	customFoodID, err := db.CreateCustomFood(database, name, description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get or create daily log
	dayLog, err := db.GetOrCreateDailyLog(database, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add custom food to meal
	if err := db.AddCustomFoodToMeal(database, dayLog.ID, mealType, customFoodID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the created custom food to render
	customFood, err := db.GetCustomFood(database, customFoodID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render single checkbox (appended to the list)
	data := CustomFoodCheckboxData{
		CustomFood: *customFood,
		Checked:    true,
	}

	if err := templates.ExecuteTemplate(w, "custom_food_checkbox", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
