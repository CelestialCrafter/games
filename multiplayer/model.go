package multiplayer

import (
	"container/list"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/puzpuzpuz/xsync/v3"
)

type SelfPlayerMsg string
type Player struct {
	Program *tea.Program
	Lobby   *Lobby
	ID      string
}

type InitialReadyMsg struct{}
type ConnectMsg string
type DisconnectMsg string
type Lobby struct {
	Players    *xsync.MapOf[string, *Player]
	ID         string
	MaxPlayers int
	Ready      bool
	Data       interface{}
}

func (l Lobby) hasMaxPlayers() bool {
	return l.Players.Size() >= l.MaxPlayers
}

func (l Lobby) Broadcast(msg tea.Msg) {
	l.Players.Range(func(_ string, player *Player) bool {
		go player.Program.Send(msg)
		return true
	})
}

// extra dependency but i like xsync api more so deal with it ;3
var Players = xsync.NewMapOf[string, *Player]()
var lobbies = map[uint]*list.List{
	common.TicTacToe.ID: list.New(),
	common.Chess.ID:     list.New(),
}

type initializeDataFunc func(*xsync.MapOf[string, *Player]) interface{}

type Model struct {
	game           uint
	initializeData initializeDataFunc
	Self           *Player
	Element        *list.Element
	Lobby          *Lobby
}

func NewModel(players int, game uint, initializeData initializeDataFunc) Model {
	_, ok := lobbies[game]
	if !ok {
		panic(fmt.Sprintf("game id not in lobbies map: %v", game))
	}

	var selectedElement *list.Element = nil

	gameLobbies := lobbies[game]
	for element := gameLobbies.Front(); element != nil; element = element.Next() {
		lobby, _ := element.Value.(*Lobby)

		if lobby.Ready || lobby.hasMaxPlayers() {
			continue
		}

		selectedElement = element
		break
	}

	if selectedElement == nil {
		hasher := sha1.New()
		hasher.Write([]byte(strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(int(game))))

		lobby := &Lobby{
			MaxPlayers: players,
			Players:    xsync.NewMapOf[string, *Player](),
			ID:         hex.EncodeToString(hasher.Sum(nil)),
		}

		selectedElement = gameLobbies.PushFront(lobby)
	}

	lobby := selectedElement.Value.(*Lobby)
	return Model{
		game:           game,
		initializeData: initializeData,
		Element:        selectedElement,
		Lobby:          lobby,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.BackMsg:
		m.Lobby.Broadcast(DisconnectMsg(m.Self.ID))
		m.Lobby.Players.Delete(m.Self.ID)
		if m.Lobby.Players.Size() < 1 {
			lobbies[m.game].Remove(m.Element)
		}
	case ConnectMsg:
	case SelfPlayerMsg:
		var ok bool
		m.Self, ok = Players.Load(string(msg))
		if !ok {
			break
		}

		lobby, _ := m.Element.Value.(*Lobby)
		lobby.Players.Store(m.Self.ID, m.Self)
		m.Self.Lobby = lobby

		m.Lobby.Broadcast(ConnectMsg(m.Self.ID))
		if lobby.hasMaxPlayers() {
			lobby, _ := m.Element.Value.(*Lobby)
			lobby.Data = m.initializeData(lobby.Players)
			lobby.Ready = true
			m.Lobby.Broadcast(InitialReadyMsg{})
		}

	}

	return m, nil
}

func (m Model) View() string {
	lobby, _ := m.Element.Value.(*Lobby)
	var readyString string
	if lobby.Ready {
		readyString = "game started"
	} else {
		readyString = "waiting for players"
	}

	return styles.StatusStyle.Render(fmt.Sprintf("%d/%d players • %s", lobby.Players.Size(), lobby.MaxPlayers, readyString))
}
