package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) diaryView() string {
	var b strings.Builder

	viewName := "Giorno"
	if m.state == viewDiaryWeek {
		viewName = "Settimana"
	} else if m.state == viewDiaryMonth {
		viewName = "Mese"
	}

	b.WriteString(titleStyle.Render(fmt.Sprintf("📊 Diario — Vista %s", viewName)))
	b.WriteString("\n")

	// Nav bar
	dateStr := m.historyDate.Format("2006-01-02")
	b.WriteString(fmt.Sprintf("  %s  %s  %s\n\n",
		keyStyle.Render("←"),
		highlightStyle.Render(dateStr),
		keyStyle.Render("→"),
	))

	// Mode switcher
	currentMode := " [Giorno] "
	if m.state == viewDiaryWeek {
		currentMode = "  Giorno   [Settimana]  Mese  "
	} else if m.state == viewDiaryMonth {
		currentMode = "  Giorno   Settimana  [Mese]  "
	} else {
		currentMode = " [Giorno]  Settimana  Mese  "
	}
	b.WriteString(dimmedStyle.Render(fmt.Sprintf("  %s\n\n", currentMode)))

	switch m.state {
	case viewDiaryDay:
		m.diaryDayView(&b)
	case viewDiaryWeek:
		m.diaryWeekView(&b)
	case viewDiaryMonth:
		m.diaryMonthView(&b)
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render(
		fmt.Sprintf("%s naviga  %s cambia vista  %s indietro",
			keyStyle.Render("← →"),
			keyStyle.Render("tab"),
			keyStyle.Render("esc"),
		),
	))

	return b.String()
}

func (m Model) diaryDayView(b *strings.Builder) {
	dateStr := m.historyDate.Format("2006-01-02")
	meals := m.diary.GetMealsForDate(dateStr)

	if meals == nil {
		b.WriteString(dimmedStyle.Render("  Nessun dato per questa data.\n"))
		// Show target
		b.WriteString(fmt.Sprintf("\n  Target: %s\n", highlightStyle.Render(fmt.Sprintf("%.0f kcal", m.diary.TargetKcal))))
		return
	}

	total := m.diary.DayTotalKcalForDate(dateStr)
	target := m.diary.TargetKcal

	bar := KcalBar(total, target, 25)
	barColor := KcalColor(total, target)
	totalStr := barColor.Render(fmt.Sprintf("%.0f / %.0f kcal", total, target))

	b.WriteString(fmt.Sprintf("  %s  %s\n\n", bar, totalStr))

	for i, meal := range meals {
		mealKcal := meal.TotalKcal()
		prefix := "  "
		style := itemStyle
		if i == m.historyCursor {
			prefix = "▸ "
			style = selectedItemStyle
		}
		mealLine := fmt.Sprintf("%s%-20s %s",
			prefix,
			meal.Name,
			KcalColor(mealKcal, target).Render(fmt.Sprintf("%.0f kcal", mealKcal)),
		)
		b.WriteString(style.Render(mealLine))
		b.WriteString("\n")

		if i == m.historyCursor && len(meal.Foods) > 0 {
			for _, food := range meal.Foods {
				b.WriteString(fmt.Sprintf("    • %s\n", food.String()))
			}
		}
		b.WriteString("\n")
	}
}

