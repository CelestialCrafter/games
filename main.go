package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/CelestialCrafter/games/db"
)

const (
	host = "localhost"
	port = "2222"
	// i have no idea what to name this
	file = "database.db"
)

var programOpts = []tea.ProgramOption{tea.WithAltScreen()}

func createTeaHandler() func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		key := sess.PublicKey()
		if key == nil {
			log.Error("Key was nil (enable PublicKeyAuth middleware?)")
			_, err := sess.Write([]byte("You need to ssh in with a public key!"))
			if err != nil {
				log.Error("Could not write error message to session", "error", err)
			}

			err = sess.Close()
			if err != nil {
				log.Error("Could not close session due to nil key", "error", err)
			}
		}

		keyType := key.Type()
		keyData := hex.EncodeToString(key.Marshal())

		m := NewModel(fmt.Sprintf("%v-%v", keyType, keyData))
		renderer := bubbletea.MakeRenderer(sess)
		lipgloss.SetDefaultRenderer(renderer)
		return m, programOpts
	}
}

func startProgram() {
	_, err := tea.NewProgram(NewModel("default"), programOpts...).Run()
	if err != nil {
		log.Error("Could not start program", "error", err)
	}
}

func startSSH() {
	addr := net.JoinHostPort(host, port)

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			bubbletea.Middleware(createTeaHandler()),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not stop server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func main() {
	newDb, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		log.Error("Could not open database", "error", err)
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

	if len(os.Args) >= 2 && os.Args[1] == "ssh" {
		startSSH()
	} else {
		// logging messes with the TUI, so we have to write logs to a file when not running as a ssh server
		f, err := tea.LogToFileWith("program.log", "debug", log.Default())
		if err != nil {
			log.Fatal("could not open log file", "error", err)
		}
		defer f.Close()

		startProgram()
	}
}
