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
	// ===================
	// PROTEINS
	// ===================
	{
		Name:               "Smoked Trout Fillets (Full Container)",
		Category:           "proteins",
		CaloriesPerServing: 217,
		ProteinPerServing:  27,
		FatPerServing:      11.6,
		CarbsPerServing:    0.8,
		ServingDescription: "1 container (110g drained)",
		IsFavorite:         true,
	},
	{
		Name:               "Smoked Trout Fillets (2/3 Container)",
		Category:           "proteins",
		CaloriesPerServing: 144,
		ProteinPerServing:  17.8,
		FatPerServing:      7.7,
		CarbsPerServing:    0.5,
		ServingDescription: "2/3 container (~73g)",
		IsFavorite:         false,
	},
	{
		Name:               "Smoked Salmon (Half Pack)",
		Category:           "proteins",
		CaloriesPerServing: 156,
		ProteinPerServing:  15,
		FatPerServing:      10.6,
		CarbsPerServing:    0,
		ServingDescription: "half pack (62.5g)",
		IsFavorite:         true,
	},
	{
		Name:               "Smoked Salmon (Full Pack)",
		Category:           "proteins",
		CaloriesPerServing: 311,
		ProteinPerServing:  30,
		FatPerServing:      21.3,
		CarbsPerServing:    0,
		ServingDescription: "full pack (125g)",
		IsFavorite:         false,
	},
	{
		Name:               "Veggie Mince (Taco Serving)",
		Category:           "proteins",
		CaloriesPerServing: 383,
		ProteinPerServing:  71,
		FatPerServing:      8.8,
		CarbsPerServing:    1.9,
		ServingDescription: "100g",
		IsFavorite:         true,
	},
	{
		Name:               "Veggie Mince (Rice Bowl)",
		Category:           "proteins",
		CaloriesPerServing: 287,
		ProteinPerServing:  53,
		FatPerServing:      6.6,
		CarbsPerServing:    1.4,
		ServingDescription: "75g",
		IsFavorite:         false,
	},
	{
		Name:               "Chicken Breast Sliced (5 slices)",
		Category:           "proteins",
		CaloriesPerServing: 40,
		ProteinPerServing:  8,
		FatPerServing:      0.8,
		CarbsPerServing:    0.4,
		ServingDescription: "5 slices (~40g)",
		IsFavorite:         false,
	},
	{
		Name:               "Weißwurst",
		Category:           "proteins",
		CaloriesPerServing: 280,
		ProteinPerServing:  12,
		FatPerServing:      25,
		CarbsPerServing:    0,
		ServingDescription: "1 sausage (~100g)",
		IsFavorite:         false,
	},
	{
		Name:               "Egg (Fried/Boiled)",
		Category:           "proteins",
		CaloriesPerServing: 75,
		ProteinPerServing:  6,
		FatPerServing:      5,
		CarbsPerServing:    0,
		ServingDescription: "1 large egg",
		IsFavorite:         false,
	},

	// ===================
	// BREADS & GRAINS
	// ===================
	{
		Name:               "Pumpkin Seed Bread (Slice)",
		Category:           "breads",
		CaloriesPerServing: 93,
		ProteinPerServing:  3.7,
		FatPerServing:      2.6,
		CarbsPerServing:    12.3,
		ServingDescription: "1 slice (~35g)",
		IsFavorite:         true,
	},
	{
		Name:               "Pumpkin Seed Roll (Half)",
		Category:           "breads",
		CaloriesPerServing: 100,
		ProteinPerServing:  4,
		FatPerServing:      3,
		CarbsPerServing:    14,
		ServingDescription: "half roll (~40g)",
		IsFavorite:         false,
	},
	{
		Name:               "Pumpkin Seed Roll (Full)",
		Category:           "breads",
		CaloriesPerServing: 200,
		ProteinPerServing:  8,
		FatPerServing:      6,
		CarbsPerServing:    28,
		ServingDescription: "1 roll (~80g)",
		IsFavorite:         false,
	},
	{
		Name:               "Pretzel",
		Category:           "breads",
		CaloriesPerServing: 330,
		ProteinPerServing:  9,
		FatPerServing:      3,
		CarbsPerServing:    65,
		ServingDescription: "1 standard (~100g)",
		IsFavorite:         false,
	},
	{
		Name:               "Pretzel (Half)",
		Category:           "breads",
		CaloriesPerServing: 165,
		ProteinPerServing:  4.5,
		FatPerServing:      1.5,
		CarbsPerServing:    32.5,
		ServingDescription: "half pretzel",
		IsFavorite:         false,
	},
	{
		Name:               "Pretzel Bun",
		Category:           "breads",
		CaloriesPerServing: 220,
		ProteinPerServing:  6,
		FatPerServing:      2,
		CarbsPerServing:    44,
		ServingDescription: "1 bun (~80g)",
		IsFavorite:         false,
	},
	{
		Name:               "Bagel (Plain)",
		Category:           "breads",
		CaloriesPerServing: 280,
		ProteinPerServing:  10,
		FatPerServing:      1.5,
		CarbsPerServing:    54,
		ServingDescription: "1 bagel (100g)",
		IsFavorite:         false,
	},
	{
		Name:               "Bagel (Seed/Korn)",
		Category:           "breads",
		CaloriesPerServing: 251,
		ProteinPerServing:  9.4,
		FatPerServing:      3,
		CarbsPerServing:    45,
		ServingDescription: "1 bagel (85g)",
		IsFavorite:         false,
	},

	// ===================
	// DAIRY
	// ===================
	{
		Name:               "Cream Cheese",
		Category:           "dairy",
		CaloriesPerServing: 83,
		ProteinPerServing:  3.9,
		FatPerServing:      7.2,
		CarbsPerServing:    0.8,
		ServingDescription: "30g",
		IsFavorite:         true,
	},
	{
		Name:               "Greek Yogurt",
		Category:           "dairy",
		CaloriesPerServing: 174,
		ProteinPerServing:  15,
		FatPerServing:      14.1,
		CarbsPerServing:    6,
		ServingDescription: "150g serving",
		IsFavorite:         true,
	},
	{
		Name:               "Black Currant Yogurt (Full)",
		Category:           "dairy",
		CaloriesPerServing: 268,
		ProteinPerServing:  10,
		FatPerServing:      6,
		CarbsPerServing:    44,
		ServingDescription: "400g container",
		IsFavorite:         false,
	},
	{
		Name:               "Black Currant Yogurt (2/3)",
		Category:           "dairy",
		CaloriesPerServing: 179,
		ProteinPerServing:  6.7,
		FatPerServing:      4,
		CarbsPerServing:    29,
		ServingDescription: "2/3 container (~267g)",
		IsFavorite:         false,
	},
	{
		Name:               "Raspberry-Lemon Yogurt Drink",
		Category:           "dairy",
		CaloriesPerServing: 136,
		ProteinPerServing:  4.8,
		FatPerServing:      3,
		CarbsPerServing:    22,
		ServingDescription: "half container (200g)",
		IsFavorite:         false,
	},
	{
		Name:               "Burrata (Half Ball)",
		Category:           "dairy",
		CaloriesPerServing: 300,
		ProteinPerServing:  18,
		FatPerServing:      24,
		CarbsPerServing:    1,
		ServingDescription: "~100g",
		IsFavorite:         false,
	},
	{
		Name:               "Burrata (Typical Serving)",
		Category:           "dairy",
		CaloriesPerServing: 150,
		ProteinPerServing:  9,
		FatPerServing:      12,
		CarbsPerServing:    0.5,
		ServingDescription: "~50g",
		IsFavorite:         false,
	},

	// ===================
	// BREAKFAST ITEMS
	// ===================
	{
		Name:               "Coffee (Oat Milk Flat White)",
		Category:           "breakfast",
		CaloriesPerServing: 72,
		ProteinPerServing:  1.7,
		FatPerServing:      3.3,
		CarbsPerServing:    7.8,
		ServingDescription: "110ml oat milk + 60ml espresso",
		IsFavorite:         true,
	},
	{
		Name:               "Granola",
		Category:           "breakfast",
		CaloriesPerServing: 214,
		ProteinPerServing:  5,
		FatPerServing:      6,
		CarbsPerServing:    33,
		ServingDescription: "50g serving",
		IsFavorite:         false,
	},

	// ===================
	// BEVERAGES
	// ===================
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
		Name:               "Mio Mio Mate",
		Category:           "snacks",
		CaloriesPerServing: 83,
		ProteinPerServing:  0,
		FatPerServing:      0,
		CarbsPerServing:    20,
		ServingDescription: "1 bottle (330ml)",
		IsFavorite:         false,
	},
	{
		Name:               "Charitea Mate",
		Category:           "snacks",
		CaloriesPerServing: 56,
		ProteinPerServing:  0,
		FatPerServing:      0,
		CarbsPerServing:    14,
		ServingDescription: "1 bottle (330ml)",
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
		Name:               "Radler",
		Category:           "snacks",
		CaloriesPerServing: 200,
		ProteinPerServing:  1,
		FatPerServing:      0,
		CarbsPerServing:    20,
		ServingDescription: "1 pint (500ml)",
		IsFavorite:         false,
	},

	// ===================
	// GERMAN TRADITIONAL
	// ===================
	{
		Name:               "Schupfnudeln",
		Category:           "dinner",
		CaloriesPerServing: 280,
		ProteinPerServing:  5,
		FatPerServing:      2,
		CarbsPerServing:    58,
		ServingDescription: "typical portion (~200g)",
		IsFavorite:         false,
	},
	{
		Name:               "Sauerkraut",
		Category:           "dinner",
		CaloriesPerServing: 30,
		ProteinPerServing:  1,
		FatPerServing:      0,
		CarbsPerServing:    6,
		ServingDescription: "side (~150g)",
		IsFavorite:         false,
	},
	{
		Name:               "Speck (Bacon Bits)",
		Category:           "dinner",
		CaloriesPerServing: 100,
		ProteinPerServing:  6,
		FatPerServing:      8,
		CarbsPerServing:    0,
		ServingDescription: "30g",
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

	// ===================
	// RESTAURANT/TAKEOUT
	// ===================
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
		Name:               "Shawarma Plate",
		Category:           "lunch",
		CaloriesPerServing: 594,
		ProteinPerServing:  52,
		FatPerServing:      22,
		CarbsPerServing:    50,
		ServingDescription: "1 plate",
		IsFavorite:         true,
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
	{
		Name:               "Tonkotsu Ramen + Karaage",
		Category:           "dinner",
		CaloriesPerServing: 950,
		ProteinPerServing:  45,
		FatPerServing:      40,
		CarbsPerServing:    90,
		ServingDescription: "1 bowl with chicken",
		IsFavorite:         false,
	},
	{
		Name:               "Sushi (Big Dinner)",
		Category:           "dinner",
		CaloriesPerServing: 800,
		ProteinPerServing:  53,
		FatPerServing:      20,
		CarbsPerServing:    90,
		ServingDescription: "mixed salmon/tuna sashimi + rolls",
		IsFavorite:         false,
	},
	{
		Name:               "Pizza (Homemade, 3 squares)",
		Category:           "dinner",
		CaloriesPerServing: 695,
		ProteinPerServing:  37,
		FatPerServing:      25,
		CarbsPerServing:    75,
		ServingDescription: "3 squares (~3/8 of 30cm pizza)",
		IsFavorite:         false,
	},
	{
		Name:               "Gnocchi + Tomato + Burrata",
		Category:           "dinner",
		CaloriesPerServing: 630,
		ProteinPerServing:  26,
		FatPerServing:      30,
		CarbsPerServing:    60,
		ServingDescription: "medium portion",
		IsFavorite:         false,
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

	// ===================
	// PASTA & GRAINS
	// ===================
	{
		Name:               "Pasta (Cooked, Half Plate)",
		Category:           "dinner",
		CaloriesPerServing: 140,
		ProteinPerServing:  5,
		FatPerServing:      1,
		CarbsPerServing:    28,
		ServingDescription: "~100g cooked",
		IsFavorite:         false,
	},
	{
		Name:               "Pasta + Butter + Parmesan",
		Category:           "dinner",
		CaloriesPerServing: 280,
		ProteinPerServing:  8,
		FatPerServing:      12,
		CarbsPerServing:    32,
		ServingDescription: "half plate with toppings",
		IsFavorite:         false,
	},
	{
		Name:               "Rice (Cooked)",
		Category:           "dinner",
		CaloriesPerServing: 260,
		ProteinPerServing:  5,
		FatPerServing:      0.5,
		CarbsPerServing:    57,
		ServingDescription: "typical side (~200g)",
		IsFavorite:         false,
	},
	{
		Name:               "Tomato Rice",
		Category:           "dinner",
		CaloriesPerServing: 260,
		ProteinPerServing:  5,
		FatPerServing:      2,
		CarbsPerServing:    54,
		ServingDescription: "200g",
		IsFavorite:         false,
	},

	// ===================
	// VEGETABLES & SIDES
	// ===================
	{
		Name:               "Avocado",
		Category:           "snacks",
		CaloriesPerServing: 80,
		ProteinPerServing:  1,
		FatPerServing:      7.5,
		CarbsPerServing:    4,
		ServingDescription: "50g",
		IsFavorite:         false,
	},
	{
		Name:               "Green Beans (Cooked)",
		Category:           "dinner",
		CaloriesPerServing: 35,
		ProteinPerServing:  2,
		FatPerServing:      0,
		CarbsPerServing:    7,
		ServingDescription: "100g",
		IsFavorite:         false,
	},
	{
		Name:               "Carrots (Cooked)",
		Category:           "dinner",
		CaloriesPerServing: 35,
		ProteinPerServing:  1,
		FatPerServing:      0,
		CarbsPerServing:    8,
		ServingDescription: "100g",
		IsFavorite:         false,
	},
	{
		Name:               "Peas (Cooked)",
		Category:           "dinner",
		CaloriesPerServing: 80,
		ProteinPerServing:  5,
		FatPerServing:      0,
		CarbsPerServing:    14,
		ServingDescription: "100g",
		IsFavorite:         false,
	},

	// ===================
	// SNACKS & ADDITIONS
	// ===================
	{
		Name:               "Butter",
		Category:           "snacks",
		CaloriesPerServing: 45,
		ProteinPerServing:  0,
		FatPerServing:      5,
		CarbsPerServing:    0,
		ServingDescription: "5g pat",
		IsFavorite:         false,
	},
	{
		Name:               "Olive Oil",
		Category:           "snacks",
		CaloriesPerServing: 45,
		ProteinPerServing:  0,
		FatPerServing:      5,
		CarbsPerServing:    0,
		ServingDescription: "1 tsp (5g)",
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
	{
		Name:               "Banana Bread (Homemade)",
		Category:           "snacks",
		CaloriesPerServing: 150,
		ProteinPerServing:  3,
		FatPerServing:      4,
		CarbsPerServing:    25,
		ServingDescription: "1 slice (~80g)",
		IsFavorite:         false,
	},

	// ===================
	// BREAKFAST COMBOS
	// ===================
	{
		Name:               "Coffee + Bread + Cheese + Salmon",
		Category:           "breakfast",
		CaloriesPerServing: 404,
		ProteinPerServing:  24,
		FatPerServing:      23,
		CarbsPerServing:    21,
		ServingDescription: "coffee + bread 35g + cheese 30g + salmon 62.5g",
		IsFavorite:         true,
	},
	{
		Name:               "Coffee + Roll + Avocado + Trout",
		Category:           "breakfast",
		CaloriesPerServing: 341,
		ProteinPerServing:  25,
		FatPerServing:      18,
		CarbsPerServing:    22,
		ServingDescription: "coffee + half roll + avocado 50g + trout 82g",
		IsFavorite:         true,
	},
	{
		Name:               "Coffee + Greek Yogurt + Granola",
		Category:           "breakfast",
		CaloriesPerServing: 460,
		ProteinPerServing:  22,
		FatPerServing:      23,
		CarbsPerServing:    47,
		ServingDescription: "coffee + yogurt 150g + granola 50g",
		IsFavorite:         false,
	},

	// ===================
	// SHAKES
	// ===================
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
