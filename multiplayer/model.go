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

// extra dependency but i like xsync api more so deal with it ;3
var Players = xsync.NewMapOf[string, *Player]()
var lobbies = map[uint]*list.List{
	common.TicTacToe.ID: list.New(),
}

type Model struct {
	game    uint
	self    *Player
	Element *list.Element
}

func NewModel(players int, game uint, initializeData func() interface{}) Model {
	_, ok := lobbies[game]
	if !ok {
		panic(fmt.Sprintf("game id not in lobbies map: %v", game))
	}

	var lobbyElement *list.Element = nil

	gameLobbies := lobbies[game]
	for element := gameLobbies.Front(); element != nil; element = element.Next() {
		lobby, _ := element.Value.(*Lobby)

		if lobby.Ready || lobby.hasMaxPlayers() {
			continue
		}

		lobbyElement = element
		break
	}

	if lobbyElement == nil {
		hasher := sha1.New()
		// timestamp + game id but with a bunch of type conversions
		hasher.Write([]byte(strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(int(game))))

		lobby := &Lobby{
			MaxPlayers: players,
			Players:    xsync.NewMapOf[string, *Player](),
			ID:         hex.EncodeToString(hasher.Sum(nil)),
			Data:       initializeData(),
		}

		lobbyElement = gameLobbies.PushFront(lobby)
	}

	return Model{
		game:    game,
		Element: lobbyElement,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Broadcast(msg tea.Msg) {
	lobby, _ := m.Element.Value.(*Lobby)
	lobby.Players.Range(func(_ string, player *Player) bool {
		go player.Program.Send(msg)
		return true
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.BackMsg:
		lobby, _ := m.Element.Value.(*Lobby)
		m.Broadcast(DisconnectMsg(m.self.ID))
		lobby.Players.Delete(m.self.ID)
		if lobby.Players.Size() < 1 {
			lobbies[m.game].Remove(m.Element)
		}
	case ConnectMsg:
	case SelfPlayerMsg:
		var ok bool
		m.self, ok = Players.Load(string(msg))
		if !ok {
			break
		}

		lobby, _ := m.Element.Value.(*Lobby)
		lobby.Players.Store(m.self.ID, m.self)

		m.Broadcast(ConnectMsg(m.self.ID))
		if lobby.hasMaxPlayers() {
			lobby.Ready = true
			m.Broadcast(InitialReadyMsg{})
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
