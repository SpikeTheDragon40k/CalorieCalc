package tui

import (
	"fmt"
	"sort"
	"strings"

	"caloriecalc/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) foodSearchView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("🔍 Cerca alimento"))
	b.WriteString("\n\n")

	// Search input
	b.WriteString(fmt.Sprintf("  🔎 %s█\n\n", m.searchQuery))

	// Category filter
	cats := m.foodCategories()
	if m.searchQuery != "" {
		b.WriteString(dimmedStyle.Render(fmt.Sprintf("  Trovati %d alimenti:\n\n", len(m.searchResults))))
	} else {
		// Show categories
		b.WriteString(dimmedStyle.Render("  Sfoglia per categoria o digita per cercare:\n\n"))
		_ = cats
		// Group by category for browsing
		grouped := m.foodsByCategory()
		// Show first few categories
		count := 0
		catKeys := make([]string, 0, len(grouped))
		for k := range grouped {
			catKeys = append(catKeys, k)
		}
		sort.Strings(catKeys)
		for _, cat := range catKeys {
			if count >= 8 {
				b.WriteString(fmt.Sprintf("    ...e altre %d categorie\n", len(catKeys)-count))
				break
			}
			b.WriteString(fmt.Sprintf("    %s (%d)\n", cat, len(grouped[cat])))
			count++
		}
		b.WriteString("\n")
	}

	// Results
	for i, food := range m.searchResults {
		if i >= 15 {
			b.WriteString(dimmedStyle.Render(fmt.Sprintf("    ...e altri %d risultati\n", len(m.searchResults)-15)))
			break
		}
		prefix := "  "
		style := itemStyle
		if i == m.searchCursor {
			prefix = "▸ "
			style = selectedItemStyle
		}
		cat := shortCategory(food.Category)
		foodLine := fmt.Sprintf("%s%-45s %s  %s",
			prefix,
			food.Name,
			greenText.Render(fmt.Sprintf("%.0f kcal/100g", food.KcalPer100)),
			dimmedStyle.Render(cat),
		)
		b.WriteString(style.Render(foodLine))
		b.WriteString("\n")
	}

	// Show selected
	if m.selectedFood != nil {
		b.WriteString(fmt.Sprintf("\n  Selezionato: %s — %s\n",
			highlightStyle.Render(m.selectedFood.Name),
			greenText.Render(fmt.Sprintf("%.0f kcal/100g", m.selectedFood.KcalPer100)),
		))
		b.WriteString(dimmedStyle.Render("  Premi Enter per confermare, o continua a cercare"))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s naviga  %s seleziona  %s indietro  %s",
			keyStyle.Render("↑↓"),
			keyStyle.Render("enter"),
			keyStyle.Render("esc"),
			dimmedStyle.Render("digita per cercare"),
		),
	))

	return b.String()
}

func (m *Model) handleFoodSearchKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.searchCursor > 0 {
			m.searchCursor--
		}
	case key.Matches(msg, m.keys.Down):
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
		}
	case key.Matches(msg, m.keys.Enter):
		if len(m.searchResults) > 0 && m.searchCursor >= 0 && m.searchCursor < len(m.searchResults) {
			m.selectedFood = &m.searchResults[m.searchCursor]
			m.prevState = m.state
			m.state = viewWeightInput
			m.weightInput = ""
			m.weightError = ""
		}
	case key.Matches(msg, m.keys.Back):
		m.state = m.prevState
		m.state = viewMealDetail
	case key.Matches(msg, m.keys.Backspace):
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.filterFoods()
		}
	case key.Matches(msg, m.keys.Clear):
		m.searchQuery = ""
		m.searchResults = m.foods
		m.searchCursor = 0
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	default:
		if msg.String() != "" && len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] < 127 {
			// Only printable ASCII
			if len(m.searchQuery) < 100 {
				m.searchQuery += msg.String()
				m.filterFoods()
			}
		}
	}
	return *m, nil
}

func (m *Model) filterFoods() {
	q := strings.ToLower(m.searchQuery)
	if q == "" {
		m.searchResults = m.foods
		return
	}
	var results []models.Food
	for _, f := range m.foods {
		name := strings.ToLower(f.Name)
		if strings.Contains(name, q) {
			results = append(results, f)
		}
	}
	// If no exact substring matches, try fuzzy (char-by-char matching)
	if len(results) == 0 {
		for _, f := range m.foods {
			if fuzzyMatch(q, strings.ToLower(f.Name)) {
				results = append(results, f)
			}
		}
	}
	m.searchResults = results
	if m.searchCursor >= len(m.searchResults) {
		m.searchCursor = 0
	}
}

func fuzzyMatch(query, target string) bool {
	qi := 0
	for ti := 0; ti < len(target) && qi < len(query); ti++ {
		if query[qi] == target[ti] {
			qi++
		}
	}
	return qi == len(query)
}

func (m Model) foodCategories() []string {
	seen := make(map[string]bool)
	var cats []string
	for _, f := range m.foods {
		if !seen[f.Category] {
			seen[f.Category] = true
			cats = append(cats, f.Category)
		}
	}
	sort.Strings(cats)
	return cats
}

func (m Model) foodsByCategory() map[string][]models.Food {
	grouped := make(map[string][]models.Food)
	for _, f := range m.foods {
		cat := f.Category
		if cat == "" {
			cat = "Varie"
		}
		grouped[cat] = append(grouped[cat], f)
	}
	for k := range grouped {
		sort.Slice(grouped[k], func(i, j int) bool {
			return grouped[k][i].Name < grouped[k][j].Name
		})
	}
	return grouped
}

func shortCategory(cat string) string {
	parts := strings.Split(cat, "/")
	if len(parts) >= 2 {
		return parts[0]
	}
	return cat
}
