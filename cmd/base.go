package cmd

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/northwind/internal/nw"
)

// Application contains the fields common to all commands
type Application struct {
	Config   *nw.Config
	InfoLog  *slog.Logger
	ErrorLog *slog.Logger
	Db       *pgxpool.Pool // app-level connection for queries
	OwnerDb  *pgxpool.Pool // db owner connection for monitoring
	Validate *validator.Validate
}

// NewApplication returns an Application with default settings. Not all fields get initialized.
func NewApplication(conf *nw.Config) (app *Application) {

	// declare and configure logs
	var infoLog, errorLog *slog.Logger
	if conf.General.Debug {
		infoLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		errorLog = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		infoLog = slog.New(slog.NewTextHandler(os.Stdout, nil))
		errorLog = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	return &Application{
		Config:   conf,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}
