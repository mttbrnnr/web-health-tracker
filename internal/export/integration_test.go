package export

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	// Create schema
	schema := `
	CREATE TABLE foods (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT,
		calories_per_serving REAL,
		protein_per_serving REAL,
		fat_per_serving REAL,
		carbs_per_serving REAL,
		serving_description TEXT,
		is_favorite BOOLEAN DEFAULT FALSE,
		source_line TEXT
	);

	CREATE TABLE custom_foods (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE daily_logs (
		id INTEGER PRIMARY KEY,
		date DATE NOT NULL UNIQUE,
		day_type TEXT CHECK(day_type IN ('rest', 'workout')) DEFAULT 'rest',
		weight_kg REAL,
		notes TEXT,
		synced_to_obsidian BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE meal_entries (
		id INTEGER PRIMARY KEY,
		daily_log_id INTEGER REFERENCES daily_logs(id),
		meal_type TEXT CHECK(meal_type IN ('breakfast', 'lunch', 'dinner', 'snacks', 'shake')),
		food_id INTEGER REFERENCES foods(id),
		custom_food_id INTEGER REFERENCES custom_foods(id),
		quantity REAL DEFAULT 1,
		saved_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	return db
}

func TestGetDayExportData_EmptyDay(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	export, err := GetDayExportData(db, "2026-03-15")
	if err != nil {
		t.Fatalf("GetDayExportData: %v", err)
	}

	if export.DayType != "rest" {
		t.Errorf("expected rest day type, got %s", export.DayType)
	}

	if export.WeightKg != nil {
		t.Errorf("expected nil weight, got %v", *export.WeightKg)
	}

	if len(export.Meals) != 5 {
		t.Errorf("expected 5 meals, got %d", len(export.Meals))
	}
}

func TestGetDayExportData_WithData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec(`
		INSERT INTO foods (id, name, category, calories_per_serving, protein_per_serving) VALUES
		(1, 'Coffee', 'breakfast', 5, 0),
		(2, 'Eggs', 'breakfast', 150, 12),
		(3, 'Chicken', 'lunch', 200, 30)
	`)
	if err != nil {
		t.Fatalf("insert foods: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO custom_foods (id, name, description) VALUES
		(1, 'Banana', 'large')
	`)
	if err != nil {
		t.Fatalf("insert custom foods: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO daily_logs (id, date, day_type, weight_kg) VALUES
		(1, '2026-03-15', 'workout', 82.5)
	`)
	if err != nil {
		t.Fatalf("insert daily log: %v", err)
	}

	now := time.Now()
	_, err = db.Exec(`
		INSERT INTO meal_entries (daily_log_id, meal_type, food_id, saved_at) VALUES
		(1, 'breakfast', 1, ?),
		(1, 'breakfast', 2, ?),
		(1, 'lunch', 3, ?)
	`, now, now, now)
	if err != nil {
		t.Fatalf("insert meal entries: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO meal_entries (daily_log_id, meal_type, custom_food_id, saved_at) VALUES
		(1, 'snacks', 1, ?)
	`, now)
	if err != nil {
		t.Fatalf("insert custom meal entry: %v", err)
	}

	// Get export data
	export, err := GetDayExportData(db, "2026-03-15")
	if err != nil {
		t.Fatalf("GetDayExportData: %v", err)
	}

	// Verify day type
	if export.DayType != "workout" {
		t.Errorf("expected workout day type, got %s", export.DayType)
	}

	// Verify weight
	if export.WeightKg == nil || *export.WeightKg != 82.5 {
		t.Errorf("expected weight 82.5, got %v", export.WeightKg)
	}

	// Verify breakfast
	var breakfast MealExport
	for _, m := range export.Meals {
		if m.MealType == "breakfast" {
			breakfast = m
			break
		}
	}
	if len(breakfast.Foods) != 2 {
		t.Errorf("expected 2 breakfast foods, got %d", len(breakfast.Foods))
	}
	if breakfast.Calories != 155 { // 5 + 150
		t.Errorf("expected breakfast calories 155, got %f", breakfast.Calories)
	}
	if breakfast.Protein != 12 { // 0 + 12
		t.Errorf("expected breakfast protein 12, got %f", breakfast.Protein)
	}

	// Verify snacks has custom food
	var snacks MealExport
	for _, m := range export.Meals {
		if m.MealType == "snacks" {
			snacks = m
			break
		}
	}
	if len(snacks.CustomFoods) != 1 {
		t.Errorf("expected 1 custom food in snacks, got %d", len(snacks.CustomFoods))
	}
	if snacks.CustomFoods[0] != "[Custom: Banana, large]" {
		t.Errorf("expected custom food format, got %s", snacks.CustomFoods[0])
	}

	// Verify totals
	if export.TotalCal != 355 { // 5 + 150 + 200
		t.Errorf("expected total cal 355, got %f", export.TotalCal)
	}
	if export.TotalProt != 42 { // 0 + 12 + 30
		t.Errorf("expected total protein 42, got %f", export.TotalProt)
	}
}

