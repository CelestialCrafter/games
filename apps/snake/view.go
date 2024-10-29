package snake

import (
	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

var cellStyle = lipgloss.NewStyle().
Width(2).
Align(lipgloss.Center)

func (m Model) View() string {
	status := ""

	if m.Finished {
		status = styles.StatusStyle.Render("you lose!")
	}

	snakedBoard := common.CreateBoard[uint8](boardWidth, boardHeight)
	for i := range m.Board {
	    copy(snakedBoard[i], m.Board[i])
	}

	for _, point := range m.snake {
		snakedBoard[point.X][point.Y] = snake
	}

	board := common.RenderBoard(snakedBoard, func(cell uint8) string {
		if cell == empty {
			return cellStyle.Render()
		}

		var color lipgloss.Color
		if cell == apple {
			color = styles.Colors.Secondary
		} else {
			color = styles.Colors.Primary
		}

		return cellStyle.
		Copy().
		Background(color).
		Render()
	})

	board = lipgloss.
	NewStyle().
	BorderForeground(styles.Colors.Accent).
	Border(lipgloss.RoundedBorder()).
	Render(board)

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
