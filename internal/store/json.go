package store

import (
	"caloriecalc/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var diaryFile = "diario.json"
const maxSize = 10 * 1024 * 1024

func GetDiaryFile() string {
	return diaryFile
}

func SetDiaryFile(name string) {
	diaryFile = name
}

func SaveDiary(d *models.Diary) error {
	data, err := d.Serialize()
	if err != nil {
		return fmt.Errorf("serialize diary: %w", err)
	}
	tmpPath := diaryFile + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmpPath, diaryFile); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func LoadDiary() (*models.Diary, error) {
	fi, err := os.Lstat(diaryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return models.NewDiary(), nil
		}
		return nil, fmt.Errorf("stat diary file: %w", err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("refusing to follow symlink: %s", diaryFile)
	}
	f, err := os.OpenFile(diaryFile, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("open diary file: %w", err)
	}
	defer f.Close()

	r := io.LimitReader(f, maxSize)
	var d models.Diary
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&d); err != nil {
		// If corrupted or old format, reset
		return models.NewDiary(), nil
	}
	if d.Meals == nil {
		d.Meals = models.DefaultMeals()
	}
	if d.Entries == nil {
		d.Entries = make(map[string][]models.Meal)
	}
	if d.TargetKcal <= 0 {
		d.TargetKcal = 2000
	}
	// Clean old keys that have all empty meals
	for key, meals := range d.Entries {
		allEmpty := true
		for _, m := range meals {
			if len(m.Foods) > 0 {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			delete(d.Entries, key)
		}
	}
	return &d, nil
}

func DiaryExists() bool {
	_, err := os.Lstat(diaryFile)
	return err == nil
}

func DiaryPath() string {
	p, _ := filepath.Abs(diaryFile)
	return p
}

func ValidateGrams(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	g, err := parseFloatStrict(s)
	if err != nil {
		return 0, fmt.Errorf("inserisci un numero valido (es. 150)")
	}
	if g <= 0 {
		return 0, fmt.Errorf("il peso deve essere maggiore di 0")
	}
	if g > 5000 {
		return 0, fmt.Errorf("il peso non può superare 5000g")
	}
	return g, nil
}

func parseFloatStrict(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
