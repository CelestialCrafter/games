package main

import (
	"github.com/CelestialCrafter/games/metadata"

	tea "github.com/charmbracelet/bubbletea"
)

type Game interface {
	GetMetadata() metadata.Metadata
	tea.Model
}

func main() {
	program := tea.NewProgram(NewModel())
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}
