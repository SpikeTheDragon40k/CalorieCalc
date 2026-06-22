package tui

import (
	"fmt"
	"strings"

	"caloriecalc/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) todayView() string {
	var b strings.Builder

	total := m.diary.TodayTotalKcal()
	target := m.diary.TargetKcal

	bar := KcalBar(total, target, 30)
	barColor := KcalColor(total, target)
	totalStr := barColor.Render(fmt.Sprintf("%.0f / %.0f kcal", total, target))

	b.WriteString(titleStyle.Render(fmt.Sprintf("📅 Oggi — %s", models.TodayKey())))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  %s\n\n", bar, totalStr))

	// Meal list
	for i, meal := range m.meals {
		mealKcal := meal.TotalKcal()
		prefix := "  "
		style := itemStyle
		if i == m.cursor {
			prefix = "▸ "
			style = selectedItemStyle
		}

		var foodNames string
		if len(meal.Foods) > 0 {
			var names []string
			for _, f := range meal.Foods {
				names = append(names, f.FoodName)
			}
			foodNames = dimmedStyle.Render(strings.Join(names, ", "))
		} else {
			foodNames = dimmedStyle.Render("(vuoto)")
		}

		mealLine := fmt.Sprintf("%s%-20s %s",
			prefix,
			meal.Name,
			KcalColor(mealKcal, target).Render(fmt.Sprintf("%.0f kcal", mealKcal)),
		)
		b.WriteString(style.Render(mealLine))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("    %s\n", foodNames))
		b.WriteString("\n")
	}

	// Help footer
	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s naviga  %s dettaglio  %s aggiungi pasto  %s elimina  %s indietro",
			keyStyle.Render("↑↓"),
			keyStyle.Render("enter"),
			keyStyle.Render("a"),
			keyStyle.Render("d"),
			keyStyle.Render("esc"),
		),
	))

	return b.String()
}

func (m *Model) handleTodayKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.meals)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.keys.Enter):
		if m.cursor >= 0 && m.cursor < len(m.meals) {
			m.prevState = m.state
			m.state = viewMealDetail
			m.detailCursor = 0
		}
	case key.Matches(msg, m.keys.Add):
		m.prevState = m.state
		m.state = viewAddMeal
		m.mealInput = ""
		m.mealError = ""
	case key.Matches(msg, m.keys.Delete):
		if len(m.meals) > 1 && m.cursor >= 0 && m.cursor < len(m.meals) {
			m.meals = append(m.meals[:m.cursor], m.meals[m.cursor+1:]...)
			if m.cursor >= len(m.meals) {
				m.cursor = len(m.meals) - 1
			}
			m.saveAndUpdate()
			m.setStatus(fmt.Sprintf("Pasto eliminato"))
		}
	case key.Matches(msg, m.keys.Edit):
		if m.cursor >= 0 && m.cursor < len(m.meals) {
			m.prevState = m.state
			m.state = viewAddMeal
			m.mealInput = m.meals[m.cursor].Name
			m.mealError = ""
		}
	case key.Matches(msg, m.keys.Back):
		m.state = viewHome
	case key.Matches(msg, m.keys.Quit):
		m.saveAndUpdate()
		return *m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	}
	return *m, nil
}

// Meal detail view
func (m Model) mealDetailView() string {
	var b strings.Builder

	if m.cursor < 0 || m.cursor >= len(m.meals) {
		return errorText.Render("Nessun pasto selezionato")
	}

	meal := m.meals[m.cursor]
	mealKcal := meal.TotalKcal()

	b.WriteString(titleStyle.Render(fmt.Sprintf("🍽 %s — %.0f kcal", meal.Name, mealKcal)))
	b.WriteString("\n\n")

	if len(meal.Foods) == 0 {
		b.WriteString(dimmedStyle.Render("  Nessun alimento inserito.\n  Premi 'a' per aggiungere.\n"))
	} else {
		for i, food := range meal.Foods {
			prefix := "  "
			style := itemStyle
			if i == m.detailCursor {
				prefix = "▸ "
				style = selectedItemStyle
			}
			foodLine := fmt.Sprintf("%s%s", prefix, food.String())
			b.WriteString(style.Render(foodLine))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Help footer
	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s naviga  %s aggiungi  %s elimina  %s modifica nome  %s indietro",
			keyStyle.Render("↑↓"),
			keyStyle.Render("a"),
			keyStyle.Render("d"),
			keyStyle.Render("e"),
			keyStyle.Render("esc"),
		),
	))

	return b.String()
}

