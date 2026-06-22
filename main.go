package main

import (
	"caloriecalc/internal/tui"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Panic: %v\n", r)
			fmt.Println("Terminal may be in raw mode. Run 'stty sane' to reset.")
			os.Exit(1)
		}
	}()

	p := tea.NewProgram(
		tui.NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		<-sigCh
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
