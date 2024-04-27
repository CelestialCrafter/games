package twenty48

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

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

type Model struct {
	Keys     KeyMap
	Help     help.Model
	Board    [][]uint16
	Finished bool
}

func (m Model) process(msg tea.Msg) {
	changed := false

	up := func() {
		var board [][]uint16
		board, changed1 := push(m.Board)
		board, changed2 := merge(board)
		board, _ = push(board)

		if changed1 || changed2 {
			changed = true
		}

		copy(m.Board, board)
	}

	right := func() {
		m.Board = rotate90(m.Board)
		up()
		m.Board = rotateN90(m.Board)
	}

	left := func() {
		m.Board = rotateN90(m.Board)
		up()
		m.Board = rotate90(m.Board)
	}

	down := func() {
		m.Board = rotate90(rotate90(m.Board))
		up()
		m.Board = rotateN90(rotateN90(m.Board))
	}

	switch {
	case key.Matches(msg.(tea.KeyMsg), m.Keys.Up):
		up()
	case key.Matches(msg.(tea.KeyMsg), m.Keys.Down):
		down()
	case key.Matches(msg.(tea.KeyMsg), m.Keys.Left):
		left()
	case key.Matches(msg.(tea.KeyMsg), m.Keys.Right):
		right()
	}

	if changed {
		var err error
		m.Board, err = addSquare(m.Board)
		if err != nil {
			// @TODO loop over each cell and check if its adjacent cells == current
			// if atleast one is true then dont set finished to true
			m.Finished = true
		}
	}
}

func NewModel() Model {
	board := createBoard(4, 4)

	for i := range 2 {
		_ = i

		var err error
		board, err = addSquare(board)
		if err != nil {
			// @TODO handle this gracefully
			panic(err)
		}
	}

	return Model{
		Keys: KeyMap{
			Up:    key.NewBinding(key.WithKeys("k", "up", "w"), key.WithHelp("↑/k/w", "move up")),
			Down:  key.NewBinding(key.WithKeys("j", "down", "s"), key.WithHelp("↑/j/s", "move down")),
			Left:  key.NewBinding(key.WithKeys("h", "left", "a"), key.WithHelp("↑/h/a", "move left")),
			Right: key.NewBinding(key.WithKeys("l", "right", "d"), key.WithHelp("↑/l/d", "move right")),
			Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		},
		Help:  help.New(),
		Board: board,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.Keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		case key.Matches(msg, m.Keys.Up, m.Keys.Down, m.Keys.Left, m.Keys.Right):
			m.process(msg)
		}
	case tea.WindowSizeMsg:
		// @TODO handle this
		m.Help.Width = msg.Width
	}

	return m, nil
}

func (m Model) View() string {
	board := ""
	status := ""

	cellStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Width(7).
		Align(lipgloss.Center)

	if m.Finished {
		status = "u lose"
	} else {
		var boardRows []string
		for y := range m.Board[0] {
			var row []string
			for x := range m.Board {
				cell := m.Board[x][y]
				color := lipgloss.Color(fmt.Sprint(math.Log2(float64(cell))))
				cellString := fmt.Sprint(cell)
				row = append(row, cellStyle.Background(color).Render(cellString))
			}

			boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
			board = lipgloss.JoinVertical(lipgloss.Left, boardRows...)
		}
	}

	return fmt.Sprintf("%v\n%v\n\n%v", board, status, m.Help.View(m.Keys))
}
