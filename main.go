package main

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/db"
)

const (
	host = "0.0.0.0"
	// i have no idea what to name this
	file = "database.db"
)

var programOpts = []tea.ProgramOption{tea.WithAltScreen()}

func main() {
	flag.Parse()

	newDb, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		log.Fatal("Could not open database", "error", err)
	}

	_ = newDb.MustExec(`
		CREATE TABLE IF NOT EXISTS saves (
		save_id TEXT PRIMARY KEY,
		owner_id TEXT NOT NULL,
		game_id INTEGER NOT NULL,
		data TEXT NOT NULL,
		file INTEGER DEFAULT 0
	)`)

	db.DB = newDb

	if *common.EnableSsh {
		startSSH()
		return
	}

	startProgram()
}
