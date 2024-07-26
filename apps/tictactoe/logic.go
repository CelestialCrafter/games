package tictactoe

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func getBoardPosition(position uint) (uint, uint) {
	x := (position - 1) % 3
	y := 2 - (position-1)/3

	return x, y
}

// ./wincondition-explanation.png
func (m Model) checkGameState() int {
	// win condition 1
	for i := 0; i < len(m.board); i++ {
		failed := false

		initial := int(m.board[i][0])
		if initial == 0 {
			failed = true
		}

		for j := 0; j < len(m.board[0]); j++ {
			current := int(m.board[i][j])

			if initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial
		}
	}

	// win condition 2
	// top left -> bottom right
	{
		failed := false

		initial := int(m.board[0][0])
		if initial == 0 {
			failed = true
		}

		for i := 0; i < len(m.board); i++ {
			current := int(m.board[i][i])
			if initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial
		}
	}

	// top right -> bottom left
	{

		failed := false

		initial := int(m.board[0][len(m.board)-1])
		if initial == 0 {
			failed = true
		}

		for i := len(m.board) - 1; i >= 0; i-- {
			current := int(m.board[len(m.board)-1-i][i])
			if initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial
		}

	}
	// win condition 3
	for i := 0; i < len(m.board[0]); i++ {
		initial := int(m.board[0][i])
		failed := false

		for j := 0; j < len(m.board); j++ {
			current := int(m.board[j][i])
			if current == 0 || initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial
		}
	}

	// draw
	failed := false
	for i := 0; i < len(m.board[0]); i++ {
		for j := 0; j < len(m.board[0]); j++ {
			if int(m.board[i][j]) == 0 {
				failed = true
				break
			}
		}
	}

	if !failed {
		return 3
	}

	return 0
}

func (m Model) place(position uint) (ok bool) {
	x, y := getBoardPosition(position)
	cell := &m.board[x][y]

	if *cell != 0 {
		return
	}

	*cell = m.turn

	return true
}

func (m *Model) process(msg tea.Msg) {
	var ok bool

	if m.winner != 0 {
		return
	}

	switch {
	// top to bottom
	case key.Matches(msg.(tea.KeyMsg), m.keys.One):
		ok = m.place(1)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Two):
		ok = m.place(2)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Three):
		ok = m.place(3)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Four):
		ok = m.place(4)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Five):
		ok = m.place(5)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Six):
		ok = m.place(6)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Seven):
		ok = m.place(7)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Eight):
		ok = m.place(8)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Nine):
		ok = m.place(9)
	}

	if ok {
		m.winner = m.checkGameState()

		if m.turn == 2 {
			m.turn = 1
		} else {
			m.turn = 2
		}
	}
}
