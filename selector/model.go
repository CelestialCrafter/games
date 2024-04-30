package selector

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	twenty48 "github.com/CelestialCrafter/games/2048"
	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	common.ArrowsKeyMap
	Play key.Binding
	Help key.Binding
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Play, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Play, k.Help, k.Quit},
	}
}

type PlayMsg struct {
	GameID uint
}

type Model struct {
	keys          KeyMap
	help          help.Model
	gamesMetadata []common.Metadata
	selectedGame  int
	rowLength     int
}

func NewModel() Model {
	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
			Play:         key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "play game")),
			Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		},
		help: help.New(),
		gamesMetadata: []common.Metadata{
			twenty48.GetMetadata(),
		},
		rowLength: 5,
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
			return m, tea.Quit
		case key.Matches(msg, m.keys.Play):
			return m, func() tea.Msg {
				return PlayMsg{
					GameID: m.gamesMetadata[m.selectedGame].ID,
				}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Up):
			m.selectedGame -= m.rowLength
		case key.Matches(msg, m.keys.Down):
			m.selectedGame += m.rowLength
		case key.Matches(msg, m.keys.Right):
			m.selectedGame++
		case key.Matches(msg, m.keys.Left):
			m.selectedGame--
		}

		m.selectedGame = min(max(m.selectedGame, 0), len(m.gamesMetadata)-1)
	case tea.WindowSizeMsg:
		// -1 is to account for margin
		m.rowLength = msg.Width/common.ICON_WIDTH - 1
		m.help.Width = msg.Width
	}

	return m, nil
}

func (m Model) View() string {
	selectedBar := fmt.Sprintf(
		"\n%v", lipgloss.NewStyle().
			Width(common.ICON_WIDTH).
			Height(1).
			Margin(0, 1).
			Render(strings.Repeat("━", common.ICON_WIDTH)),
	)

	rowAmount := int(math.Ceil(float64(len(m.gamesMetadata)) / float64(m.rowLength)))
	menu := make([][]string, rowAmount)
	for i := range menu {
		menu[i] = make([]string, m.rowLength)
	}

	menuRows := make([]string, rowAmount)

	for i := 0; i < len(m.gamesMetadata); i++ {
		currentRow := int(math.Floor(float64(i) / float64(m.rowLength)))
		currentColumn := i % m.rowLength
		current := &menu[currentRow][currentColumn]

		*current = m.gamesMetadata[i].Icon
		if m.selectedGame == i {
			*current = strings.TrimRightFunc(*current, unicode.IsSpace) + selectedBar
		}
	}

	for i := 0; i < len(menu); i++ {
		menuRows[i] = lipgloss.JoinHorizontal(lipgloss.Top, menu[i]...)
	}

	menuString := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	return fmt.Sprintf("%v\n%v", menuString, m.help.View(m.keys))
}
