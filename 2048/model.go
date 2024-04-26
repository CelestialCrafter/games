package twenty48

import (
	"fmt"
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

func initializeBoard(board [][]uint16) ([][]uint16, error) {
	if len(board) <= 0 {
		return nil, fmt.Errorf("expected board width to be > 0, got %v", len(board))
	}

	if len(board[0]) <= 0 {
		return nil, fmt.Errorf("expected board height to be > 0, got %v", len(board[0]))
	}

	for i := range 2 {
		_ = i
		x := rand.Intn(len(board))
		y := rand.Intn(len(board[0]))

		v := (rand.Intn(2) + 1) * 2

		board[x][y] = uint16(v)
	}

	return board, nil
}

type Model struct {
	Keys  KeyMap
	Help  help.Model
	Board [][]uint16
}

func (m Model) process(msg tea.Msg) {
	m.Board[0][1] = 20
	// @TODO game logic
}

func NewModel(boardWidth uint8, boardHeight uint8) Model {
	board := make([][]uint16, boardWidth)
	for i := range board {
		board[i] = make([]uint16, boardHeight)
	}

	board, err := initializeBoard(board)

	if err != nil {
		// @TODO handle this error
		panic(err)
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
	cellStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#888888")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(3).
		Height(1)

	var boardRows []string
	for y := range m.Board[0] {
		var row []string
		for x := range m.Board {
			cell := fmt.Sprint(m.Board[x][y])
			row = append(row, cellStyle.Render(cell))
		}
		boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	boardView := lipgloss.JoinVertical(lipgloss.Left, boardRows...)

	return boardView + "\n\n" + m.Help.View(m.Keys)
}
