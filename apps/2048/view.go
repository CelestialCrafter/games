package twenty48

import (
	"fmt"
	"math"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	cellStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Width(7).
		Align(lipgloss.Center)

	status := ""

	if m.finished {
		status = styles.StatusStyle.Render("you lose!")
	}

	board := common.RenderBoard(m.board, func(cell uint16) string {
		color := lipgloss.Color(fmt.Sprint(math.Log2(float64(cell))))
		cellString := fmt.Sprint(cell)
		return cellStyle.Background(color).Render(cellString)
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