func (m Model) diaryWeekView(b *strings.Builder) {
	dates := m.diary.GetWeekDates(m.historyDate)
	target := m.diary.TargetKcal
	weekTotal := m.diary.WeekTotalKcal(m.historyDate)

	b.WriteString(fmt.Sprintf("  Totale settimana: %s\n\n",
		KcalColor(weekTotal, target*7).Render(fmt.Sprintf("%.0f / %.0f kcal", weekTotal, target*7)),
	))

	for _, date := range dates {
		total := m.diary.DayTotalKcalForDate(date)
		bar := KcalBar(total, target, 20)
		displayDate := date
		if dt, err := time.Parse("2006-01-02", date); err == nil {
			displayDate = dt.Format("lun 02/01")
			// Italian day abbreviations
			switch dt.Weekday() {
			case time.Monday:
				displayDate = "Lun"
			case time.Tuesday:
				displayDate = "Mar"
			case time.Wednesday:
				displayDate = "Mer"
			case time.Thursday:
				displayDate = "Gio"
			case time.Friday:
				displayDate = "Ven"
			case time.Saturday:
				displayDate = "Sab"
			case time.Sunday:
				displayDate = "Dom"
			}
			displayDate += fmt.Sprintf(" %02d/%02d", dt.Day(), dt.Month())
		}

		prefix := "  "
		style := dimmedStyle
		if date == m.historyDate.Format("2006-01-02") {
			prefix = "▸ "
			style = itemStyle
		}

		if total > 0 {
			line := fmt.Sprintf("%s%s │%s│ %s",
				prefix,
				displayDate,
				bar,
				KcalColor(total, target).Render(fmt.Sprintf("%.0f kcal", total)),
			)
			b.WriteString(style.Render(line))
		} else {
			line := fmt.Sprintf("%s%s │%s│ %s",
				prefix,
				displayDate,
				strings.Repeat("░", 20),
				dimmedStyle.Render("—"),
			)
			b.WriteString(style.Render(line))
		}
		b.WriteString("\n")
	}
}

func (m Model) diaryMonthView(b *strings.Builder) {
	year, month, _ := m.historyDate.Date()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)
	target := m.diary.TargetKcal
	monthTotal := m.diary.MonthTotalKcal(m.historyDate)
	daysInMonth := lastDay.Day()

	b.WriteString(fmt.Sprintf("  %s %d\n", highlightStyle.Render(month.String()[:3]), year))
	b.WriteString(fmt.Sprintf("  Totale mese: %s\n\n",
		KcalColor(monthTotal, target*float64(daysInMonth)).Render(fmt.Sprintf("%.0f kcal", monthTotal)),
	))

	// Simple list view for month
	b.WriteString(fmt.Sprintf("  %-12s %s\n", "Giorno", dimmedStyle.Render("Kcal")))
	b.WriteString(fmt.Sprintf("  %s\n", strings.Repeat("─", 35)))

	currentDate := m.historyDate.Format("2006-01-02")
	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		dateStr := date.Format("2006-01-02")
		total := m.diary.DayTotalKcalForDate(dateStr)

		prefix := "  "
		if dateStr == currentDate {
			prefix = "▸ "
		}

		dayLabel := fmt.Sprintf("%02d/%02d", day, int(month))
		var totalStr string
		if total > 0 {
			bar := KcalBar(total, target, 10)
			totalStr = fmt.Sprintf("%s %s", bar, KcalColor(total, target).Render(fmt.Sprintf("%.0f", total)))
		} else {
			totalStr = dimmedStyle.Render("—")
		}

		b.WriteString(fmt.Sprintf("%s%-8s %s\n", prefix, dayLabel, totalStr))
	}
}

func (m *Model) handleDiaryKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left):
		switch m.state {
		case viewDiaryDay:
			m.historyDate = m.historyDate.AddDate(0, 0, -1)
		case viewDiaryWeek:
			m.historyDate = m.historyDate.AddDate(0, 0, -7)
		case viewDiaryMonth:
			m.historyDate = m.historyDate.AddDate(0, -1, 0)
		}
		m.historyCursor = 0
	case key.Matches(msg, m.keys.Right):
		switch m.state {
		case viewDiaryDay:
			m.historyDate = m.historyDate.AddDate(0, 0, 1)
		case viewDiaryWeek:
			m.historyDate = m.historyDate.AddDate(0, 0, 7)
		case viewDiaryMonth:
			m.historyDate = m.historyDate.AddDate(0, 1, 0)
		}
		m.historyCursor = 0
	case key.Matches(msg, m.keys.Tab):
		switch m.state {
		case viewDiaryDay:
			m.state = viewDiaryWeek
		case viewDiaryWeek:
			m.state = viewDiaryMonth
		case viewDiaryMonth:
			m.state = viewDiaryDay
		}
		m.historyCursor = 0
	case key.Matches(msg, m.keys.Back):
		m.state = viewHome
	case key.Matches(msg, m.keys.Help):
		m.prevState = m.state
		m.state = viewHelp
	}
	return *m, nil
}
