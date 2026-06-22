package tui

import (
	"caloriecalc/internal/models"
	"caloriecalc/internal/parser"
	"caloriecalc/internal/store"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState int

const (
	viewHome viewState = iota
	viewToday
	viewDiaryDay
	viewDiaryWeek
	viewDiaryMonth
	viewConverter
	viewSettings
	viewFoodSearch
	viewWeightInput
	viewMealDetail
	viewAddMeal
	viewHelp
	viewError
)

type screen int

const (
	screenMain screen = iota
	screenDetail
)

type Model struct {
	state    viewState
	prevState viewState
	diary    *models.Diary
	foods    []models.Food
	width    int
	height   int
	help     help.Model
	keys     keyMap

	// Navigation
	cursor     int
	detailCursor int
	foodCursor int

	// Today view
	meals []models.Meal

	// History
	historyDate time.Time
	historyCursor int

	// Food search
	searchQuery    string
	searchCursor   int
	searchResults  []models.Food
	selectedFood   *models.Food

	// Weight input
	weightInput    string
	weightError    string

	// Converter
	convMode       int // 0: raw->cooked, 1: cooked->raw
	convInput      string
	convResult     string
	convFood       int
	convError      string

	// Settings
	settingsCursor int
	targetInput    string
	targetError    string
	mealInput      string
	mealError      string

	// Messages
	statusMsg      string
	statusTimer    int

	err error
}

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Enter     key.Binding
	Back      key.Binding
	Quit      key.Binding
	Help      key.Binding
	Add       key.Binding
	Delete    key.Binding
	Tab       key.Binding
	Edit      key.Binding
	Number    key.Binding
	Backspace key.Binding
	Clear     key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "su")),
		Down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "giù")),
		Left:      key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "indietro")),
		Right:     key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "avanti")),
		Enter:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "seleziona")),
		Back:      key.NewBinding(key.WithKeys("esc", "backspace"), key.WithHelp("esc", "indietro")),
		Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "esci")),
		Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Add:       key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "aggiungi")),
		Delete:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "elimina")),
		Tab:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "pannello")),
		Edit:      key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "modifica")),
		Number:    key.NewBinding(key.WithKeys("1", "2", "3", "4"), key.WithHelp("1-4", "menu")),
		Backspace: key.NewBinding(key.WithKeys("backspace"), key.WithHelp("backspace", "cancella")),
		Clear:     key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("ctrl+u", "pulisci")),
	}
}

