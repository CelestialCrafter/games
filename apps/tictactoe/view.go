package tictactoe

import (
	"fmt"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

var xColor = styles.CellColors[11]
var oColor = styles.CellColors[5]

var cellStyle = lipgloss.NewStyle().
	Padding(1, 0).
	Width(7).
	Align(lipgloss.Center)

var xStyle = lipgloss.NewStyle().
	Inherit(styles.StatusStyle).
	Foreground(xColor)

var oStyle = lipgloss.NewStyle().
	Inherit(styles.StatusStyle).
	Foreground(oColor)

func winStatus(winner uint8) string {
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

func turnTextStatus(turn uint8) string {
	var turnText string
	if turn == 1 {
		turnText = xStyle.Render("x")
	} else {
		turnText = oStyle.Render("o")
	}

	return fmt.Sprintf("%v%v", turnText, styles.StatusStyle.Render("'s turn"))

}

func playerTextStatus(player uint8) string {
	var playerText string
	if player == 1 {
		playerText = xStyle.Render("x")
	} else {
		playerText = oStyle.Render("o")
	}

	return fmt.Sprintf(styles.StatusStyle.Render("you're %v"), playerText)
}

func (m Model) View() string {

	var status string
	if m.winner == 0 {
		status = turnTextStatus(m.turn)
	} else {
		status = winStatus(m.winner)
	}

	var playerStatus string
	if m.multiplayer.Lobby != nil {
		playerStatus = playerTextStatus(m.player)
	} else {
		playerStatus = ""
	}

	status = lipgloss.JoinVertical(lipgloss.Top, status, playerStatus)

	// render cell colors
	board := common.RenderBoard(m.board, func(cell uint8) string {
		newCellStyle := cellStyle.Copy()

		if cell == 1 {
			newCellStyle = newCellStyle.Background(xColor)
		} else if cell == 2 {
			newCellStyle = newCellStyle.Background(oColor)
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

	board = lipgloss.NewStyle().BorderForeground(styles.Colors.Accent).Border(lipgloss.RoundedBorder()).Render(board)

	board = lipgloss.JoinVertical(lipgloss.Top, board, status)
	board = lipgloss.Place(m.width, lipgloss.Height(board), lipgloss.Center, lipgloss.Top, board)

	help := m.help.View(m.keys)
	var multiplayer string
	if m.multiplayer.Lobby != nil {
		multiplayer = m.multiplayer.View()
	} else {
		multiplayer = ""
	}

	availableHeight := m.height
	availableHeight -= lipgloss.Height(help)
	availableHeight -= lipgloss.Height(multiplayer)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(availableHeight).Render(board),
		multiplayer,
		help,
	)
}
