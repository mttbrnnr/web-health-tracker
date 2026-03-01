package db

import (
	"database/sql"
	"fmt"
)

type seedFood struct {
	Name               string
	Category           string
	CaloriesPerServing float64
	ProteinPerServing  float64
	FatPerServing      float64
	CarbsPerServing    float64
	ServingDescription string
	IsFavorite         bool
}

var presetFoods = []seedFood{
	// Breakfast
	{
		Name:               "Coffee (Oat Milk Flat White)",
		Category:           "breakfast",
		CaloriesPerServing: 72,
		ProteinPerServing:  1.7,
		FatPerServing:      3.3,
		CarbsPerServing:    7.8,
		ServingDescription: "1 cup (170ml)",
		IsFavorite:         true,
	},
	{
		Name:               "Pumpkin Bread + Cream Cheese + Salmon",
		Category:           "breakfast",
		CaloriesPerServing: 332,
		ProteinPerServing:  22.6,
		FatPerServing:      20.5,
		CarbsPerServing:    15,
		ServingDescription: "1 serving (bread 35g + cheese 30g + salmon 62.5g)",
		IsFavorite:         true,
	},
	{
		Name:               "Greek Yogurt + Granola",
		Category:           "breakfast",
		CaloriesPerServing: 388,
		ProteinPerServing:  20,
		FatPerServing:      15,
		CarbsPerServing:    45,
		ServingDescription: "1 bowl (yogurt 150g + granola 50g)",
		IsFavorite:         false,
	},

	// Lunch/Dinner
	{
		Name:               "Chicken Dürüm (Large)",
		Category:           "lunch",
		CaloriesPerServing: 690,
		ProteinPerServing:  50,
		FatPerServing:      25,
		CarbsPerServing:    65,
		ServingDescription: "1 large wrap",
		IsFavorite:         true,
	},
	{
		Name:               "Work Canteen Turkey + Rice",
		Category:           "lunch",
		CaloriesPerServing: 565,
		ProteinPerServing:  43,
		FatPerServing:      15,
		CarbsPerServing:    55,
		ServingDescription: "1 plate",
		IsFavorite:         false,
	},
	{
		Name:               "Weißwurst + Pretzel",
		Category:           "lunch",
		CaloriesPerServing: 500,
		ProteinPerServing:  18,
		FatPerServing:      28,
		CarbsPerServing:    45,
		ServingDescription: "1 sausage + 1 pretzel",
		IsFavorite:         false,
	},
	{
		Name:               "Gnocchi + Burrata",
		Category:           "dinner",
		CaloriesPerServing: 630,
		ProteinPerServing:  26,
		FatPerServing:      30,
		CarbsPerServing:    60,
		ServingDescription: "1 medium portion",
		IsFavorite:         false,
	},
	{
		Name:               "Pizza (Homemade)",
		Category:           "dinner",
		CaloriesPerServing: 695,
		ProteinPerServing:  37,
		FatPerServing:      25,
		CarbsPerServing:    75,
		ServingDescription: "3 squares (~3/8 of 30cm pizza)",
		IsFavorite:         false,
	},
	{
		Name:               "Pho Bo",
		Category:           "dinner",
		CaloriesPerServing: 500,
		ProteinPerServing:  25,
		FatPerServing:      15,
		CarbsPerServing:    60,
		ServingDescription: "1 bowl",
		IsFavorite:         false,
	},

	// Snacks
	{
		Name:               "Club Mate",
		Category:           "snacks",
		CaloriesPerServing: 100,
		ProteinPerServing:  0,
		FatPerServing:      0,
		CarbsPerServing:    24,
		ServingDescription: "1 bottle (500ml)",
		IsFavorite:         false,
	},
	{
		Name:               "Beer (Helles)",
		Category:           "snacks",
		CaloriesPerServing: 210,
		ProteinPerServing:  2,
		FatPerServing:      0,
		CarbsPerServing:    13,
		ServingDescription: "1 pint (500ml)",
		IsFavorite:         false,
	},
	{
		Name:               "Grapes",
		Category:           "snacks",
		CaloriesPerServing: 100,
		ProteinPerServing:  1,
		FatPerServing:      0,
		CarbsPerServing:    25,
		ServingDescription: "1 handful (150g)",
		IsFavorite:         false,
	},

	// Shakes
	{
		Name:               "Protein Shake (Single Scoop)",
		Category:           "shake",
		CaloriesPerServing: 111,
		ProteinPerServing:  24,
		FatPerServing:      1,
		CarbsPerServing:    2,
		ServingDescription: "30g powder + water",
		IsFavorite:         false,
	},
	{
		Name:               "Protein Shake (Double Scoop)",
		Category:           "shake",
		CaloriesPerServing: 222,
		ProteinPerServing:  48,
		FatPerServing:      2,
		CarbsPerServing:    4,
		ServingDescription: "60g powder + water",
		IsFavorite:         true,
	},
	{
		Name:               "Protein Shake (Double + Milk)",
		Category:           "shake",
		CaloriesPerServing: 318,
		ProteinPerServing:  53,
		FatPerServing:      7,
		CarbsPerServing:    12,
		ServingDescription: "60g powder + 150ml whole milk",
		IsFavorite:         true,
	},
}

// SeedFoods inserts preset foods into the database if the foods table is empty.
func SeedFoods(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM foods").Scan(&count); err != nil {
		return fmt.Errorf("check foods count: %w", err)
	}

	if count > 0 {
		return nil // Already seeded
	}

	stmt, err := db.Prepare(`
		INSERT INTO foods (name, category, calories_per_serving, protein_per_serving,
			fat_per_serving, carbs_per_serving, serving_description, is_favorite)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, food := range presetFoods {
		_, err := stmt.Exec(
			food.Name,
			food.Category,
			food.CaloriesPerServing,
			food.ProteinPerServing,
			food.FatPerServing,
			food.CarbsPerServing,
			food.ServingDescription,
			food.IsFavorite,
		)
		if err != nil {
			return fmt.Errorf("insert food %s: %w", food.Name, err)
		}
	}

	return nil
}