func TestExportDayToObsidian_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create temp vault directory
	vaultPath := t.TempDir()
	os.Setenv("OBSIDIAN_VAULT_PATH", vaultPath)
	defer os.Unsetenv("OBSIDIAN_VAULT_PATH")

	// Insert test data
	_, err := db.Exec(`
		INSERT INTO foods (id, name, category, calories_per_serving, protein_per_serving) VALUES
		(1, 'Coffee', 'breakfast', 5, 0),
		(2, 'Toast', 'breakfast', 120, 4)
	`)
	if err != nil {
		t.Fatalf("insert foods: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO daily_logs (id, date, day_type, weight_kg) VALUES
		(1, '2026-03-15', 'rest', 83.1)
	`)
	if err != nil {
		t.Fatalf("insert daily log: %v", err)
	}

	now := time.Now()
	_, err = db.Exec(`
		INSERT INTO meal_entries (daily_log_id, meal_type, food_id, saved_at) VALUES
		(1, 'breakfast', 1, ?),
		(1, 'breakfast', 2, ?)
	`, now, now)
	if err != nil {
		t.Fatalf("insert meal entries: %v", err)
	}

	// Export to Obsidian
	err = ExportDayToObsidian(db, "2026-03-15")
	if err != nil {
		t.Fatalf("ExportDayToObsidian: %v", err)
	}

	// Read the exported file
	filePath := filepath.Join(vaultPath, "02-Areas", "Health & Fitness", "Nutrition Logs", "2026-03 March.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read exported file: %v", err)
	}

	contentStr := string(content)

	// Verify content
	if !strings.Contains(contentStr, "# Nutrition Log - March 2026") {
		t.Errorf("missing header in:\n%s", contentStr)
	}

	if !strings.Contains(contentStr, "### Sunday, March 15") {
		t.Errorf("missing day header in:\n%s", contentStr)
	}

	if !strings.Contains(contentStr, "**Type:** Rest Day") {
		t.Errorf("missing day type in:\n%s", contentStr)
	}

	if !strings.Contains(contentStr, "**Weight:** 83.1kg") {
		t.Errorf("missing weight in:\n%s", contentStr)
	}

	if !strings.Contains(contentStr, "Coffee + Toast") {
		t.Errorf("missing breakfast foods in:\n%s", contentStr)
	}

	if !strings.Contains(contentStr, "125 kcal | 4g protein") {
		t.Errorf("missing breakfast macros in:\n%s", contentStr)
	}

	// Verify synced flag was set
	var synced bool
	err = db.QueryRow("SELECT synced_to_obsidian FROM daily_logs WHERE date = '2026-03-15'").Scan(&synced)
	if err != nil {
		t.Fatalf("query synced flag: %v", err)
	}
	if !synced {
		t.Error("expected synced_to_obsidian to be true")
	}

	// Test updating the same day
	_, err = db.Exec("UPDATE daily_logs SET weight_kg = 83.2 WHERE date = '2026-03-15'")
	if err != nil {
		t.Fatalf("update weight: %v", err)
	}

	err = ExportDayToObsidian(db, "2026-03-15")
	if err != nil {
		t.Fatalf("ExportDayToObsidian (update): %v", err)
	}

	// Read updated file
	content, err = os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read updated file: %v", err)
	}
	contentStr = string(content)

	// Should have new weight
	if !strings.Contains(contentStr, "**Weight:** 83.2kg") {
		t.Errorf("missing updated weight in:\n%s", contentStr)
	}

	// Should NOT have old weight (no duplicates)
	if strings.Contains(contentStr, "83.1kg") {
		t.Errorf("old weight still present in:\n%s", contentStr)
	}

	// Should only have one day header (no duplicates)
	count := strings.Count(contentStr, "### Sunday, March 15")
	if count != 1 {
		t.Errorf("expected 1 day header, found %d in:\n%s", count, contentStr)
	}
}

func TestExportDayToObsidian_NoVaultPath(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Ensure no vault path is set
	os.Unsetenv("OBSIDIAN_VAULT_PATH")

	// Should not error when vault path is not configured
	err := ExportDayToObsidian(db, "2026-03-15")
	if err != nil {
		t.Errorf("expected no error when vault path not set, got: %v", err)
	}
}
