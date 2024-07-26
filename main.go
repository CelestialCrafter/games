package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
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
	"github.com/muesli/termenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/CelestialCrafter/games/db"
)

const (
	host = "0.0.0.0"
	port = "2222"
	// i have no idea what to name this
	file = "database.db"
)

var (
	enableSsh = flag.Bool("ssh", false, "turns into a ssh server")
	// https://github.com/muesli/termenv/blob/51d72d34e2b9778a31aa5dd79fbdd8cdac50b4d5/profile.go#L12
	forceColorProfile = flag.Int("force-profile", -1, "force a color profile (seems to only work in ssh mode)")
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
	if *forceColorProfile != -1 {
		lipgloss.DefaultRenderer().SetColorProfile(termenv.Profile(*forceColorProfile))
	}

	_, err := tea.NewProgram(NewModel("default"), programOpts...).Run()
	if err != nil {
		log.Error("Could not start program", "error", err)
	}
}

func startSSH() {
	addr := net.JoinHostPort(host, port)

	teaHandler := createTeaHandler()
	var teaMiddleware wish.Middleware

	if *forceColorProfile != -1 {
		teaMiddleware = bubbletea.MiddlewareWithColorProfile(
			createTeaHandler(),
			termenv.Profile(*forceColorProfile),
		)
	} else {
		teaMiddleware = bubbletea.Middleware(teaHandler)
	}

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			teaMiddleware,
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
	flag.Parse()

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

	if *enableSsh {
		startSSH()
		return
	}

	// logging messes with the TUI, so we have to write logs to a file when not running as a ssh server
	f, err := tea.LogToFileWith("program.log", "debug", log.Default())
	if err != nil {
		log.Fatal("could not open log file", "error", err)
	}
	defer f.Close()

	startProgram()
}
