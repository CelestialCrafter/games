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
func (m Model) checkGameState() uint8 {
	// win condition 1
	for i := 0; i < len(m.board); i++ {
		failed := false

		initial := m.board[i][0]
		if initial == 0 {
			failed = true
		}

		for j := 0; j < len(m.board[0]); j++ {
			current := m.board[i][j]

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

		initial := m.board[0][0]
		if initial == 0 {
			failed = true
		}

		for i := 0; i < len(m.board); i++ {
			current := m.board[i][i]
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

		initial := m.board[0][len(m.board)-1]
		if initial == 0 {
			failed = true
		}

		for i := len(m.board) - 1; i >= 0; i-- {
			current := m.board[len(m.board)-1-i][i]
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
		initial := m.board[0][i]
		failed := false

		for j := 0; j < len(m.board); j++ {
			current := m.board[j][i]
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
			if m.board[i][j] == 0 {
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

func (m Model) place(position uint, player uint8) {
	x, y := getBoardPosition(position)
	cell := &m.board[x][y]

	*cell = player
}

func (m Model) placeCheck(position uint) (ok bool) {
	x, y := getBoardPosition(position)
	cell := &m.board[x][y]
	return *cell == 0

}

func (m *Model) process(msg tea.Msg) (tea.Msg, bool) {
	if m.winner != 0 || m.turn != m.player || !m.ready {
		return nil, false
	}

	var p uint

	switch {
	// top to bottom
	case key.Matches(msg.(tea.KeyMsg), m.keys.One):
		p = 1
	case key.Matches(msg.(tea.KeyMsg), m.keys.Two):
		p = 2
	case key.Matches(msg.(tea.KeyMsg), m.keys.Three):
		p = 3
	case key.Matches(msg.(tea.KeyMsg), m.keys.Four):
		p = 4
	case key.Matches(msg.(tea.KeyMsg), m.keys.Five):
		p = 5
	case key.Matches(msg.(tea.KeyMsg), m.keys.Six):
		p = 6
	case key.Matches(msg.(tea.KeyMsg), m.keys.Seven):
		p = 7
	case key.Matches(msg.(tea.KeyMsg), m.keys.Eight):
		p = 8
	case key.Matches(msg.(tea.KeyMsg), m.keys.Nine):
		p = 9
	}

	if !m.placeCheck(p) {
		return nil, false
	}

	var nextTurn uint8
	if m.player == 1 {
		nextTurn = 2
	} else {
		nextTurn = 1
	}

	return moveMsg{
		position: p,
		player:   m.player,
		turn:     nextTurn,
	}, true
}
