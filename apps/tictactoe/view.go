package tictactoe

import (
	"fmt"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	cellStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Width(7).
		Align(lipgloss.Center)

	xStyle := lipgloss.NewStyle().
		Inherit(styles.StatusStyle).
		Foreground(lipgloss.Color("4"))
	oStyle := lipgloss.NewStyle().
		Inherit(styles.StatusStyle).
		Foreground(lipgloss.Color("2"))

	status := ""

	var turn string
	if m.turn == 1 {
		turn = xStyle.Render("x")
	} else {
		turn = oStyle.Render("o")
	}

	status = fmt.Sprintf("%v%v", turn, styles.StatusStyle.Render("'s turn"))

	if m.winner != 0 {
		var winner string
		if m.winner == 1 {
			winner = xStyle.Render("x")
		} else {
			winner = oStyle.Render("o")
		}

		status = fmt.Sprintf("%v %v", winner, styles.StatusStyle.Render("wins!"))
	}

	board := common.RenderBoard(m.board, func(cell uint8) string {
		var color lipgloss.Color

		if cell == 0 {
			color = lipgloss.Color("0")
		} else if cell == 1 {
			color = lipgloss.Color("4")
		} else {
			color = lipgloss.Color("2")
		}

		var cellString string
		if cell == 0 {
			cellString = "-"
		} else if cell == 1 {
			cellString = "x"
		} else {
			cellString = "o"
		}

		return cellStyle.Background(color).Render(cellString)
	})

	return fmt.Sprintf("%v\n\n%v\n%v", board, status, m.help.View(m.keys))
}
