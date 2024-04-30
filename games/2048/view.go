package twenty48

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	board := ""
	status := ""

	cellStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Width(7).
		Align(lipgloss.Center)

	if m.finished {
		status = "u lose"
	} else {
		var boardRows []string
		for y := range m.board[0] {
			var row []string
			for x := range m.board {
				cell := m.board[x][y]
				color := lipgloss.Color(fmt.Sprint(math.Log2(float64(cell))))
				cellString := fmt.Sprint(cell)
				row = append(row, cellStyle.Background(color).Render(cellString))
			}

			boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
			board = lipgloss.JoinVertical(lipgloss.Left, boardRows...)
		}
	}

	return fmt.Sprintf("%v\n%v\n\n%v", board, status, m.help.View(m.keys))
}
