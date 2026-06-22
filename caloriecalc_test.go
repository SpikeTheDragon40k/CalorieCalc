package main

import (
	"caloriecalc/internal/models"
	"caloriecalc/internal/parser"
	"caloriecalc/internal/store"
	"os"
	"strings"
	"testing"
)

func TestLoadFoods(t *testing.T) {
	foods, err := parser.LoadFoods()
	if err != nil {
		t.Fatalf("Failed to load foods: %v", err)
	}
	if len(foods) == 0 {
		t.Fatal("No foods loaded")
	}
	if len(foods) < 350 {
		t.Fatalf("Expected at least 350 foods, got %d", len(foods))
	}
	// Check known foods
	found := false
	for _, f := range foods {
		if strings.Contains(f.Name, "Yogurt") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find Yogurt in foods")
	}
}

func TestDiarySerialization(t *testing.T) {
	d := models.NewDiary()
	d.TargetKcal = 2000
	d.Meals = models.DefaultMeals()

	// Add food to Colazione
	meal := &d.Meals[0]
	meal.Foods = append(meal.Foods, models.FoodEntry{
		FoodName: "Yogurt, naturale",
		Grams:    150,
		Kcal:     99,
	})

	d.SaveTodayMeals(d.Meals)

	data, err := d.Serialize()
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	if !strings.Contains(string(data), "Yogurt") {
		t.Error("Expected Yogurt in serialized data")
	}
	if !strings.Contains(string(data), models.TodayKey()) {
		t.Error("Expected today's date in serialized data")
	}

	d2, err := models.Deserialize(data)
	if err != nil {
		t.Fatalf("Deserialize: %v", err)
	}
	if d2.TargetKcal != 2000 {
		t.Errorf("TargetKcal = %f, want 2000", d2.TargetKcal)
	}
	meals := d2.GetTodayMeals()
	if len(meals) == 0 {
		t.Fatal("No meals for today")
	}
	if len(meals[0].Foods) == 0 {
		t.Error("Expected foods in first meal")
	}
	if meals[0].Foods[0].FoodName != "Yogurt, naturale" {
		t.Errorf("FoodName = %q, want Yogurt", meals[0].Foods[0].FoodName)
	}
}

func TestValidateGrams(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		err   bool
	}{
		{"100", 100, false},
		{"150.5", 150.5, false},
		{"200,5", 200.5, false},
		{"0", 0, true},
		{"-10", 0, true},
		{"abc", 0, true},
		{"5001", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := store.ValidateGrams(tt.input)
		if tt.err && err == nil {
			t.Errorf("ValidateGrams(%q) expected error", tt.input)
		}
		if !tt.err && err != nil {
			t.Errorf("ValidateGrams(%q) unexpected error: %v", tt.input, err)
		}
		if !tt.err && got != tt.want {
			t.Errorf("ValidateGrams(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}
}

func TestNewFoodEntry(t *testing.T) {
	f := models.Food{Name: "Mela", Category: "Frutta", KcalPer100: 55}
	entry := models.NewFoodEntry(f, 200)
	if entry.FoodName != "Mela" {
		t.Errorf("FoodName = %q", entry.FoodName)
	}
	if entry.Grams != 200 {
		t.Errorf("Grams = %f", entry.Grams)
	}
	if entry.Kcal != 110.0 {
		t.Errorf("Kcal = %f, want 110", entry.Kcal)
	}
}

func TestMealTotal(t *testing.T) {
	meal := models.Meal{Name: "Pranzo"}
	meal.Foods = append(meal.Foods, models.FoodEntry{FoodName: "Pasta", Grams: 100, Kcal: 353})
	meal.Foods = append(meal.Foods, models.FoodEntry{FoodName: "Pomodoro", Grams: 200, Kcal: 42})
	total := meal.TotalKcal()
	if total != 395.0 {
		t.Errorf("Total = %f, want 395", total)
	}
}

func TestDiaryTodayTotal(t *testing.T) {
	d := models.NewDiary()
	d.TargetKcal = 2000

	meals := []models.Meal{
		{Name: "Colazione", Foods: []models.FoodEntry{{FoodName: "Yogurt", Grams: 150, Kcal: 99}}},
		{Name: "Pranzo", Foods: []models.FoodEntry{{FoodName: "Pasta", Grams: 100, Kcal: 353}}},
	}
	d.SaveTodayMeals(meals)

	total := d.TodayTotalKcal()
	if total != 452.0 {
		t.Errorf("TodayTotal = %f, want 452", total)
	}
}

func TestFoodEntryString(t *testing.T) {
	entry := models.FoodEntry{FoodName: "Pasta", Grams: 100, Kcal: 353}
	s := entry.String()
	if !strings.Contains(s, "Pasta") || !strings.Contains(s, "100") {
		t.Errorf("String = %q", s)
	}
}

func TestDiaryDefaultMeals(t *testing.T) {
	d := models.NewDiary()
	if len(d.Meals) != 6 {
		t.Errorf("Expected 6 default meals, got %d", len(d.Meals))
	}
	if d.Meals[0].Name != "Colazione" {
		t.Errorf("First meal = %q", d.Meals[0].Name)
	}
}

func TestStoreSaveLoad(t *testing.T) {
	tmpFile := "test_diario.json"
	defer os.Remove(tmpFile)
	defer os.Remove(tmpFile + ".tmp")

	origFile := store.GetDiaryFile()
	store.SetDiaryFile(tmpFile)
	defer store.SetDiaryFile(origFile)

	d := models.NewDiary()
	d.TargetKcal = 1800
	meals := []models.Meal{
		{Name: "Colazione", Foods: []models.FoodEntry{{FoodName: "Caffè", Grams: 100, Kcal: 2}}},
	}
	d.SaveTodayMeals(meals)

	if err := store.SaveDiary(d); err != nil {
		t.Fatalf("Save: %v", err)
	}

	d2, err := store.LoadDiary()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if d2.TargetKcal != 1800 {
		t.Errorf("TargetKcal = %f", d2.TargetKcal)
	}
}

func TestCaloriesCalculation(t *testing.T) {
	f := models.Food{Name: "Test", KcalPer100: 200}
	entry := models.NewFoodEntry(f, 50) // 50g of 200kcal/100g = 100kcal
	if entry.Kcal != 100.0 {
		t.Errorf("Expected 100 kcal, got %f", entry.Kcal)
	}
}
