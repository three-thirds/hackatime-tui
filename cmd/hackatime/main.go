package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/three-thirds/hackatime-tui/internal/app"
)

func main() {
	p := tea.NewProgram(app.NewApp())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting TUI: %v\n", err)
		os.Exit(1)
	}
}
