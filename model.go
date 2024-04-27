package main

import (
	twenty48 "github.com/CelestialCrafter/games/2048"
	"github.com/CelestialCrafter/games/game"
	"github.com/CelestialCrafter/games/selector"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	selectorView sessionState = iota
	gameView
)

type MainModel struct {
	state    sessionState
	selector tea.Model
	game     tea.Model
	gameID   uint
}

func NewModel() MainModel {
	return MainModel{
		state:    selectorView,
		selector: selector.NewModel(),
		game:     nil,
		gameID:   0,
	}
}

func NewGame(id uint) tea.Model {
	switch id {
	case twenty48.GetMetadata().ID:
		return twenty48.NewModel()
	}

	return EmptyModel{}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case game.QuitMsg:
		m.state = selectorView
		m.game = nil
	case selector.PlayMsg:
		m.state = gameView
		m.gameID = msg.GameID

		m.game = NewGame(m.gameID)
	}

	switch m.state {
	case selectorView:
		m.selector, cmd = m.selector.Update(msg)
	case gameView:
		m.game, cmd = m.game.Update(msg)
	}

	return m, cmd
}

func (m MainModel) View() string {
	switch m.state {
	case selectorView:
		return m.selector.View()
	case gameView:
		return m.game.View()
	}

	return ""
}
