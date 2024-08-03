package tictactoe

import (
	"bytes"
	"encoding/gob"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/CelestialCrafter/games/saveManager"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type gameSave struct {
	Board [][]uint8
}
type moveMsg struct {
	position uint
	player   uint8
	turn     uint8
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case multiplayer.InitialReadyMsg:
		lobby, _ := m.multiplayer.Element.Value.(*multiplayer.Lobby)
		data, _ := lobby.Data.(*lobbyData)
		m.turn = data.startingTurn
	case moveMsg:
		m.place(msg.position, msg.player)
		m.turn = msg.turn
		m.winner = m.checkGameState()
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

			m.multiplayer.Broadcast(move)
		case key.Matches(msg, m.keys.Save):
			cmds = append(cmds, func() tea.Msg {
				bytes := bytes.Buffer{}
				encoder := gob.NewEncoder(&bytes)
				err := encoder.Encode(gameSave{
					Board: m.board,
				})

				if err != nil {
					return common.ErrorMsg{
						Err: err,
					}
				}

				return saveManager.SaveMsg{
					Data: bytes.Bytes(),
					ID:   common.TicTacToe.ID,
				}
			})
		}
	case saveManager.LoadMsg:
		saveData := gameSave{}

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
		m.winner = m.checkGameState()
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	}

	var multiplayerCmd tea.Cmd
	m.multiplayer, multiplayerCmd = m.multiplayer.Update(msg)

	cmds = append(cmds, multiplayerCmd)

	return m, tea.Batch(cmds...)
}
