package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"web-health-tracker/internal/db"
	"web-health-tracker/internal/models"
)

var database *sql.DB

// SetDB sets the database connection for handlers.
func SetDB(d *sql.DB) {
	database = d
}

// mealDisplayNames maps meal types to display names.
var mealDisplayNames = map[string]string{
	"breakfast": "Breakfast",
	"lunch":     "Lunch",
	"dinner":    "Dinner",
	"snacks":    "Snacks",
	"shake":     "Shake",
}

// mealCategories maps meal types to their food categories.
// Some meals show foods from multiple categories.
var mealCategories = map[string][]string{
	"breakfast": {"breakfast"},
	"lunch":     {"lunch"},
	"dinner":    {"dinner"},
	"snacks":    {"snacks"},
	"shake":     {"shake"},
}

// DayViewData holds all data needed to render the day view.
type DayViewData struct {
	Title         string
	Date          string
	DisplayDate   string
	PrevDate      string
	NextDate      string
	DayLog        *models.DailyLog
	WeightDisplay string
	MealSections  []models.MealSection
	DayTotals     models.DayTotals
}

// Day handles GET /day/:date - renders the day view page.
func Day(w http.ResponseWriter, r *http.Request) {
	dateStr := chi.URLParam(r, "date")

	// Parse and validate date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Default to today if invalid date
		date = time.Now()
		dateStr = date.Format("2006-01-02")
	}

	// Get or create daily log
	dayLog, err := db.GetOrCreateDailyLog(database, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all foods
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

	// Get all meal entries for this day
	allEntries, err := db.GetMealEntries(database, dayLog.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build set of checked food IDs per meal type
	checkedByMeal := make(map[string]map[int64]bool)
	for _, mealType := range models.MealTypes {
		checkedByMeal[mealType] = make(map[int64]bool)
	}
	for _, entry := range allEntries {
		if entry.FoodID != nil {
			checkedByMeal[entry.MealType][*entry.FoodID] = true
		}
	}

	// Build meal sections
	var mealSections []models.MealSection
	var totalCalories, totalProtein float64

	for _, mealType := range models.MealTypes {
		section := models.MealSection{
			MealType:    mealType,
			DisplayName: mealDisplayNames[mealType],
		}

		// Get foods for this meal's categories
		categories := mealCategories[mealType]
		checkedFoods := checkedByMeal[mealType]

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

		totalCalories += section.TotalCalories
		totalProtein += section.TotalProtein
		mealSections = append(mealSections, section)
	}

	// Format dates for navigation
	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	// Format display date
	displayDate := date.Format("Monday, Jan 2")

	// Short date for nav links
	prevShort := date.AddDate(0, 0, -1).Format("Jan 2")
	nextShort := date.AddDate(0, 0, 1).Format("Jan 2")

	// Format weight for display
	var weightDisplay string
	if dayLog.WeightKg != nil {
		weightDisplay = strconv.FormatFloat(*dayLog.WeightKg, 'f', 1, 64)
	}

	data := DayViewData{
		Title:         "Day View - " + displayDate,
		Date:          dateStr,
		DisplayDate:   displayDate,
		PrevDate:      prevDate,
		NextDate:      nextDate,
		DayLog:        dayLog,
		WeightDisplay: weightDisplay,
		MealSections:  mealSections,
		DayTotals: models.DayTotals{
			Calories: totalCalories,
			Protein:  totalProtein,
		},
	}

	// Add short dates to data for template
	type extendedData struct {
		DayViewData
		PrevShort string
		NextShort string
	}

	fullData := extendedData{
		DayViewData: data,
		PrevShort:   prevShort,
		NextShort:   nextShort,
	}

	if err := templates.ExecuteTemplate(w, "day.html", fullData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// formatMealType converts meal type to title case for display.
func formatMealType(mealType string) string {
	if mealType == "" {
		return ""
	}
	return strings.ToUpper(mealType[:1]) + mealType[1:]
}

// WeightInputData holds data for rendering the weight input partial.
type WeightInputData struct {
	Date        string
	WeightValue string
}

// UpdateWeight handles POST /day/:date/weight - saves weight and returns updated input.
func UpdateWeight(w http.ResponseWriter, r *http.Request) {
	dateStr := chi.URLParam(r, "date")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weightStr := r.FormValue("weight_kg")
	var weightKg *float64

	if weightStr != "" {
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			http.Error(w, "invalid weight value", http.StatusBadRequest)
			return
		}
		weightKg = &weight
	}

	// Get or create daily log
	dayLog, err := db.GetOrCreateDailyLog(database, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update weight if provided
	if weightKg != nil {
		if err := db.UpdateDayWeight(database, dayLog.ID, *weightKg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Render updated weight input
	var weightValue string
	if weightKg != nil {
		weightValue = strconv.FormatFloat(*weightKg, 'f', 1, 64)
	}
	data := WeightInputData{
		Date:        dateStr,
		WeightValue: weightValue,
	}

	if err := templates.ExecuteTemplate(w, "weight_input", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DayTypeData holds data for rendering the day type select partial.
type DayTypeData struct {
	Date    string
	DayType string
}

// UpdateDayType handles POST /day/:date/type - saves day type and returns updated select.
func UpdateDayType(w http.ResponseWriter, r *http.Request) {
	dateStr := chi.URLParam(r, "date")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dayType := r.FormValue("day_type")
	if dayType != "rest" && dayType != "workout" {
		http.Error(w, "invalid day type", http.StatusBadRequest)
		return
	}

	// Get or create daily log
	dayLog, err := db.GetOrCreateDailyLog(database, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update day type
	if err := db.UpdateDayType(database, dayLog.ID, dayType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render updated day type select
	data := DayTypeData{
		Date:    dateStr,
		DayType: dayType,
	}

	if err := templates.ExecuteTemplate(w, "day_type_select", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
