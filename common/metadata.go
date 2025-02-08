package common

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	ICON_WIDTH  = 15
	ICON_HEIGHT = 7
)

type Metadata struct {
	Name string
	// icons should be 15x7
	Icon string
	// multiplayer, and any other common features the game may support
	Features []string
	ID       uint
}

var (
	Twenty48 = (func() Metadata {
		icon := lipgloss.NewStyle().
			Background(lipgloss.Color("214")).
			Align(lipgloss.Center).
			Padding(3, 6).
			// this has weird positioning and it annoys me
			PaddingRight(5).
			Margin(1, 2).
			Render("2048")

		return Metadata{
			Name:     "2048",
			Features: []string{},
			Icon:     icon,
			ID:       0,
		}
	})()

	TicTacToe = (func() Metadata {
		icon := lipgloss.NewStyle().
			Background(lipgloss.Color("168")).
			Align(lipgloss.Center).
			Padding(3, 6).
			Margin(1, 2).
			Render("X/O")

		return Metadata{
			Name:     "TicTacToe",
			Features: []string{"multiplayer"},
			Icon:     icon,
			ID:       1,
		}
	})()

	Chess = (func() Metadata {
		icon := lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Align(lipgloss.Center).
			Padding(3, 6).
			Margin(1, 2).
			Render("♚ ♖")

		return Metadata{
			Name:     "TicTacToe",
			Features: []string{"multiplayer"},
			Icon:     icon,
			ID:       2,
		}
	})()

	Snake = (func() Metadata {
		icon := lipgloss.NewStyle().
			Background(lipgloss.Color("70")).
			Align(lipgloss.Center).
			Padding(3, 5).
			Margin(1, 2).
			Render("Snake")

		return Metadata{
			Name:     "Snake",
			Features: []string{},
			Icon:     icon,
			ID:       3,
		}
	})()

	BlockBlast = (func() Metadata {
		icon := lipgloss.NewStyle().
			Background(lipgloss.Color("27")).
			Align(lipgloss.Center).
			Padding(3, 2).
			Margin(1, 2).
			Render("Block Blast")

		return Metadata{
			Name:     "Block Blast",
			Features: []string{},
			Icon:     icon,
			ID:       4,
		}
	})()
)

var Games = map[uint]Metadata{
	Twenty48.ID:   Twenty48,
	TicTacToe.ID:  TicTacToe,
	Chess.ID:      Chess,
	Snake.ID:      Snake,
	BlockBlast.ID: BlockBlast,
}
