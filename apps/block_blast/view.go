package blockblast

import (
	"fmt"
	"strings"

	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
)

var cellStyle = lipgloss.NewStyle().Width(2)

func (m Model) View() string {
	boardCopy := make(boardType, len(m.board))
	for i := range m.board {
		boardCopy[i] = make([]*lipgloss.Color, len(m.board[i]))
		copy(boardCopy[i], m.board[i])
	}

	if m.valid(m.position, false) {
		m.integrate(boardCopy, styles.Colors.Accent)
	}

	status := styles.StatusStyle.Render(fmt.Sprint(m.points, " points"))
	board := common.RenderBoard(boardCopy, func(position [2]int, cell *lipgloss.Color) string {
		if position == m.position {
			var color lipgloss.Color
			if m.valid(m.position, true) {
				color = styles.Colors.Primary
			} else {
				color = lipgloss.Color("197")
			}

			return cellStyle.Copy().Background(color).Render()
		}

		if cell == nil {
			return cellStyle.Render()
		}

		return cellStyle.Copy().Background(*cell).Render()
	})
	board = lipgloss.NewStyle().BorderForeground(styles.Colors.Accent).Border(lipgloss.RoundedBorder()).Render(board)

	colors := []lipgloss.Color{styles.Colors.Primary, styles.Colors.Secondary, styles.Colors.Accent}
	pieces := []string{}
	for i, piece := range m.pieces {
		if i == m.piece {
			continue
		}

		color := colors[i%len(colors)]
		pieces = append(pieces, common.RenderBoard(piece, func(_ [2]int, cell bool) string {
			if cell {
				return cellStyle.Copy().Background(color).Render()
			} else {
				return cellStyle.Render()
			}
		}))
	}

	combinedPieces := strings.Join(pieces, "\n\n")
	combinedPieces = lipgloss.Place(lipgloss.Width(combinedPieces)+2, lipgloss.Height(board), lipgloss.Right, lipgloss.Center, combinedPieces)

	grouped := lipgloss.JoinHorizontal(lipgloss.Top, board, combinedPieces)
	grouped = lipgloss.JoinVertical(lipgloss.Top, grouped, status)
	grouped = lipgloss.Place(m.width, lipgloss.Height(board), lipgloss.Center, lipgloss.Top, grouped)

	help := m.help.View(m.keys)

	availableHeight := m.height
	availableHeight -= lipgloss.Height(help)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(availableHeight).Render(grouped),
		help,
	)
}
