package main

import (
	"fmt"
	"os"

	"tunepipe/internal/player"
	"tunepipe/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	pipeName := "tunepipe-mpv"
	eng, err := player.NewEngine(pipeName)
	if err != nil {
		fmt.Printf("Erreur au demarrage de mpv: %v\n", err)
		os.Exit(1)
	}
	defer eng.Stop()

	p := tea.NewProgram(ui.NewModel(eng), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erreur TUI: %v\n", err)
		os.Exit(1)
	}
}
