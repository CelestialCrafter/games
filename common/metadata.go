package common

import (
	"github.com/charmbracelet/lipgloss"
)

const ICON_WIDTH, ICON_HEIGHT = 15, 7

type Metadata struct {
	Name string
	// icons should be 15x7
	Icon string
	// loading, saving, and any other common features the game may support
	Features []string
	ID       uint
}

var (
	Twenty48 = (func() Metadata {
		logo := lipgloss.NewStyle().
			Background(lipgloss.Color("#ffcc33")).
			BorderForeground(lipgloss.Color("#ffcc33")).
			Align(lipgloss.Center).
			Padding(3, 6).
			// this has weird positioning and it annoys me
			PaddingRight(5).
			Margin(1, 2).
			Render("2048")

		return Metadata{
			Name:     "2048",
			Features: []string{"saving"},
			Icon:     logo,
			ID:       0,
		}
	})()

	TicTacToe = (func() Metadata {
		logo := lipgloss.NewStyle().
			Background(lipgloss.Color("1")).
			Align(lipgloss.Center).
			Padding(3, 6).
			Margin(1, 2).
			Render("X/O")

		return Metadata{
			Name:     "TicTacToe",
			Features: []string{"saving"},
			Icon:     logo,
			ID:       1,
		}
	})()

	Saves = (func() Metadata {
		logo := lipgloss.NewStyle().
			Background(lipgloss.Color("8")).
			BorderForeground(lipgloss.Color("8")).
			Align(lipgloss.Center).
			Padding(3, 5).
			Margin(1, 2).
			Render("Saves")

		return Metadata{
			Name: "Saves",
			Icon: logo,
			ID:   2,
		}
	})()
)

var Games = map[uint]Metadata{
	Twenty48.ID:  Twenty48,
	TicTacToe.ID: TicTacToe,
	Saves.ID:     Saves,
}
