package common

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

func CreateBoard[T any](w int, h int) [][]T {
	board := make([][]T, w)
	for i := range board {
		board[i] = make([]T, h)
	}

	return board
}

func CompareBoards[T ~[][]E, E comparable](a T, b T) (equal bool) {
	if len(a) != len(b) {
		return false
	}

	for i, aRow := range a {
		bRow := b[i]
		if !slices.Equal(aRow, bRow) {
			return false
		}
	}

	return true
}

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
