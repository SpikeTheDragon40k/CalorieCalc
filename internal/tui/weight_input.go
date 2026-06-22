package tui

import (
	"fmt"
	"strings"

	"caloriecalc/internal/models"
	"caloriecalc/internal/store"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) weightInputView() string {
	var b strings.Builder
	if m.selectedFood == nil {
		return errorText.Render("Nessun alimento selezionato")
	}

	b.WriteString(titleStyle.Render("⚖️  Inserisci peso"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  Alimento: %s\n", highlightStyle.Render(m.selectedFood.Name)))
	b.WriteString(fmt.Sprintf("  Calorie:  %s\n\n", greenText.Render(fmt.Sprintf("%.0f kcal/100g", m.selectedFood.KcalPer100))))

	b.WriteString("  Grammi: ")
	b.WriteString(fmt.Sprintf("%s█\n\n", m.weightInput))

	if m.weightError != "" {
		b.WriteString(errorText.Render(fmt.Sprintf("  ⚠ %s", m.weightError)))
		b.WriteString("\n\n")
	}

	// Preview
	if g, err := store.ValidateGrams(m.weightInput); err == nil {
		kcal := m.selectedFood.KcalPer100 * g / 100
		b.WriteString(greenText.Render(fmt.Sprintf("  ≈ %.0fg → %.0f kcal", g, kcal)))
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

func (m *Model) handleWeightInputKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Enter):
		if m.selectedFood == nil {
			m.weightError = "Nessun alimento selezionato"
			return *m, nil
		}
		grams, err := store.ValidateGrams(m.weightInput)
		if err != nil {
			m.weightError = err.Error()
			return *m, nil
		}
		entry := models.NewFoodEntry(*m.selectedFood, grams)
		if m.cursor >= 0 && m.cursor < len(m.meals) {
			m.meals[m.cursor].Foods = append(m.meals[m.cursor].Foods, entry)
			m.saveAndUpdate()
			m.setStatus(fmt.Sprintf("Aggiunto: %s (%.0fg, %.0f kcal)",
				m.selectedFood.Name, grams, entry.Kcal))
		}
		m.selectedFood = nil
		m.state = viewMealDetail
		m.detailCursor = len(m.meals[m.cursor].Foods) - 1
		if m.detailCursor < 0 {
			m.detailCursor = 0
		}

	case key.Matches(msg, m.keys.Back):
		m.state = viewFoodSearch
	case key.Matches(msg, m.keys.Backspace):
		if len(m.weightInput) > 0 {
			m.weightInput = m.weightInput[:len(m.weightInput)-1]
		}
	case key.Matches(msg, m.keys.Clear):
		m.weightInput = ""
	default:
		// Accept digits, comma, dot, and backspace
		ch := msg.String()
		if ch >= "0" && ch <= "9" || ch == "," || ch == "." {
			if len(m.weightInput) < 10 {
				m.weightInput += ch
			}
		}
	}
	return *m, nil
}
