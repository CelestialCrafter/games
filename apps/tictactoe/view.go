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
		newCellStyle := cellStyle.Copy()

		if cell == 1 {
			newCellStyle = newCellStyle.Background(lipgloss.Color("2"))
		} else if cell == 2 {
			newCellStyle = newCellStyle.Background(lipgloss.Color("6"))
		}

		var cellString string
		if cell == 0 {
			cellString = "-"
		} else if cell == 1 {
			cellString = "x"
		} else {
			cellString = "o"
		}

		return newCellStyle.Render(cellString)
	})

	board = lipgloss.NewStyle().BorderForeground(lipgloss.Color("2")).Border(lipgloss.RoundedBorder()).Render(board)

	board = lipgloss.JoinVertical(lipgloss.Top, board, status)
	board = lipgloss.Place(m.width, lipgloss.Height(board), lipgloss.Center, lipgloss.Top, board)

	help := m.help.View(m.keys)

	availableHeight := m.height
	availableHeight -= lipgloss.Height(help)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(availableHeight).Render(board),
		help,
	)
}
