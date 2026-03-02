package export

import (
	"strings"
	"testing"
	"time"
)

func TestFormatDayMarkdown_FullDay(t *testing.T) {
	weight := 83.1
	export := &DayExport{
		Date:     time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
		DayType:  "rest",
		WeightKg: &weight,
		Meals: []MealExport{
			{
				MealType:     "breakfast",
				DisplayName:  "Breakfast",
				Foods:        []string{"Coffee", "Pumpkin bread", "Cream cheese", "Smoked salmon"},
				Calories:     404,
				Protein:      24,
				HasSavedMeal: true,
			},
			{
				MealType:     "lunch",
				DisplayName:  "Lunch",
				Foods:        []string{"Chicken dürüm (large)"},
				Calories:     690,
				Protein:      50,
				HasSavedMeal: true,
			},
			{
				MealType:     "snacks",
				DisplayName:  "Snacks",
				Foods:        []string{"Club Mate"},
				CustomFoods:  []string{"[Custom: handful of nuts, medium portion]"},
				Calories:     100,
				Protein:      0,
				HasSavedMeal: true,
			},
			{
				MealType:    "dinner",
				DisplayName: "Dinner",
			},
			{
				MealType:     "shake",
				DisplayName:  "Shake",
				Foods:        []string{"Double scoop"},
				Calories:     222,
				Protein:      48,
				HasSavedMeal: true,
			},
		},
		TotalCal:  1416,
		TotalProt: 122,
	}

	result := FormatDayMarkdown(export)

	// Check header
	if !strings.Contains(result, "### Saturday, January 31") {
		t.Errorf("Expected header with day name, got:\n%s", result)
	}

	// Check day type
	if !strings.Contains(result, "**Type:** Rest Day") {
		t.Errorf("Expected Rest Day type, got:\n%s", result)
	}

	// Check weight
	if !strings.Contains(result, "**Weight:** 83.1kg") {
		t.Errorf("Expected weight, got:\n%s", result)
	}

	// Check breakfast format
	if !strings.Contains(result, "- **Breakfast:** Coffee + Pumpkin bread + Cream cheese + Smoked salmon") {
		t.Errorf("Expected breakfast with + separator, got:\n%s", result)
	}

	// Check custom food format
	if !strings.Contains(result, "[Custom: handful of nuts, medium portion]") {
		t.Errorf("Expected custom food format, got:\n%s", result)
	}

	// Check custom item note
	if !strings.Contains(result, "(+ custom item)") {
		t.Errorf("Expected custom item note, got:\n%s", result)
	}

	// Check pending dinner
	if !strings.Contains(result, "- **Dinner:** [pending]") {
		t.Errorf("Expected pending dinner, got:\n%s", result)
	}

	// Check total with comma formatting
	if !strings.Contains(result, "**Total:** 1,416 kcal | 122g protein") {
		t.Errorf("Expected formatted total, got:\n%s", result)
	}

	// Check analysis placeholder
	if !strings.Contains(result, "**Analysis:** [left blank for Claude]") {
		t.Errorf("Expected analysis placeholder, got:\n%s", result)
	}
}

func TestFormatDayMarkdown_WorkoutDay(t *testing.T) {
	export := &DayExport{
		Date:    time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		DayType: "workout",
		Meals:   buildEmptyMeals(),
	}

	result := FormatDayMarkdown(export)

	if !strings.Contains(result, "**Type:** Workout Day") {
		t.Errorf("Expected Workout Day type, got:\n%s", result)
	}

	// No weight should be present
	if strings.Contains(result, "**Weight:**") {
		t.Errorf("Should not have weight when nil, got:\n%s", result)
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{500, "500"},
		{999, "999"},
		{1000, "1,000"},
		{1416, "1,416"},
		{2500, "2,500"},
	}

	for _, tc := range tests {
		result := formatNumber(tc.input)
		if result != tc.expected {
			t.Errorf("formatNumber(%v) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestReplaceDaySection_NewSection(t *testing.T) {
	content := "# Nutrition Log - January 2026\n\n"
	dayHeader := "### Friday, January 31"
	newSection := "### Friday, January 31\n**Type:** Rest Day  \n\n"

	result := replaceDaySection(content, dayHeader, newSection)

	if !strings.Contains(result, newSection) {
		t.Errorf("Expected new section to be appended, got:\n%s", result)
	}
}

func TestReplaceDaySection_ReplaceExisting(t *testing.T) {
	content := `# Nutrition Log - January 2026

### Thursday, January 30
**Type:** Rest Day

### Friday, January 31
**Type:** Rest Day
**Weight:** 83.0kg

**Meals:**
- **Breakfast:** [pending]
`
	dayHeader := "### Friday, January 31"
	newSection := "### Friday, January 31\n**Type:** Workout Day  \n**Weight:** 83.1kg\n"

	result := replaceDaySection(content, dayHeader, newSection)

	// Should contain new section
	if !strings.Contains(result, "**Type:** Workout Day") {
		t.Errorf("Expected updated day type, got:\n%s", result)
	}

	// Should not contain old weight
	if strings.Contains(result, "83.0kg") {
		t.Errorf("Should have replaced old weight, got:\n%s", result)
	}

	// Should still contain Thursday
	if !strings.Contains(result, "### Thursday, January 30") {
		t.Errorf("Should preserve other days, got:\n%s", result)
	}
}
