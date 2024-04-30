package twenty48

import (
	"fmt"
	"math/rand"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// @TODO score
// @TODO game end logic

func addSquare(board [][]uint16) ([][]uint16, error) {
	empty := make([]*uint16, 0)

	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[0]); y++ {
			c := &board[x][y]
			if *c == 0 {
				empty = append(empty, c)
			}
		}
	}

	if len(empty) <= 0 {
		return board, fmt.Errorf("no empty spaces in board")
	}

	*empty[rand.Intn(len(empty))] = uint16((rand.Intn(2) + 1) * 2)

	return board, nil
}

func reverse(matrix [][]uint16) [][]uint16 {
	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}

	return matrix
}

func transpose(matrix [][]uint16) [][]uint16 {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < i; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	return matrix
}

func rotate90(matrix [][]uint16) [][]uint16 {
	return transpose(reverse(matrix))
}

func rotateN90(matrix [][]uint16) [][]uint16 {
	return reverse(transpose(matrix))
}

func createBoard(w int, h int) [][]uint16 {
	board := make([][]uint16, w)
	for i := range board {
		board[i] = make([]uint16, h)
	}

	return board
}

func push(board [][]uint16) ([][]uint16, bool) {
	newBoard := createBoard(len(board), len(board[0]))
	changed := false

	for i := 0; i < len(board); i++ {
		position := 0
		for j := 0; j < len(board[0]); j++ {
			current := &board[i][j]
			next := &newBoard[i][position]

			if *current != 0 {
				*next = *current
				if j != position {
					changed = true
				}
				position++
			}
		}
	}

	return newBoard, changed
}

func merge(board [][]uint16) ([][]uint16, bool) {
	changed := false

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board)-1; j++ {
			current := &board[i][j]
			next := &board[i][j+1]
			if *current == *next && *current != 0 {
				*current = *current * 2
				*next = 0
				changed = true
			}
		}
	}

	return board, changed
}

func (m Model) process(msg tea.Msg) {
	changed := false

	up := func() {
		var board [][]uint16
		board, changed1 := push(m.board)
		board, changed2 := merge(board)
		board, _ = push(board)

		if changed1 || changed2 {
			changed = true
		}

		copy(m.board, board)
	}

	right := func() {
		m.board = rotate90(m.board)
		up()
		m.board = rotateN90(m.board)
	}

	left := func() {
		m.board = rotateN90(m.board)
		up()
		m.board = rotate90(m.board)
	}

	down := func() {
		m.board = rotate90(rotate90(m.board))
		up()
		m.board = rotateN90(rotateN90(m.board))
	}

	switch {
	case key.Matches(msg.(tea.KeyMsg), m.keys.Up):
		up()
	case key.Matches(msg.(tea.KeyMsg), m.keys.Down):
		down()
	case key.Matches(msg.(tea.KeyMsg), m.keys.Left):
		left()
	case key.Matches(msg.(tea.KeyMsg), m.keys.Right):
		right()
	}

	if changed {
		var err error
		m.board, err = addSquare(m.board)
		if err != nil {
			// @TODO loop over each cell and check if its adjacent cells == current
			// if atleast one is true then dont set finished to true
			m.finished = true
		}
	}
}
