package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	twenty48 "github.com/CelestialCrafter/games/2048"
	"github.com/CelestialCrafter/games/metadata"
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
	Play  key.Binding
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

type Model struct {
	Keys         KeyMap
	Help         help.Model
	Games        []Game
	SelectedGame int
	RowLength    int
}

func NewModel() Model {
	return Model{
		Keys: KeyMap{
			Up:    key.NewBinding(key.WithKeys("k", "up", "w"), key.WithHelp("↑/k/w", "move up")),
			Down:  key.NewBinding(key.WithKeys("j", "down", "s"), key.WithHelp("↑/j/s", "move down")),
			Left:  key.NewBinding(key.WithKeys("h", "left", "a"), key.WithHelp("↑/h/a", "move left")),
			Right: key.NewBinding(key.WithKeys("l", "right", "d"), key.WithHelp("↑/l/d", "move right")),
			Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		},
		Help: help.New(),
		Games: []Game{
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
			twenty48.NewModel(),
		},
		RowLength: 5,
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
		case key.Matches(msg, m.Keys.Up):
			m.SelectedGame -= m.RowLength
		case key.Matches(msg, m.Keys.Down):
			m.SelectedGame += m.RowLength
		case key.Matches(msg, m.Keys.Right):
			m.SelectedGame++
		case key.Matches(msg, m.Keys.Left):
			m.SelectedGame--
		}

		m.SelectedGame = min(max(m.SelectedGame, 0), len(m.Games)-1)
	case tea.WindowSizeMsg:
		// -1 is to account for margin
		m.RowLength = msg.Width/metadata.ICON_WIDTH - 1
		m.Help.Width = msg.Width
	}

	return m, nil
}

func (m Model) View() string {
	selectedBar := fmt.Sprintf(
		"\n%v", lipgloss.NewStyle().
			Width(metadata.ICON_WIDTH).
			Height(1).
			Margin(0, 1).
			Render(strings.Repeat("━", metadata.ICON_WIDTH)),
	)

	rowAmount := int(math.Ceil(float64(len(m.Games)) / float64(m.RowLength)))
	menu := make([][]string, rowAmount)
	for i := range menu {
		menu[i] = make([]string, m.RowLength)
	}

	menuRows := make([]string, rowAmount)

	for i := 0; i < len(m.Games); i++ {
		currentRow := int(math.Floor(float64(i) / float64(m.RowLength)))
		currentColumn := i % m.RowLength
		current := &menu[currentRow][currentColumn]

		*current = m.Games[i].GetMetadata().Icon
		if m.SelectedGame == i {
			*current = strings.TrimRightFunc(*current, unicode.IsSpace) + selectedBar
		}
	}

	for i := 0; i < len(menu); i++ {
		menuRows[i] = lipgloss.JoinHorizontal(lipgloss.Top, menu[i]...)
	}

	menuString := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	return fmt.Sprintf("%v\n%v", menuString, m.Help.View(m.Keys))
}
