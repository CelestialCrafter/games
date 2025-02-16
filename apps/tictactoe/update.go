package tictactoe

import (
	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type moveMsg struct {
	position uint
	player   uint8
	turn     uint8
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case multiplayer.InitialReadyMsg:
		data, _ := m.multiplayer.Lobby.Data.(*lobbyData)
		m.turn = data.startingTurn
		m.player = data.colors[m.multiplayer.Self.ID]
	case multiplayer.DisconnectMsg:
		data, _ := m.multiplayer.Lobby.Data.(*lobbyData)
		loser := data.colors[string(msg)]

		if loser == 1 {
			m.winner = 2
		} else {
			m.winner = 1
		}
	case moveMsg:
		m.place(msg.position, msg.player)
		m.turn = msg.turn
		m.winner = m.checkGameState()
		if m.multiplayer.Lobby == nil {
			m.player = msg.turn
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(
			msg,
			m.keys.One,
			m.keys.Two,
			m.keys.Three,
			m.keys.Four,
			m.keys.Five,
			m.keys.Six,
			m.keys.Seven,
			m.keys.Eight,
			m.keys.Nine,
		):
			move, ok := m.process(msg)
			if !ok {
				break
			}

			if m.multiplayer.Lobby != nil {
				m.multiplayer.Lobby.Broadcast(move)
			} else {
				cmds = append(cmds, func() tea.Msg {
					return move
				})
			}
		}
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.multiplayer.Lobby != nil {
		var multiplayerCmd tea.Cmd
		m.multiplayer, multiplayerCmd = m.multiplayer.Update(msg)

		cmds = append(cmds, multiplayerCmd)

	}

	return m, tea.Batch(cmds...)
}
