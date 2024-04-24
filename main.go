package main

import (
	twenty48 "games/2048"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(twenty48.NewModel(8, 8))
	_, err := program.Run()
	if err != nil {
		panic(err)
	}
}
