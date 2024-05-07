package twenty48

import (
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

type gameSave struct {
	Board [][]uint16
}

type Model struct {
	keys     KeyMap
	help     help.Model
	board    [][]uint16
	finished bool
}

func NewModel() Model {
	board := createBoard(4, 4)

	for i := range 2 {
		_ = i

		// this should never error
		board = addSquare(board)
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

func GetMetadata() common.Metadata {
	logo := lipgloss.NewStyle().
		Background(lipgloss.Color("#ffcc33")).
		BorderForeground(lipgloss.Color("#ffcc33")).
		Align(lipgloss.Center).
		Padding(3, 6).
		// this has weird positioning and it annoys me
		PaddingRight(5).
		Margin(1, 2).
		Render("2048")

	return common.Metadata{
		Name:     "2048",
		Features: []string{"saving"},
		Icon:     logo,
		ID:       1,
	}
}
