package models

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type Food struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	KcalPer100 float64 `json:"kcal_per_100"`
}

type FoodEntry struct {
	FoodName string  `json:"food_name"`
	Grams    float64 `json:"grams"`
	Kcal     float64 `json:"kcal"`
}

func NewFoodEntry(f Food, grams float64) FoodEntry {
	kcal := math.Round(f.KcalPer100*grams/100*10) / 10
	return FoodEntry{FoodName: f.Name, Grams: grams, Kcal: kcal}
}

func (fe FoodEntry) String() string {
	return fmt.Sprintf("%s — %.0fg (%.0f kcal)", fe.FoodName, fe.Grams, fe.Kcal)
}

type Meal struct {
	Name  string      `json:"name"`
	Foods []FoodEntry `json:"foods"`
}

func (m Meal) TotalKcal() float64 {
	var t float64
	for _, f := range m.Foods {
		t += f.Kcal
	}
	return math.Round(t*10) / 10
}

func (m Meal) FoodNames() string {
	var names []string
	for _, f := range m.Foods {
		names = append(names, f.FoodName)
	}
	if len(names) == 0 {
		return "(vuoto)"
	}
	return strings.Join(names, ", ")
}

type Diary struct {
	TargetKcal float64            `json:"target_kcal"`
	Meals      []Meal             `json:"meals"`
	Entries    map[string][]Meal  `json:"entries"`
}

func NewDiary() *Diary {
	return &Diary{
		TargetKcal: 2000,
		Meals:      DefaultMeals(),
		Entries:    make(map[string][]Meal),
	}
}

func DefaultMeals() []Meal {
	return []Meal{
		{Name: "Colazione"},
		{Name: "Spuntino 1"},
		{Name: "Pranzo"},
		{Name: "Spuntino 2"},
		{Name: "Cena"},
		{Name: "Post-Cena"},
	}
}

func TodayKey() string {
	return time.Now().Format("2006-01-02")
}

func (d *Diary) GetTodayMeals() []Meal {
	key := TodayKey()
	if meals, ok := d.Entries[key]; ok {
		return meals
	}
	cp := make([]Meal, len(d.Meals))
	for i, m := range d.Meals {
		cp[i] = Meal{Name: m.Name}
	}
	return cp
}

func (d *Diary) SaveTodayMeals(meals []Meal) {
	d.Entries[TodayKey()] = meals
}

func (d *Diary) TodayTotalKcal() float64 {
	var t float64
	for _, m := range d.GetTodayMeals() {
		t += m.TotalKcal()
	}
	return math.Round(t*10) / 10
}

func (d *Diary) GetDates() []string {
	var dates []string
	for k := range d.Entries {
		dates = append(dates, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	return dates
}

func (d *Diary) GetMealsForDate(date string) []Meal {
	if meals, ok := d.Entries[date]; ok {
		return meals
	}
	return nil
}

func (d *Diary) GetWeekDates(ref time.Time) []string {
	weekday := ref.Weekday()
	monday := ref.AddDate(0, 0, -int(weekday)+1)
	var dates []string
	for i := 0; i < 7; i++ {
		dates = append(dates, monday.AddDate(0, 0, i).Format("2006-01-02"))
	}
	return dates
}

func (d *Diary) WeekTotalKcal(ref time.Time) float64 {
	var t float64
	for _, date := range d.GetWeekDates(ref) {
		for _, m := range d.GetMealsForDate(date) {
			t += m.TotalKcal()
		}
	}
	return math.Round(t*10) / 10
}

func (d *Diary) MonthTotalKcal(ref time.Time) float64 {
	year, month, _ := ref.Date()
	var t float64
	for date := range d.Entries {
		dt, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}
		if dt.Year() == year && dt.Month() == month {
			for _, m := range d.Entries[date] {
				t += m.TotalKcal()
			}
		}
	}
	return math.Round(t*10) / 10
}

func (d *Diary) DayTotalKcalForDate(date string) float64 {
	var t float64
	for _, m := range d.GetMealsForDate(date) {
		t += m.TotalKcal()
	}
	return math.Round(t*10) / 10
}

func (d *Diary) Serialize() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

func Deserialize(data []byte) (*Diary, error) {
	var d Diary
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	if d.Meals == nil {
		d.Meals = DefaultMeals()
	}
	if d.Entries == nil {
		d.Entries = make(map[string][]Meal)
	}
	if d.TargetKcal <= 0 {
		d.TargetKcal = 2000
	}
	return &d, nil
}
