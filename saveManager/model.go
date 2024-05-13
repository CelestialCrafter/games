package saveManager

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/CelestialCrafter/games/apps/saves"
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

type Model struct {
	db       *sqlx.DB
	userId   string
	username string
}

func NewModel(db *sqlx.DB, userId string, username string) Model {
	return Model{
		db,
		userId,
		username,
	}
}

func (m Model) Init() tea.Cmd {
	// if username is blank, just insert the user
	// else, upsert the username to the current one
	upsert := " ON CONFLICT(user_id) DO UPDATE SET username=$2;"
	ignore := " OR IGNORE"
	if m.username == "" {
		upsert = ";"
	} else {
		ignore = ""
	}

	_, err := m.db.Exec(fmt.Sprintf("INSERT%v INTO users(user_id,username) VALUES($1,$2)%v", ignore, upsert),
		m.userId,
		m.username,
	)

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
		h := sha1.New()
		io.WriteString(h, fmt.Sprintf("%v-%v-%v", m.userId, msg.ID, saveFile))

		_, err := m.db.Exec(`
				INSERT INTO games(game_id,owner_id,game,data,save,last_save_time) VALUES($1,$2,$3,$4,$5,$6)
					ON CONFLICT(game_id) DO UPDATE SET data=$4;`,
			hex.EncodeToString(h.Sum(nil)),
			m.userId,
			msg.ID,
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
		// @TODO multiple save files
		save := []saves.Save{}
		saveFile := 0

		err := m.db.Select(&save, "SELECT data FROM games WHERE owner_id=$1 AND game=$2 AND save=$3;", m.userId, msg.ID, saveFile)
		if err != nil {
			log.Error("couldnt load save from database", "error", err)
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err:    fmt.Errorf("couldnt load save: %v", err),
					Action: nil,
				}
			}
		}

		if len(save) > 0 {
			return m, func() tea.Msg {
				return LoadMsg{
					Data: []byte(save[0].Data),
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return ""
}
