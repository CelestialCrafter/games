package tictactoe

import (
	"fmt"

	"github.com/CelestialCrafter/games/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct{}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func GetMetadata() common.Metadata {
	r1 := "─┘─└─\n"
	r2 := "─┤─├─\n"
	r3 := "─┐─┌─\n"
	board := fmt.Sprint(r1, r2, r3)

	logo := lipgloss.NewStyle().
		Background(lipgloss.Color("1")).
		Align(lipgloss.Center).
		Padding(1, 5).
		Margin(1, 2).
		Render(board)

	return common.Metadata{
		Name:     "TicTacToe",
		Features: []string{},
		Icon:     logo,
		ID:       2,
	}
}
