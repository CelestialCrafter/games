package main

import (
	"fmt"

	twenty48 "github.com/CelestialCrafter/games/apps/2048"
	"github.com/CelestialCrafter/games/apps/chess"
	"github.com/CelestialCrafter/games/apps/saves"
	"github.com/CelestialCrafter/games/apps/tictactoe"
	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/CelestialCrafter/games/saveManager"
	"github.com/CelestialCrafter/games/selector"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type sessionState int

const (
	selectorView sessionState = iota
	gameView
)

type KeyMap struct {
	common.ArrowsKeyMap
	Select key.Binding
	Quit   key.Binding
	Reset  key.Binding
}

type MainModel struct {
	state        sessionState
	saveManager  tea.Model
	selector     tea.Model
	app          tea.Model
	currentAppId *uint
	keys         KeyMap
	err          *common.ErrorMsg
	userId       string
	selected     int
	width        int
	height       int
}

func NewModel(userId string) MainModel {
	return MainModel{
		state:       selectorView,
		saveManager: saveManager.NewModel(userId),
		selector:    selector.NewModel(),
		app:         nil,
		userId:      userId,
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Select:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			Quit:         key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
			Reset:        common.NewResetBinding(),
		},
	}
}

func (m MainModel) NewGame(id uint) tea.Model {
	switch id {
	case common.Twenty48.ID:
		return twenty48.NewModel()
	case common.TicTacToe.ID:
		return tictactoe.NewModel()
	case common.Chess.ID:
		return chess.NewModel()
	case common.Saves.ID:
		return saves.NewModel(m.userId)
	}

	return EmptyModel{}
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(m.saveManager.Init(), m.selector.Init())
}

func (m MainModel) initializeApp(msg selector.PlayMsg) tea.Cmd {
	var loadCmd tea.Cmd
	if msg.Load {
		loadCmd = func() tea.Msg {
			return saveManager.TryLoad{
				ID: msg.ID,
			}
		}
	}

	return tea.Sequence(
		m.app.Init(),
		func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  m.width,
				Height: m.height,
			}
		},
		func() tea.Msg {
			return multiplayer.SelfPlayerMsg(m.userId)
		},
		loadCmd,
	)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Reset):
			if m.app == nil || m.currentAppId == nil {
				break
			}

			return m, tea.Sequence(
				func() tea.Msg { return common.BackMsg{} },
				func() tea.Msg {
					return selector.PlayMsg{ID: *m.currentAppId}
				},
			)
		case key.Matches(msg, m.keys.Quit):
			if m.app != nil {
				m.app.Update(common.BackMsg{})
			}

			return m, tea.Quit
		case m.err != nil && key.Matches(msg, m.keys.Left):
			m.selected = max(m.selected-1, 0)
		case m.err != nil && key.Matches(msg, m.keys.Right):
			// 2nd value is the amount of option buttons (0 based)
			m.selected = min(m.selected+1, 1)
		case m.err != nil && key.Matches(msg, m.keys.Select):
			if m.selected == 1 {
				if m.app != nil {
					m.app.Update(common.BackMsg{})
				}

				return m, tea.Quit
			}

			action := m.err.Action
			m.err = nil
			return m, action
		}
	case common.BackMsg:
		m.state = selectorView
		m.app.Update(msg)
		m.app = nil
		m.currentAppId = nil
	case selector.PlayMsg:
		m.state = gameView
		m.app = m.NewGame(msg.ID)
		m.currentAppId = &msg.ID

		return m, m.initializeApp(msg)
	case common.ErrorMsg:
		if msg.Err != nil {
			log.Error("game sent error message", "error", msg.Err)
			m.err = &msg
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.saveManager, cmd = m.saveManager.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	switch m.state {
	case selectorView:
		m.selector, cmd = m.selector.Update(msg)
	case gameView:
		m.app, cmd = m.app.Update(msg)
	}

	return m, cmd
}

func (m MainModel) getButtonStyle(option int) lipgloss.Style {
	if option == m.selected {
		return styles.ButtonSelected
	}

	return styles.Button
}

func (m MainModel) View() string {
	var s string
	switch m.state {
	case selectorView:
		s = m.selector.View()
	case gameView:
		s = m.app.View()
	}

	// error handling
	if m.err != nil {
		err := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(fmt.Sprint(m.err.Err))

		var actionText string
		if m.err.ActionText == "" {
			actionText = "Continue"
		} else {
			actionText = m.err.ActionText
		}

		actionButton := m.getButtonStyle(0).Render(actionText)
		quitButton := m.getButtonStyle(1).Render("Quit")

		var buttons string
		if m.err.Fatal {
			buttons = quitButton
		} else {
			buttons = lipgloss.JoinHorizontal(lipgloss.Top, actionButton, quitButton)
		}

		ui := lipgloss.JoinVertical(lipgloss.Center, err, buttons)
		s = lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			styles.DialogBox.Render(ui),
		)
	}

	return s
}
