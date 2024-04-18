package cmd

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/northwind/internal/nw"
)

// Application contains the fields common to all commands
type Application struct {
	Config   *nw.Config
	InfoLog  *slog.Logger
	ErrorLog *slog.Logger
	Db       *pgxpool.Pool
	Validate *validator.Validate
}
