package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Food struct {
	ID                  int     `json:"id" db:"id"`
	Name                string  `json:"name" db:"name"`
	CaloriesPer100g     float64 `json:"kcal" db:"kcal"`
	ProteinPer100g      float64 `json:"protein" db:"protein"`
	StandardServiceSize float64 `json:"serving" db:"serving"`
	Category            string  `json:"category" db:"category"`
}

func (f Food) TotalCalories() (float64, error) {
	isIngredient := f.Category == "Ingredient"
	if f.CaloriesPer100g <= 0 || (!isIngredient && f.StandardServiceSize <= 0) {
		return 0.0, fmt.Errorf("Invalid food item \"%s\"", f.Name)
	}

	if isIngredient {
		return f.CaloriesPer100g, nil
	}

	return (f.CaloriesPer100g / 100) * f.StandardServiceSize, nil
}

func ParseString(line string) (Food, error) {
	if line == "" {
		return Food{}, fmt.Errorf("Input string is null")
	}

	// formatting example: "Chicken Breast: 165 kcal, 31g protein per 100g, Ingredient"
	splits := strings.Split(line, ":")
	name := splits[0]
	subSplits := strings.Split(splits[1], ",")

	trimmedKcal := strings.TrimSpace(subSplits[0])
	kcal, err := strconv.ParseFloat(trimmedKcal[:strings.Index(trimmedKcal, " ")], 64)
	if err != nil {
		return Food{}, err
	}

	protein, err := strconv.ParseFloat(strings.TrimSpace(subSplits[1])[:strings.Index(subSplits[1], "g")-1], 64)
	if err != nil {
		return Food{}, err
	}

	return Food{Name: name, CaloriesPer100g: kcal, ProteinPer100g: protein, Category: strings.TrimSpace(subSplits[2])}, nil
}

func main() {
	fmt.Println("Hello World!")

	foods := []Food{
		{ID: 1, Name: "Oat milk Flat White", CaloriesPer100g: 50, ProteinPer100g: 3,
			StandardServiceSize: 200, Category: "Drink"},
	}
	foods = append(foods, Food{ID: 2, Name: "Chicken", CaloriesPer100g: 80, ProteinPer100g: 30,
		StandardServiceSize: 50, Category: "Ingredient"})
	foods = append(foods, Food{ID: 3, Name: "Invalid Food", CaloriesPer100g: 0, ProteinPer100g: -30,
		StandardServiceSize: -50, Category: "Ingredient"})

	inputString := "Chicken Breast: 165 kcal, 31g protein per 100g, Ingredient"
	stringFood, err := ParseString(inputString)
	if err == nil {
		foods = append(foods, stringFood)
	}

	for _, food := range foods {
		calories, err := food.TotalCalories()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("ID: %d, Name: %s, Calories:%f\n", food.ID, food.Name, calories)
		}
	}
}
