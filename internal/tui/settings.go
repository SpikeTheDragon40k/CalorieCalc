package tui

import (
	"fmt"
	"strings"

	"caloriecalc/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) settingsView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("⚙️  Configurazione"))
	b.WriteString("\n\n")

	settings := []struct {
		name  string
		value string
	}{
		{"Target calorico giornaliero", fmt.Sprintf("%.0f kcal", m.diary.TargetKcal)},
		{"Pasti predefiniti", fmt.Sprintf("%d pasti", len(m.diary.Meals))},
	}

	for i, s := range settings {
		prefix := "  "
		style := itemStyle
		if i == m.settingsCursor {
			prefix = "▸ "
			style = selectedItemStyle
		}
		b.WriteString(style.Render(fmt.Sprintf("%s%-30s %s", prefix, s.name, highlightStyle.Render(s.value))))
		b.WriteString("\n\n")
	}

	// If editing target
	if m.settingsCursor == 0 && m.targetInput != "" {
		b.WriteString(fmt.Sprintf("  Nuovo target (kcal): %s█\n\n", m.targetInput))
		if m.targetError != "" {
			b.WriteString(errorText.Render(fmt.Sprintf("  ⚠ %s", m.targetError)))
			b.WriteString("\n\n")
		}
	}

	b.WriteString(dimmedStyle.Render(fmt.Sprintf("\n  Pasti attuali: %s\n", strings.Join(mealNames(m.diary.Meals), ", "))))
	b.WriteString(dimmedStyle.Render("  Modifica i pasti predefiniti dalla schermata Oggi.\n\n"))

	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s naviga  %s modifica  %s indietro",
			keyStyle.Render("↑↓"),
			keyStyle.Render("enter"),
			keyStyle.Render("esc"),
		),
	))

	return b.String()
}

func mealNames(meals []models.Meal) []string {
	names := make([]string, len(meals))
	for i, m := range meals {
		names[i] = m.Name
	}
	return names
}

func (m *Model) handleSettingsKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		m.targetInput = ""
		m.targetError = ""
	case key.Matches(msg, m.keys.Down):
		if m.settingsCursor < 1 {
			m.settingsCursor++
		}
		m.targetInput = ""
		m.targetError = ""
	case key.Matches(msg, m.keys.Enter):
		if m.settingsCursor == 0 {
			// Toggle target input
			if m.targetInput == "" {
				m.targetInput = fmt.Sprintf("%.0f", m.diary.TargetKcal)
			} else {
				// Validate and save
				val := parseFloatSafe(m.targetInput)
				if val < 100 || val > 10000 {
					m.targetError = "Il target deve essere tra 100 e 10000 kcal"
					return *m, nil
				}
				m.diary.TargetKcal = val
				m.saveAndUpdate()
				m.setStatus(fmt.Sprintf("Target impostato a %.0f kcal", val))
				m.targetInput = ""
				m.targetError = ""
			}
		}
	case key.Matches(msg, m.keys.Back):
		m.state = viewHome
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	case key.Matches(msg, m.keys.Backspace):
		if len(m.targetInput) > 0 {
			m.targetInput = m.targetInput[:len(m.targetInput)-1]
		}
	case key.Matches(msg, m.keys.Clear):
		m.targetInput = ""
	default:
		if m.targetInput != "" {
			ch := msg.String()
			if ch >= "0" && ch <= "9" {
				if len(m.targetInput) < 10 {
					m.targetInput += ch
				}
			}
		}
	}
	return *m, nil
}
