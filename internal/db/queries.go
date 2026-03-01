package db

import (
	"database/sql"
	"fmt"
	"time"

	"web-health-tracker/internal/models"
)

// GetOrCreateDailyLog retrieves the daily log for a date, creating one if it doesn't exist.
func GetOrCreateDailyLog(db *sql.DB, date string) (*models.DailyLog, error) {
	log := &models.DailyLog{}

	err := db.QueryRow(`
		SELECT id, date, day_type, weight_kg, notes, synced_to_obsidian, created_at, updated_at
		FROM daily_logs
		WHERE date = ?
	`, date).Scan(
		&log.ID, &log.Date, &log.DayType, &log.WeightKg, &log.Notes,
		&log.SyncedToObsidian, &log.CreatedAt, &log.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		result, err := db.Exec(`
			INSERT INTO daily_logs (date, day_type)
			VALUES (?, 'rest')
		`, date)
		if err != nil {
			return nil, fmt.Errorf("create daily log: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("get last insert id: %w", err)
		}

		log.ID = id
		log.Date = date
		log.DayType = "rest"
		log.CreatedAt = time.Now()
		log.UpdatedAt = time.Now()
		return log, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query daily log: %w", err)
	}

	return log, nil
}

// GetAllFoods retrieves all foods from the database.
func GetAllFoods(db *sql.DB) ([]models.Food, error) {
	rows, err := db.Query(`
		SELECT id, name, category, calories_per_serving, protein_per_serving,
			fat_per_serving, carbs_per_serving, serving_description, is_favorite
		FROM foods
		ORDER BY is_favorite DESC, category, name
	`)
	if err != nil {
		return nil, fmt.Errorf("query foods: %w", err)
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var f models.Food
		err := rows.Scan(
			&f.ID, &f.Name, &f.Category, &f.CaloriesPerServing, &f.ProteinPerServing,
			&f.FatPerServing, &f.CarbsPerServing, &f.ServingDescription, &f.IsFavorite,
		)
		if err != nil {
			return nil, fmt.Errorf("scan food: %w", err)
		}
		foods = append(foods, f)
	}

	return foods, rows.Err()
}

// GetFoodsByCategory retrieves foods filtered by category.
func GetFoodsByCategory(db *sql.DB, category string) ([]models.Food, error) {
	rows, err := db.Query(`
		SELECT id, name, category, calories_per_serving, protein_per_serving,
			fat_per_serving, carbs_per_serving, serving_description, is_favorite
		FROM foods
		WHERE category = ?
		ORDER BY is_favorite DESC, name
	`, category)
	if err != nil {
		return nil, fmt.Errorf("query foods by category: %w", err)
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var f models.Food
		err := rows.Scan(
			&f.ID, &f.Name, &f.Category, &f.CaloriesPerServing, &f.ProteinPerServing,
			&f.FatPerServing, &f.CarbsPerServing, &f.ServingDescription, &f.IsFavorite,
		)
		if err != nil {
			return nil, fmt.Errorf("scan food: %w", err)
		}
		foods = append(foods, f)
	}

	return foods, rows.Err()
}

// GetMealEntries retrieves all meal entries for a daily log.
func GetMealEntries(db *sql.DB, dailyLogID int64) ([]models.MealEntry, error) {
	rows, err := db.Query(`
		SELECT id, daily_log_id, meal_type, food_id, custom_food_id, quantity, saved_at, created_at
		FROM meal_entries
		WHERE daily_log_id = ?
	`, dailyLogID)
	if err != nil {
		return nil, fmt.Errorf("query meal entries: %w", err)
	}
	defer rows.Close()

	var entries []models.MealEntry
	for rows.Next() {
		var e models.MealEntry
		err := rows.Scan(
			&e.ID, &e.DailyLogID, &e.MealType, &e.FoodID, &e.CustomFoodID,
			&e.Quantity, &e.SavedAt, &e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan meal entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// GetMealEntriesByType retrieves meal entries for a specific meal type.
func GetMealEntriesByType(db *sql.DB, dailyLogID int64, mealType string) ([]models.MealEntry, error) {
	rows, err := db.Query(`
		SELECT id, daily_log_id, meal_type, food_id, custom_food_id, quantity, saved_at, created_at
		FROM meal_entries
		WHERE daily_log_id = ? AND meal_type = ?
	`, dailyLogID, mealType)
	if err != nil {
		return nil, fmt.Errorf("query meal entries by type: %w", err)
	}
	defer rows.Close()

	var entries []models.MealEntry
	for rows.Next() {
		var e models.MealEntry
		err := rows.Scan(
			&e.ID, &e.DailyLogID, &e.MealType, &e.FoodID, &e.CustomFoodID,
			&e.Quantity, &e.SavedAt, &e.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan meal entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// GetCheckedFoodIDs returns a set of food IDs that are checked for a given daily log and meal type.
func GetCheckedFoodIDs(db *sql.DB, dailyLogID int64, mealType string) (map[int64]bool, error) {
	entries, err := GetMealEntriesByType(db, dailyLogID, mealType)
	if err != nil {
		return nil, err
	}

	checked := make(map[int64]bool)
	for _, e := range entries {
		if e.FoodID != nil {
			checked[*e.FoodID] = true
		}
	}

	return checked, nil
}

// SaveMealEntries replaces all meal entries for a given daily log and meal type.
// It deletes existing entries and inserts new ones for the provided food IDs.
func SaveMealEntries(db *sql.DB, dailyLogID int64, mealType string, foodIDs []int64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing entries for this meal
	_, err = tx.Exec(`
		DELETE FROM meal_entries
		WHERE daily_log_id = ? AND meal_type = ?
	`, dailyLogID, mealType)
	if err != nil {
		return fmt.Errorf("delete existing entries: %w", err)
	}

	// Insert new entries
	now := time.Now()
	for _, foodID := range foodIDs {
		_, err = tx.Exec(`
			INSERT INTO meal_entries (daily_log_id, meal_type, food_id, quantity, saved_at, created_at)
			VALUES (?, ?, ?, 1, ?, ?)
		`, dailyLogID, mealType, foodID, now, now)
		if err != nil {
			return fmt.Errorf("insert meal entry: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetDayTotals calculates the total calories and protein for all saved meals in a daily log.
func GetDayTotals(db *sql.DB, dailyLogID int64) (models.DayTotals, error) {
	var totals models.DayTotals

	err := db.QueryRow(`
		SELECT COALESCE(SUM(f.calories_per_serving * me.quantity), 0),
		       COALESCE(SUM(f.protein_per_serving * me.quantity), 0)
		FROM meal_entries me
		JOIN foods f ON me.food_id = f.id
		WHERE me.daily_log_id = ? AND me.saved_at IS NOT NULL
	`, dailyLogID).Scan(&totals.Calories, &totals.Protein)

	if err != nil {
		return totals, fmt.Errorf("calculate day totals: %w", err)
	}

	return totals, nil
}

// GetPreviousDayMealEntries retrieves meal entries for the previous day's same meal type.
func GetPreviousDayMealEntries(db *sql.DB, currentDate string, mealType string) ([]int64, error) {
	// Parse the current date and subtract one day
	current, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		return nil, fmt.Errorf("parse date: %w", err)
	}
	previousDate := current.AddDate(0, 0, -1).Format("2006-01-02")

	rows, err := db.Query(`
		SELECT me.food_id
		FROM meal_entries me
		JOIN daily_logs dl ON me.daily_log_id = dl.id
		WHERE dl.date = ? AND me.meal_type = ? AND me.food_id IS NOT NULL
	`, previousDate, mealType)
	if err != nil {
		return nil, fmt.Errorf("query previous day entries: %w", err)
	}
	defer rows.Close()

	var foodIDs []int64
	for rows.Next() {
		var foodID int64
		if err := rows.Scan(&foodID); err != nil {
			return nil, fmt.Errorf("scan food id: %w", err)
		}
		foodIDs = append(foodIDs, foodID)
	}

	return foodIDs, rows.Err()
}

// UpdateDayWeight updates the weight for a daily log.
func UpdateDayWeight(db *sql.DB, dailyLogID int64, weightKg float64) error {
	_, err := db.Exec(`
		UPDATE daily_logs
		SET weight_kg = ?, updated_at = ?
		WHERE id = ?
	`, weightKg, time.Now(), dailyLogID)
	if err != nil {
		return fmt.Errorf("update weight: %w", err)
	}
	return nil
}

// UpdateDayType updates the day type for a daily log.
func UpdateDayType(db *sql.DB, dailyLogID int64, dayType string) error {
	_, err := db.Exec(`
		UPDATE daily_logs
		SET day_type = ?, updated_at = ?
		WHERE id = ?
	`, dayType, time.Now(), dailyLogID)
	if err != nil {
		return fmt.Errorf("update day type: %w", err)
	}
	return nil
}

// CreateCustomFood inserts a new custom food and returns its ID.
func CreateCustomFood(db *sql.DB, name, description string) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO custom_foods (name, description, created_at)
		VALUES (?, ?, ?)
	`, name, description, time.Now())
	if err != nil {
		return 0, fmt.Errorf("insert custom food: %w", err)
	}
	return result.LastInsertId()
}

// AddCustomFoodToMeal adds a custom food entry to a meal.
func AddCustomFoodToMeal(db *sql.DB, dailyLogID int64, mealType string, customFoodID int64) error {
	now := time.Now()
	_, err := db.Exec(`
		INSERT INTO meal_entries (daily_log_id, meal_type, custom_food_id, quantity, saved_at, created_at)
		VALUES (?, ?, ?, 1, ?, ?)
	`, dailyLogID, mealType, customFoodID, now, now)
	if err != nil {
		return fmt.Errorf("add custom food to meal: %w", err)
	}
	return nil
}

// GetCustomFoodsForMeal retrieves custom foods for a daily log and meal type.
func GetCustomFoodsForMeal(db *sql.DB, dailyLogID int64, mealType string) ([]models.CustomFood, error) {
	rows, err := db.Query(`
		SELECT cf.id, cf.name, cf.description, cf.created_at
		FROM custom_foods cf
		JOIN meal_entries me ON me.custom_food_id = cf.id
		WHERE me.daily_log_id = ? AND me.meal_type = ?
	`, dailyLogID, mealType)
	if err != nil {
		return nil, fmt.Errorf("query custom foods: %w", err)
	}
	defer rows.Close()

	var foods []models.CustomFood
	for rows.Next() {
		var f models.CustomFood
		if err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan custom food: %w", err)
		}
		foods = append(foods, f)
	}
	return foods, rows.Err()
}

// GetCustomFood retrieves a single custom food by ID.
func GetCustomFood(db *sql.DB, id int64) (*models.CustomFood, error) {
	var cf models.CustomFood
	err := db.QueryRow(`
		SELECT id, name, description, created_at
		FROM custom_foods
		WHERE id = ?
	`, id).Scan(&cf.ID, &cf.Name, &cf.Description, &cf.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get custom food: %w", err)
	}
	return &cf, nil
}
