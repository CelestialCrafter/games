package saves

import (
	"fmt"
	"time"

	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type Save struct {
	Id           string    `db:"game_id"`
	UserId       string    `db:"user_id"`
	GameId       uint      `db:"game"`
	Save         uint      `db:"save"`
	Data         string    `db:"data"`
	LastSaveTime time.Time `db:"last_save_time"`
}

func (s Save) Title() string {
	return fmt.Sprintf("%v - File %v", common.Games[s.GameId].Name, s.Save)
}
func (s Save) Description() string {
	return fmt.Sprintf("Created at %v", s.LastSaveTime.Format(time.Stamp))
}
func (s Save) FilterValue() string { return s.Title() }

type KeyMap struct {
	common.ArrowsKeyMap
	Delete key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Delete, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Delete, k.Help, k.Quit},
	}
}

type Model struct {
	keys   KeyMap
	help   help.Model
	list   list.Model
	db     *sqlx.DB
	userId string
	saves  []Save
}

func NewModel(db *sqlx.DB, userId string) Model {
	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Delete:       key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete save")),
			Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
		},
		help:   help.New(),
		list:   list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		db:     db,
		userId: userId,
	}
}

func (m Model) Init() tea.Cmd {
	m.list.Title = "Saves"

	err := m.db.Select(&m.saves, "SELECT id, game_id, last_save_time FROM games WHERE user_id=$1", m.userId)
	// life if messages were commands 🤤
	if err != nil {
		return func() tea.Msg {
			return common.ErrorMsg{
				Err: err,
				Action: func() tea.Msg {
					return common.BackMsg{}
				},
				ActionText: "Back",
			}
		}
	}

	items := make([]list.Item, len(m.saves))
	for i, v := range m.saves {
		items[i] = v
	}

	m.list.SetItems(items)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return listStyle.Render(m.list.View())
}