func (m *Model) handleMealDetailKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	meal := &m.meals[m.cursor]

	switch {
	case key.Matches(msg, m.keys.Up):
		if m.detailCursor > 0 {
			m.detailCursor--
		}
	case key.Matches(msg, m.keys.Down):
		if m.detailCursor < len(meal.Foods)-1 {
			m.detailCursor++
		}
	case key.Matches(msg, m.keys.Add):
		m.prevState = m.state
		m.state = viewFoodSearch
		m.searchQuery = ""
		m.searchCursor = 0
		m.searchResults = m.foods
		m.selectedFood = nil
	case key.Matches(msg, m.keys.Delete):
		if len(meal.Foods) > 0 && m.detailCursor >= 0 && m.detailCursor < len(meal.Foods) {
			removed := meal.Foods[m.detailCursor].FoodName
			meal.Foods = append(meal.Foods[:m.detailCursor], meal.Foods[m.detailCursor+1:]...)
			if m.detailCursor >= len(meal.Foods) {
				m.detailCursor = len(meal.Foods) - 1
			}
			m.saveAndUpdate()
			m.setStatus(fmt.Sprintf("Rimosso: %s", removed))
		}
	case key.Matches(msg, m.keys.Edit):
		if m.cursor >= 0 && m.cursor < len(m.meals) {
			m.prevState = m.state
			m.state = viewAddMeal
			m.mealInput = m.meals[m.cursor].Name
			m.mealError = ""
		}
	case key.Matches(msg, m.keys.Back):
		m.state = viewToday
		m.saveAndUpdate()
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	}
	return *m, nil
}

// Add meal view
func (m Model) addMealView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("✏️  Nome pasto"))
	b.WriteString("\n\n")

	label := "Nuovo nome:"
	if m.mealInput == "" {
		label = "Inserisci nome pasto:"
	}

	b.WriteString(fmt.Sprintf("  %s\n\n", label))
	b.WriteString(fmt.Sprintf("  ┌ %s ┐\n", strings.Repeat("─", 30)))
	b.WriteString(fmt.Sprintf("  │ %-30s │\n", m.mealInput+blinkCursor()))
	b.WriteString(fmt.Sprintf("  └ %s ┘\n", strings.Repeat("─", 30)))
	b.WriteString("\n")

	if m.mealError != "" {
		b.WriteString(errorText.Render(fmt.Sprintf("  ⚠ %s", m.mealError)))
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s conferma  %s indietro",
			keyStyle.Render("enter"),
			keyStyle.Render("esc"),
		),
	))

	return b.String()
}

func (m *Model) handleAddMealKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Enter):
		name := strings.TrimSpace(m.mealInput)
		if name == "" {
			m.mealError = "Il nome non può essere vuoto"
			return *m, nil
		}
		if len(name) > 100 {
			m.mealError = "Il nome è troppo lungo (max 100 caratteri)"
			return *m, nil
		}
		// Check if we're editing an existing meal
		if m.prevState == viewMealDetail {
			// Editing meal name
			if m.cursor >= 0 && m.cursor < len(m.meals) {
				m.meals[m.cursor].Name = name
				m.saveAndUpdate()
				m.setStatus(fmt.Sprintf("Pasto rinominato: %s", name))
			}
			m.state = viewMealDetail
		} else {
			// Adding new meal
			newMeal := models.Meal{Name: name}
			// Insert after cursor or at end
			insertPos := len(m.meals)
			if m.cursor >= 0 && m.cursor < len(m.meals) {
				insertPos = m.cursor + 1
			}
			m.meals = append(m.meals[:insertPos], append([]models.Meal{newMeal}, m.meals[insertPos:]...)...)
			m.saveAndUpdate()
			m.setStatus(fmt.Sprintf("Pasto aggiunto: %s", name))
			m.state = viewToday
			if insertPos <= m.cursor {
				m.cursor = insertPos
			}
		}
	case key.Matches(msg, m.keys.Back):
		if m.prevState == viewMealDetail {
			m.state = viewMealDetail
		} else {
			m.state = viewToday
		}
	case key.Matches(msg, m.keys.Backspace):
		if len(m.mealInput) > 0 {
			m.mealInput = m.mealInput[:len(m.mealInput)-1]
		}
	case key.Matches(msg, m.keys.Clear):
		m.mealInput = ""
	default:
		if msg.String() != "" && len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] < 127 {
			if len(m.mealInput) < 100 {
				m.mealInput += msg.String()
			}
		}
	}
	return *m, nil
}

func blinkCursor() string {
	return "█"
}
