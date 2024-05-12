package common

import (
	"github.com/charmbracelet/lipgloss"
)

func RenderBoard[T any](board [][]T, renderCell func(cell T) string) string {
	var boardRows []string
	for y := range board[0] {
		var row []string
		for x := range board {
			row = append(row, renderCell(board[x][y]))
		}
		boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, boardRows...)
}
