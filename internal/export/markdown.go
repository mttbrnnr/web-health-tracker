package export

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MealExport holds exported data for a single meal.
type MealExport struct {
	MealType     string
	DisplayName  string
	Foods        []string  // standard food names
	CustomFoods  []string  // formatted as "[Custom: name, description]"
	Calories     float64
	Protein      float64
	HasSavedMeal bool
}

// DayExport holds all data needed to export a day's entry.
type DayExport struct {
	Date      time.Time
	DayType   string
	WeightKg  *float64
	Meals     []MealExport
	TotalCal  float64
	TotalProt float64
}

// mealOrder defines the order of meals in the export.
var mealOrder = []string{"breakfast", "lunch", "snacks", "dinner", "shake"}

// mealDisplayNames maps meal types to display names.
var mealDisplayNames = map[string]string{
	"breakfast": "Breakfast",
	"lunch":     "Lunch",
	"dinner":    "Dinner",
	"snacks":    "Snacks",
	"shake":     "Shake",
}

// GetDayExportData queries all data for a given date from SQLite.
func GetDayExportData(db *sql.DB, dateStr string) (*DayExport, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("parse date: %w", err)
	}

	// Get daily log
	var dayType string
	var weightKg *float64

	err = db.QueryRow(`
		SELECT day_type, weight_kg
		FROM daily_logs
		WHERE date = ?
	`, dateStr).Scan(&dayType, &weightKg)

	if err == sql.ErrNoRows {
		// No log for this day, return empty export
		return &DayExport{
			Date:    date,
			DayType: "rest",
			Meals:   buildEmptyMeals(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query daily log: %w", err)
	}

	// Query all meal entries with food names
	rows, err := db.Query(`
		SELECT
			me.meal_type,
			COALESCE(f.name, '') as food_name,
			COALESCE(f.calories_per_serving, 0) as calories,
			COALESCE(f.protein_per_serving, 0) as protein,
			COALESCE(cf.name, '') as custom_name,
			COALESCE(cf.description, '') as custom_desc,
			me.saved_at IS NOT NULL as is_saved
		FROM meal_entries me
		LEFT JOIN foods f ON me.food_id = f.id
		LEFT JOIN custom_foods cf ON me.custom_food_id = cf.id
		JOIN daily_logs dl ON me.daily_log_id = dl.id
		WHERE dl.date = ?
		ORDER BY me.meal_type, me.created_at
	`, dateStr)
	if err != nil {
		return nil, fmt.Errorf("query meal entries: %w", err)
	}
	defer rows.Close()

	// Build meals map
	mealsMap := make(map[string]*MealExport)
	for _, mt := range mealOrder {
		mealsMap[mt] = &MealExport{
			MealType:    mt,
			DisplayName: mealDisplayNames[mt],
		}
	}

	for rows.Next() {
		var mealType, foodName, customName, customDesc string
		var calories, protein float64
		var isSaved bool

		if err := rows.Scan(&mealType, &foodName, &calories, &protein, &customName, &customDesc, &isSaved); err != nil {
			return nil, fmt.Errorf("scan meal entry: %w", err)
		}

		meal := mealsMap[mealType]
		if meal == nil {
			continue
		}

		if isSaved {
			meal.HasSavedMeal = true
		}

		if foodName != "" {
			meal.Foods = append(meal.Foods, foodName)
			meal.Calories += calories
			meal.Protein += protein
		}

		if customName != "" {
			var customEntry string
			if customDesc != "" {
				customEntry = fmt.Sprintf("[Custom: %s, %s]", customName, customDesc)
			} else {
				customEntry = fmt.Sprintf("[Custom: %s]", customName)
			}
			meal.CustomFoods = append(meal.CustomFoods, customEntry)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate meal entries: %w", err)
	}

	// Build meals slice in order
	var meals []MealExport
	var totalCal, totalProt float64
	for _, mt := range mealOrder {
		meal := *mealsMap[mt]
		meals = append(meals, meal)
		totalCal += meal.Calories
		totalProt += meal.Protein
	}

	return &DayExport{
		Date:      date,
		DayType:   dayType,
		WeightKg:  weightKg,
		Meals:     meals,
		TotalCal:  totalCal,
		TotalProt: totalProt,
	}, nil
}

// buildEmptyMeals creates empty meal exports for all meal types.
func buildEmptyMeals() []MealExport {
	var meals []MealExport
	for _, mt := range mealOrder {
		meals = append(meals, MealExport{
			MealType:    mt,
			DisplayName: mealDisplayNames[mt],
		})
	}
	return meals
}

// FormatDayMarkdown formats a day's export as markdown.
func FormatDayMarkdown(export *DayExport) string {
	var sb strings.Builder

	// Header: ### Friday, January 31
	sb.WriteString(fmt.Sprintf("### %s\n", export.Date.Format("Monday, January 2")))

	// Day type
	dayTypeDisplay := "Rest Day"
	if export.DayType == "workout" {
		dayTypeDisplay = "Workout Day"
	}
	sb.WriteString(fmt.Sprintf("**Type:** %s  \n", dayTypeDisplay))

	// Weight (if present)
	if export.WeightKg != nil {
		sb.WriteString(fmt.Sprintf("**Weight:** %.1fkg\n", *export.WeightKg))
	}

	sb.WriteString("\n**Meals:**\n")

	// Meals
	hasCustomItems := false
	for _, meal := range export.Meals {
		sb.WriteString(fmt.Sprintf("- **%s:** ", meal.DisplayName))

		if !meal.HasSavedMeal && len(meal.Foods) == 0 && len(meal.CustomFoods) == 0 {
			sb.WriteString("[pending]\n")
			continue
		}

		// Combine foods
		var allItems []string
		allItems = append(allItems, meal.Foods...)
		allItems = append(allItems, meal.CustomFoods...)

		if len(allItems) == 0 {
			sb.WriteString("[pending]\n")
			continue
		}

		sb.WriteString(strings.Join(allItems, " + "))
		sb.WriteString("\n")

		// Macro line
		if len(meal.CustomFoods) > 0 {
			hasCustomItems = true
			sb.WriteString(fmt.Sprintf("  - %.0f kcal | %.0fg protein (+ custom item)\n", meal.Calories, meal.Protein))
		} else if meal.Calories > 0 || meal.Protein > 0 {
			sb.WriteString(fmt.Sprintf("  - %.0f kcal | %.0fg protein\n", meal.Calories, meal.Protein))
		}
	}

	// Total
	sb.WriteString("\n")
	totalCalStr := formatNumber(export.TotalCal)
	sb.WriteString(fmt.Sprintf("**Total:** %s kcal | %.0fg protein\n", totalCalStr, export.TotalProt))

	// Analysis placeholder
	sb.WriteString("\n**Analysis:** [left blank for Claude]\n")

	// Add trailing newline for section separation
	_ = hasCustomItems // used in macro line formatting above

	return sb.String()
}

// formatNumber formats a number with comma thousands separator.
func formatNumber(n float64) string {
	intPart := int(n)
	if intPart >= 1000 {
		return fmt.Sprintf("%d,%03d", intPart/1000, intPart%1000)
	}
	return fmt.Sprintf("%d", intPart)
}

// ExportDayToObsidian exports a day's data to the Obsidian vault.
// It finds and replaces the day's section in the monthly file, or creates it.
func ExportDayToObsidian(db *sql.DB, dateStr string) error {
	vaultPath := os.Getenv("OBSIDIAN_VAULT_PATH")
	if vaultPath == "" {
		// Skip export if vault path not configured
		return nil
	}

	// Get day data
	export, err := GetDayExportData(db, dateStr)
	if err != nil {
		return fmt.Errorf("get day export data: %w", err)
	}

	// Format as markdown
	markdown := FormatDayMarkdown(export)

	// Build file path: 02-Areas/Health & Fitness/Nutrition Logs/YYYY-MM Month.md
	monthFile := export.Date.Format("2006-01 January") + ".md"
	filePath := filepath.Join(vaultPath, "02-Areas", "Health & Fitness", "Nutrition Logs", monthFile)

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create nutrition logs directory: %w", err)
	}

	// Read existing file or start fresh
	var content string
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file with header
			content = fmt.Sprintf("# Nutrition Log - %s\n\n", export.Date.Format("January 2006"))
		} else {
			return fmt.Errorf("read existing file: %w", err)
		}
	} else {
		content = string(existingContent)
	}

	// Find and replace existing day entry, or append
	dayHeader := fmt.Sprintf("### %s", export.Date.Format("Monday, January 2"))
	content = replaceDaySection(content, dayHeader, markdown)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// Mark as synced in database
	_, err = db.Exec(`
		UPDATE daily_logs
		SET synced_to_obsidian = TRUE, updated_at = ?
		WHERE date = ?
	`, time.Now(), dateStr)
	if err != nil {
		return fmt.Errorf("mark synced: %w", err)
	}

	return nil
}

// replaceDaySection replaces an existing day section or appends a new one.
// Days are ordered chronologically (newer days at the end).
func replaceDaySection(content, dayHeader, newSection string) string {
	// Find the start of this day's section
	startIdx := strings.Index(content, dayHeader)
	if startIdx == -1 {
		// Append new section (with blank line separator)
		content = strings.TrimRight(content, "\n")
		if content != "" {
			content += "\n\n"
		}
		content += newSection
		return content
	}

	// Find the end of this day's section (next ### header or EOF)
	afterHeader := startIdx + len(dayHeader)
	endIdx := len(content)

	// Look for next day header after the current one
	nextHeaderIdx := strings.Index(content[afterHeader:], "\n### ")
	if nextHeaderIdx != -1 {
		endIdx = afterHeader + nextHeaderIdx + 1 // +1 to include the newline before ###
	}

	// Replace the section
	return content[:startIdx] + newSection + content[endIdx:]
}
