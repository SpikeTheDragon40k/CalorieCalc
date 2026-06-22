package tui

import (
	"fmt"
	"strings"
	"time"

	"caloriecalc/internal/store"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) homeView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🥗 CalorieTracker"))
	b.WriteString("\n\n")

	diaryPath := store.DiaryPath()
	if store.DiaryExists() {
		b.WriteString(dimmedStyle.Render(fmt.Sprintf("📁 %s", diaryPath)))
	} else {
		b.WriteString(dimmedStyle.Render("📁 Nuovo diario"))
	}
	b.WriteString("\n\n")

	target := m.diary.TargetKcal
	b.WriteString(fmt.Sprintf("Target giornaliero: %s\n\n",
		highlightStyle.Render(fmt.Sprintf("%.0f kcal", target)),
	))

	menuItems := []struct {
		key     string
		title   string
		desc    string
	}{
		{"1", "📅 Oggi", "Inserisci i pasti per oggi"},
		{"2", "📊 Diario", "Storico kcal (giorno/settimana/mese)"},
		{"3", "🔄 Crudo → Cotto", "Convertitore peso alimenti"},
		{"4", "⚙️  Configurazione", "Target kcal e pasti predefiniti"},
	}

	for _, item := range menuItems {
		b.WriteString(fmt.Sprintf("  [%s] %s\n", keyStyle.Render(item.key), item.title))
		b.WriteString(fmt.Sprintf("      %s\n", dimmedStyle.Render(item.desc)))
		b.WriteString("\n")
	}

	// Quick summary if today has data
	todayMeals := m.diary.GetTodayMeals()
	hasData := false
	for _, mm := range todayMeals {
		if len(mm.Foods) > 0 {
			hasData = true
			break
		}
	}
	if hasData {
		total := m.diary.TodayTotalKcal()
		kcalStyle := KcalColor(total, m.diary.TargetKcal)
		b.WriteString(separator + "\n")
		b.WriteString(fmt.Sprintf("Oggi: %s\n", kcalStyle.Render(fmt.Sprintf("%.0f kcal", total))))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s  %s  %s  %s",
			keyStyle.Render("1-4"),
			dimmedStyle.Render("seleziona"),
			keyStyle.Render("?"),
			dimmedStyle.Render("aiuto"),
		),
	))

	return b.String()
}

func (m *Model) handleHomeKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Number):
		switch msg.String() {
		case "1":
			m.state = viewToday
			m.meals = m.diary.GetTodayMeals()
			m.cursor = 0
		case "2":
			m.state = viewDiaryDay
			m.historyDate = timeNow()
			m.historyCursor = 0
		case "3":
			m.state = viewConverter
			m.convMode = 0
			m.convInput = ""
			m.convResult = ""
			m.convFood = 0
			m.convError = ""
		case "4":
			m.state = viewSettings
			m.settingsCursor = 0
			m.targetInput = ""
			m.targetError = ""
		}
	case key.Matches(msg, m.keys.Quit):
		m.saveAndUpdate()
		return *m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	}
	return *m, nil
}

func timeNow() time.Time {
	return time.Now()
}
