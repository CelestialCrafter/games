package twenty48

import (
	"fmt"
	"math"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

var cellStyle = lipgloss.NewStyle().
	Padding(1, 0).
	Width(7).
	Align(lipgloss.Center)

func (m Model) View() string {
	status := fmt.Sprint(m.points, " points")
	if m.finished {
		status += " • you lose!"
	}
	status = styles.StatusStyle.Render(status)

	board := common.RenderBoard(m.board, func(_ [2]int, cell uint16) string {
		index := int(math.Max(math.Log2(float64(cell)), 0))
		color := styles.CellColors[index]
		cellString := fmt.Sprint(cell)

		if cell == 0 {
			return cellStyle.Render(fmt.Sprint(cellString))
		}

		newCellStyle := cellStyle.Copy().Background(color)
		return newCellStyle.Render(cellString)
	})

	board = lipgloss.NewStyle().BorderForeground(styles.Colors.Accent).Border(lipgloss.RoundedBorder()).Render(board)

	grouped := lipgloss.JoinVertical(lipgloss.Top, board, status)
	grouped = lipgloss.Place(m.width, lipgloss.Height(board), lipgloss.Center, lipgloss.Top, board)

	help := m.help.View(m.keys)

	availableHeight := m.height
	availableHeight -= lipgloss.Height(help)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(availableHeight).Render(grouped),
		help,
	)
}
