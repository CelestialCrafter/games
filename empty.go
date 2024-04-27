package main

import (
	"github.com/CelestialCrafter/games/game"
	tea "github.com/charmbracelet/bubbletea"
)

type EmptyModel struct {
}

func (m EmptyModel) Init() tea.Cmd {
	return func() tea.Msg {
		return game.QuitMsg{}
	}
}

func (m EmptyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		return game.QuitMsg{}
	}
}

func (m EmptyModel) View() string {
	return ""
}
