# CalorieTracker

A terminal-based (TUI) daily calorie tracking application written in Go.

![License](https://img.shields.io/badge/license-GPLv3-blue)
![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)

Track your meals, log calories, browse history, and convert between raw and cooked food weights — all from your terminal.

## Features

- **Home screen** with quick access to all modules
- **Today's diary** — add, remove, and rename meals; add foods with fuzzy search; real-time calorie totals
- **Nutrition database** — 388+ foods pre-loaded from CREA/INRAN and CRÉATION sources, covering Italian and European food categories
- **Diary history** — day, week, and month views with visual calorie bars
- **Raw / Cooked converter** — bidirectionally convert weights for pasta, rice, meat, fish, legumes, and more
- **Configurable target** — set your daily calorie goal (default 2000 kcal)
- **Persistent JSON storage** — `diario.json` saved atomically; safe against corruption
- **Dark mode** — clean, modern TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Screenshots

```
┌──────────────────────────────────────────┐
│  🥗 CalorieTracker                        │
│                                           │
│  [1] 📅 Today — Log meals for today       │
│  [2] 📊 Diary — Day / Week / Month view   │
│  [3] 🔄 Raw → Cooked — Weight converter   │
│  [4] ⚙️  Settings — Target & meal config  │
│                                           │
│  1-4 select  ? help  q quit               │
└──────────────────────────────────────────┘
```

## Installation

### Prerequisites

- Go 1.26 or later

### From source

```bash
git clone https://github.com/yourusername/caloriecalc.git
cd caloriecalc
go build -o caloriecalc .
```

### Pre-built binary

Download the latest release for your platform from the [Releases page](https://github.com/yourusername/caloriecalc/releases).

### Run

```bash
./caloriecalc
```

The app saves your data to `diario.json` in the current directory. You can resume later by running the same command from the same directory.

## Usage

| Key | Action |
|-----|--------|
| `↑/k` `↓/j` | Navigate lists |
| `Enter` | Select / confirm |
| `Esc` / `Backspace` | Go back |
| `q` / `Ctrl+C` | Quit |
| `?` | Toggle help |
| `1`–`4` | Home screen menu |
| `a` | Add meal / food |
| `d` | Delete meal / food |
| `e` | Rename meal |
| `←` `→` | Browse dates in diary |
| `Tab` | Switch mode (converter) |

### Adding food to a meal

1. From **Today**, select a meal and press `Enter`
2. Press `a` to open the food search
3. Type to filter — both substring and fuzzy matching are supported
4. Navigate results and press `Enter`
5. Enter the weight in grams and confirm

### Raw → Cooked converter

The converter supports these foods with their expansion factors:

| Food | Raw → Cooked factor |
|------|--------------------|
| Dry pasta | ×2.5 |
| Dry egg pasta | ×2.4 |
| Dry rice | ×3.0 |
| Brown rice (dry) | ×2.8 |
| Farro, barley, couscous, quinoa (dry) | ×2.5–2.8 |
| Dry legumes | ×2.5 |
| Potatoes (raw) | ×1.0 |
| Meat (raw) | ×0.7 |
| Fish (raw) | ×0.8 |
| Vegetables (raw) | ×0.9 |

Press `Tab` to switch between raw→cooked and cooked→raw modes.

## Food database

The built-in database contains 388 foods sourced from:
- **CREA** (Consiglio per la ricerca in agricoltura e l'analisi dell'economia agraria) — Italian food composition tables
- **CRÉATION** — Swiss food database (original source of the legacy `.xls` file)
- **BDA** (Banca Dati di Composizione degli Alimenti) — European reference

Categories include: meat, fish, dairy, pasta & cereals, fruit, vegetables, legumes, nuts, oils, sweets, beverages, and more.

To add more foods, edit `data/alimenti.csv` and copy it to `internal/parser/alimenti.csv` before rebuilding.

## Security

- Atomic file writes prevent data corruption
- Symlink protection prevents file redirection attacks
- JSON input is limited to 10 MB and validated with strict decoding
- Gram input is bounded (0–5000 g)
- Panic recovery ensures terminal is always reset

## License

Copyright (C) 2026 Adriano Inghingolo

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — UI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling
- [Go](https://go.dev/) — programming language
