package northgen

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/loveyourstack/northwind/cmd"
	"github.com/loveyourstack/northwind/internal/nw"
	"github.com/spf13/cobra"
)

var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "northgen",
	Version: version,
	Short:   "northgen - code generation CLI tool",
	Long:    `northgen is a CLI tool for development use only. It creates Go and SQL code from PostgreSQL tables.`,
	// no Run function: a subcommand is always needed
}

type cliApplication struct {
	*cmd.Application
}

var cliApp *cliApplication

func init() {
	cobra.OnInitialize(initApp)
}

func initApp() {

	// load config from file
	conf := nw.Config{}
	err := conf.LoadFromFile("/usr/local/etc/nw_config.toml")
	if err != nil {
		log.Fatalf("initialization: nw_config.toml not found: %s" + err.Error())
	}

	ctx := context.Background()

	// create non-specific app
	app := &cmd.Application{
		Config:   &conf,
		InfoLog:  slog.New(slog.NewTextHandler(os.Stdout, nil)),
		ErrorLog: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	// create cli app
	cliApp = &cliApplication{app}

	// connect to db and assign pool to cliApp
	cliApp.Db, err = lyspgdb.GetPool(ctx, conf.Db, conf.DbOwnerUser)
	if err != nil {
		log.Fatalf("initialization: failed to create db connection pool: %s", err.Error())
	}
	// not deferring Db.Close() here: it is called before subcommand is reached. Defer close in subcommand instead
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}