func NewModel() Model {
	diary, err := store.LoadDiary()
	if err != nil {
		diary = models.NewDiary()
	}

	foods, err := parser.LoadFoods()
	if err != nil {
		foods = []models.Food{}
	}

	return Model{
		state:          viewHome,
		diary:          diary,
		foods:          foods,
		help:           help.New(),
		keys:           defaultKeyMap(),
		meals:          diary.GetTodayMeals(),
		historyDate:    time.Now(),
		convMode:       0,
		settingsCursor: 0,
		statusTimer:    0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) saveAndUpdate() {
	m.diary.SaveTodayMeals(m.meals)
	if err := store.SaveDiary(m.diary); err != nil {
		m.err = err
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		if m.statusTimer > 0 {
			m.statusTimer--
			if m.statusTimer <= 0 {
				m.statusMsg = ""
			}
		}
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

func (m Model) handleKeyMsg(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case viewHome:
			return m.handleHomeKey(msg)
		case viewToday:
			return m.handleTodayKey(msg)
		case viewMealDetail:
			return m.handleMealDetailKey(msg)
		case viewFoodSearch:
			return m.handleFoodSearchKey(msg)
		case viewWeightInput:
			return m.handleWeightInputKey(msg)
		case viewDiaryDay, viewDiaryWeek, viewDiaryMonth:
			return m.handleDiaryKey(msg)
		case viewConverter:
			return m.handleConverterKey(msg)
		case viewSettings:
			return m.handleSettingsKey(msg)
		case viewAddMeal:
			return m.handleAddMealKey(msg)
		case viewHelp:
			if key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Help) {
				m.state = m.prevState
			}
			return m, nil
		default:
			if key.Matches(msg, m.keys.Quit) {
				m.saveAndUpdate()
				return m, tea.Quit
			}
			if key.Matches(msg, m.keys.Back) {
				m.state = viewHome
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) setStatus(msg string) {
	m.statusMsg = msg
	m.statusTimer = 3
}

// Helper to update food list for a given meal in the today view
func (m *Model) updateFoodInMeal(mealIdx int, entry models.FoodEntry, add bool) {
	if mealIdx < 0 || mealIdx >= len(m.meals) {
		return
	}
	if add {
		m.meals[mealIdx].Foods = append(m.meals[mealIdx].Foods, entry)
	} else {
		// remove by detailCursor
		if m.detailCursor >= 0 && m.detailCursor < len(m.meals[mealIdx].Foods) {
			m.meals[mealIdx].Foods = append(
				m.meals[mealIdx].Foods[:m.detailCursor],
				m.meals[mealIdx].Foods[m.detailCursor+1:]...,
			)
		}
	}
	m.saveAndUpdate()
}

func (m Model) View() string {
	if m.err != nil {
		return errorText.Render(fmt.Sprintf("Errore: %v", m.err))
	}

	var mainContent string
	switch m.state {
	case viewHome:
		mainContent = m.homeView()
	case viewToday:
		mainContent = m.todayView()
	case viewMealDetail:
		mainContent = m.mealDetailView()
	case viewFoodSearch:
		mainContent = m.foodSearchView()
	case viewWeightInput:
		mainContent = m.weightInputView()
	case viewDiaryDay, viewDiaryWeek, viewDiaryMonth:
		mainContent = m.diaryView()
	case viewConverter:
		mainContent = m.converterView()
	case viewSettings:
		mainContent = m.settingsView()
	case viewAddMeal:
		mainContent = m.addMealView()
	case viewHelp:
		mainContent = m.helpView()
	default:
		mainContent = m.homeView()
	}

	// Build status bar
	var statusBar string
	if m.statusMsg != "" {
		statusBar = greenText.Render("✓ " + m.statusMsg)
	}

	content := lipgloss.JoinVertical(lipgloss.Top,
		mainContent,
		statusBar,
	)

	return appStyle.Render(content)
}

func (m Model) helpView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Aiuto — Comandi Disponibili"))
	b.WriteString("\n\n")

	keys := []struct {
		key string
		desc string
		ctx string
	}{
		{"↑/k, ↓/j", "Naviga lista", "Tutti"},
		{"Enter", "Seleziona / Conferma", "Tutti"},
		{"Esc / Backspace", "Torna indietro", "Tutti"},
		{"q / Ctrl+C", "Esci", "Tutti"},
		{"?", "Mostra/nascondi aiuto", "Tutti"},
		{"1-4", "Menu principale", "Home"},
		{"a", "Aggiungi pasto / alimento", "Oggi"},
		{"d", "Elimina pasto / alimento", "Oggi"},
		{"e", "Modifica nome pasto", "Oggi"},
		{"← →", "Naviga date", "Diario"},
		{"Tab", "Cambia modalità (crudo↔cotto)", "Convertitore"},
		{"Ctrl+U", "Pulisci input", "Input"},
	}

	b.WriteString(borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			lipgloss.NewStyle().Bold(true).Foreground(colorYellow).Render(fmt.Sprintf("%-20s %-35s %s", "Tasto", "Azione", "Contesto")),
			strings.Repeat("─", 70),
			func() string {
				var rows []string
				for _, k := range keys {
					rows = append(rows, fmt.Sprintf("%-20s %-35s %s", k.key, k.desc, k.ctx))
				}
				return strings.Join(rows, "\n")
			}(),
		),
	))

	b.WriteString(helpStyle.Render("\nPremi Esc per tornare indietro"))
	return b.String()
}

func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keys.Up, m.keys.Down, m.keys.Enter, m.keys.Back, m.keys.Quit, m.keys.Help,
	}
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keys.Up, m.keys.Down, m.keys.Enter, m.keys.Back, m.keys.Quit},
		{m.keys.Add, m.keys.Delete, m.keys.Edit, m.keys.Tab, m.keys.Help},
	}
}

// Truncate string with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
