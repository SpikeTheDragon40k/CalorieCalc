## CalorieTracker v1.0.0 — Initial Release

A terminal-based daily calorie tracking application with an elegant dark-mode TUI, built with Go and Bubble Tea.

### Features

- **Home screen** — 4-module launcher
- **Today's diary** — add, remove, rename meals; search and add foods with weight input
- **388-food database** — sourced from CREA/INRAN (Italy), CRÉATION (Switzerland), and BDA (European) food composition tables
- **Diary history** — day, week, and month views with visual progress bars
- **Raw → Cooked converter** — bidirectionally convert weights for pasta (×2.5), rice (×3.0), meat (×0.7), fish (×0.8), legumes (×2.5), and more
- **Configurable daily target** — default 2000 kcal, adjustable at any time
- **Persistent JSON storage** — automatic saving with atomic writes and symlink protection
- **Dark mode TUI** — powered by Bubble Tea, Bubbles, and Lip Gloss

### Keyboard shortcuts

| Key | Action |
|-----|--------|
| `↑/k` `↓/j` | Navigate |
| `Enter` | Select / confirm |
| `Esc` / `Backspace` | Go back |
| `q` / `Ctrl+C` | Quit |
| `?` | Help |
| `1`–`4` | Home menu |
| `a` | Add |
| `d` | Delete |
| `e` | Rename meal |
| `←` `→` | Browse history dates |
| `Tab` | Switch mode (converter) |

### Installation

```bash
# From source
go install github.com/yourusername/caloriecalc@latest

# Or download the pre-built binary for your platform
```

### Notes

- Data is stored in `diario.json` in the current working directory
- The food database is embedded; to add foods edit `data/alimenti.csv` and rebuild
- Requires a terminal with true color support (most modern terminals)

### Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1.3
- [Bubbles](https://github.com/charmbracelet/bubbles) v1.0
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) v1.1
- Go 1.26
