package twenty48

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"

	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	common.ArrowsKeyMap
	Save key.Binding
	Help key.Binding
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Save, k.Help, k.Quit},
	}
}

type save struct {
	Board [][]uint16
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
	keys     KeyMap
	help     help.Model
	board    [][]uint16
	finished bool
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

func NewModel() Model {
	board := createBoard(4, 4)

	for i := range 2 {
		_ = i

		// this should never error
		board, _ = addSquare(board)
	}

	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Save:         key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
			Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
		},
		help:  help.New(),
		board: board,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Up, m.keys.Down, m.keys.Left, m.keys.Right):
			m.process(msg)
		case key.Matches(msg, m.keys.Save):
			return m, func() tea.Msg {
				save := save{
					Board: m.board,
				}

				bytes := bytes.Buffer{}
				encoder := gob.NewEncoder(&bytes)
				err := encoder.Encode(save)

				if err != nil {
					return common.ErrorMsg{
						Err: err,
					}
				}

				return common.SaveMsg{
					Data: bytes.Bytes(),
				}
			}
		}
	case common.LoadMsg:
		saveData := save{}

		bytes := bytes.Buffer{}
		bytes.Write(msg.Data)
		decoder := gob.NewDecoder(&bytes)
		err := decoder.Decode(&saveData)

		if err != nil {
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err: err,
				}
			}
		}

		m.board = saveData.Board
	case tea.WindowSizeMsg:
		// @TODO handle this
		m.help.Width = msg.Width
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

	if m.finished {
		status = "u lose"
	} else {
		var boardRows []string
		for y := range m.board[0] {
			var row []string
			for x := range m.board {
				cell := m.board[x][y]
				color := lipgloss.Color(fmt.Sprint(math.Log2(float64(cell))))
				cellString := fmt.Sprint(cell)
				row = append(row, cellStyle.Background(color).Render(cellString))
			}

			boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
			board = lipgloss.JoinVertical(lipgloss.Left, boardRows...)
		}
	}

	return fmt.Sprintf("%v\n%v\n\n%v", board, status, m.help.View(m.keys))
}

func GetMetadata() common.Metadata {
	logo := lipgloss.NewStyle().
		Background(lipgloss.Color("#ffcc33")).
		BorderForeground(lipgloss.Color("#ffcc33")).
		Align(lipgloss.Center).
		Padding(3, 5).
		Margin(1, 1).
		Render("2048")

	return common.Metadata{
		Name:     "2048",
		Features: []string{},
		Icon:     logo,
		ID:       1,
	}
}
