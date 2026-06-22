package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type convFoodItem struct {
	name    string
	rawKcal float64
	factor  float64 // multiplier: raw -> cooked
}

var convFoods = []convFoodItem{
	{"Pasta secca", 353, 2.5},
	{"Pasta all'uovo secca", 365, 2.4},
	{"Riso secco", 356, 3.0},
	{"Riso integrale secco", 362, 2.8},
	{"Farro secco", 335, 2.5},
	{"Orzo perlato secco", 354, 2.8},
	{"Couscous secco", 376, 2.5},
	{"Quinoa secca", 368, 2.8},
	{"Legumi secchi (ceci, fagioli, lenticchie)", 339, 2.5},
	{"Patate crude", 76, 1.0},
	{"Carne cruda (pollo, manzo, maiale)", 130, 0.7},
	{"Pesce crudo", 80, 0.8},
	{"Verdure crude (spinaci, bietola)", 25, 0.9},
}

func (m Model) converterView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("🔄 Convertitore Crudo ↔ Cotto"))
	b.WriteString("\n\n")

	modeStr := "Crudo → Cotto"
	if m.convMode == 1 {
		modeStr = "Cotto → Crudo"
	}

	b.WriteString(fmt.Sprintf("  Modalità: %s\n\n", highlightStyle.Render(modeStr)))

	// Food list
	b.WriteString(dimmedStyle.Render("  Seleziona l'alimento:\n\n"))
	for i, item := range convFoods {
		prefix := "  "
		style := itemStyle
		if i == m.convFood {
			prefix = "▸ "
			style = selectedItemStyle
		}
		b.WriteString(style.Render(fmt.Sprintf("%s%s", prefix, item.name)))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Input
	selected := convFoods[m.convFood]
	if m.convMode == 0 {
		b.WriteString(fmt.Sprintf("  Grammi crudi di \"%s\":\n", selected.name))
	} else {
		b.WriteString(fmt.Sprintf("  Grammi cotti di \"%s\":\n", selected.name))
	}
	b.WriteString(fmt.Sprintf("  %s█\n\n", m.convInput))

	if m.convError != "" {
		b.WriteString(errorText.Render(fmt.Sprintf("  ⚠ %s", m.convError)))
		b.WriteString("\n\n")
	}

	if m.convResult != "" {
		b.WriteString(greenText.Render(fmt.Sprintf("  ✓ %s", m.convResult)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s naviga  %s cambia modo  %s calcola  %s indietro",
			keyStyle.Render("↑↓"),
			keyStyle.Render("tab"),
			keyStyle.Render("enter"),
			keyStyle.Render("esc"),
		),
	))

	return b.String()
}

func parseFloatSafe(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	var val float64
	fmt.Sscanf(s, "%f", &val)
	return val
}

func (m *Model) handleConverterKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.convFood > 0 {
			m.convFood--
			m.convResult = ""
			m.convError = ""
		}
	case key.Matches(msg, m.keys.Down):
		if m.convFood < len(convFoods)-1 {
			m.convFood++
			m.convResult = ""
			m.convError = ""
		}
	case key.Matches(msg, m.keys.Tab):
		m.convMode = 1 - m.convMode
		m.convResult = ""
		m.convError = ""
	case key.Matches(msg, m.keys.Enter):
		if m.convInput == "" {
			m.convError = "Inserisci un peso"
			return *m, nil
		}
		grams := parseFloatSafe(m.convInput)
		if grams <= 0 {
			m.convError = "Inserisci un numero positivo"
			return *m, nil
		}
		if grams > 5000 {
			m.convError = "Il peso non può superare 5000g"
			return *m, nil
		}

		item := convFoods[m.convFood]
		if m.convMode == 0 {
			// Raw -> Cooked
			cooked := grams * item.factor
			kcal := item.rawKcal * grams / 100
			m.convResult = fmt.Sprintf("%.0fg crudi → %.0fg cotti (%.0f kcal)",
				grams, cooked, kcal)
		} else {
			// Cooked -> Raw
			raw := grams / item.factor
			kcal := item.rawKcal * raw / 100
			m.convResult = fmt.Sprintf("%.0fg cotti ← %.0fg crudi (%.0f kcal)",
				grams, raw, kcal)
		}
		m.convError = ""

	case key.Matches(msg, m.keys.Back):
		m.state = viewHome
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	case key.Matches(msg, m.keys.Backspace):
		if len(m.convInput) > 0 {
			m.convInput = m.convInput[:len(m.convInput)-1]
			m.convResult = ""
		}
	case key.Matches(msg, m.keys.Clear):
		m.convInput = ""
		m.convResult = ""
	default:
		ch := msg.String()
		if ch >= "0" && ch <= "9" || ch == "," || ch == "." {
			if len(m.convInput) < 10 {
				m.convInput += ch
				m.convResult = ""
			}
		}
	}
	return *m, nil
}
