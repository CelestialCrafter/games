package save

import (
	"fmt"
	"time"

	"github.com/CelestialCrafter/games/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

type TryLoad struct {
	ID uint
}

type SaveMsg struct {
	Data []byte
	ID   uint
}

type LoadMsg struct {
	Data []byte
}

type save struct {
	Id           int       `db:"id"`
	GameId       int       `db:"game_id"`
	UserId       int       `db:"user_id"`
	Data         string    `db:"data"`
	LastSaveTime time.Time `db:"last_save_time"`
}

type Model struct {
	db     *sqlx.DB
	userID string
}

func NewModel(db *sqlx.DB, userID string) Model {
	return Model{
		db,
		userID,
	}
}

func (m Model) Init() tea.Cmd {
	// @TODO create user if they dont exist
	_, err := m.db.Exec("INSERT INTO users(user_id) VALUES($1);", m.userID)

	if err != nil {
		log.Error("couldnt create user", "error", err)
		return func() tea.Msg {
			return common.ErrorMsg{
				Err:    fmt.Errorf("couldnt create user: %v", err),
				Action: nil,
				Fatal:  true,
			}
		}
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SaveMsg:
		saveFile := 0

		_, err := m.db.Exec(`
				INSERT INTO games(owner_id,game,data,save,last_save_time) VALUES($1,$2,$3,$4,$5)
					ON CONFLICT(id) DO UPDATE SET data=$2;`,
			m.userID,
			msg.ID,
			// @TODO escape the save data
			string(msg.Data),
			saveFile,
			time.Now(),
		)

		if err != nil {
			log.Error("couldnt save to database", "error", err)
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err:    fmt.Errorf("couldnt save: %v", err),
					Action: nil,
				}
			}
		}
	case TryLoad:
		// @TODO multiple save files w a save manager
		save := save{}
		saveFile := 0

		err := m.db.Get(&save, "SELECT data from games WHERE owner_id=$1 AND game=$2 AND save=$3;", m.userID, msg.ID, saveFile)
		if err != nil {
			log.Error("couldnt load save from database", "error", err)
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err:    fmt.Errorf("couldnt load save: %v", err),
					Action: nil,
				}
			}
		}

		// @TODO multiple save files
		return m, func() tea.Msg {
			return LoadMsg{
				Data: []byte(save.Data),
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return ""
}
