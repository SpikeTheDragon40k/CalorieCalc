package parser

import (
	"caloriecalc/internal/models"
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//go:embed alimenti.csv
var csvData embed.FS

func LoadFoods() ([]models.Food, error) {
	file, err := csvData.Open("alimenti.csv")
	if err != nil {
		return nil, fmt.Errorf("cannot open embedded CSV: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ';'
	r.LazyQuotes = true

	var foods []models.Food
	line := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("CSV parse error at line %d: %w", line+1, err)
		}
		line++
		if line == 1 {
			continue
		}
		if len(record) < 3 {
			continue
		}
		name := strings.TrimSpace(record[0])
		category := strings.TrimSpace(record[1])
		// Handle multi-category entries: take the first category
		if idx := strings.Index(category, ";"); idx >= 0 {
			category = strings.TrimSpace(category[:idx])
		}
		kcals := strings.TrimSpace(record[2])
		// Replace comma with dot for decimal parsing
		kcals = strings.ReplaceAll(kcals, ",", ".")
		kcal, err := strconv.ParseFloat(kcals, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid kcal value %q at line %d: %w", record[2], line, err)
		}
		foods = append(foods, models.Food{
			Name:       name,
			Category:   category,
			KcalPer100: kcal,
		})
	}
	return foods, nil
}
