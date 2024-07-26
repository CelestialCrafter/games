package tictactoe

import (
	"fmt"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

var cellStyle = lipgloss.NewStyle().
	Padding(1, 0).
	Width(7).
	Align(lipgloss.Center)

var xStyle = lipgloss.NewStyle().
	Inherit(styles.StatusStyle).
	Foreground(lipgloss.Color("4"))

var oStyle = lipgloss.NewStyle().
	Inherit(styles.StatusStyle).
	Foreground(lipgloss.Color("2"))

func winStatus(winner int) string {

	// draw
	if winner == 3 {
		return styles.StatusStyle.Render("it's a draw!")

	}

	var winnerText string
	if winner == 1 {
		winnerText = xStyle.Render("x")
	} else {
		winnerText = oStyle.Render("o")
	}

	return fmt.Sprintf("%v %v", winnerText, styles.StatusStyle.Render("wins!"))
}

func turnTextStatus(turn int) string {

	var turnText string
	if turn == 1 {
		turnText = xStyle.Render("x")
	} else {
		turnText = oStyle.Render("o")
	}

	return fmt.Sprintf("%v%v", turnText, styles.StatusStyle.Render("'s turn"))

}

func (m Model) View() string {

	var status string
	if m.winner == 0 {
		status = turnTextStatus(int(m.turn))
	} else {
		status = winStatus(m.winner)
	}

	// render cell colors
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
