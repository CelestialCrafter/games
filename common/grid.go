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

func ReverseBoard[T ~[][]E, E comparable](board T) {
	for i, j := 0, len(board)-1; i < j; i, j = i+1, j-1 {
		board[i], board[j] = board[j], board[i]
	}
}

func TransposeBoard[T ~[][]E, E comparable](board T) {
	for i := 0; i < len(board); i++ {
		for j := 0; j < i; j++ {
			board[i][j], board[j][i] = board[j][i], board[i][j]
		}
	}
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

func RenderBoard[T any](board [][]T, renderCell func(position [2]int, cell T) string) string {
	var boardRows []string
	for y := range board[0] {
		var row []string
		for x := range board {
			row = append(row, renderCell([2]int{x, y}, board[x][y]))
		}
		boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, boardRows...)
}
