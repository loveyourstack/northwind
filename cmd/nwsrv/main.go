package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/loveyourstack/northwind/cmd"
	"github.com/loveyourstack/northwind/internal/nw"
)

func main() {

	// mandatory flag if not using default
	configFilePath := flag.String("configFilePath", "nw_config.toml", "Path to the config file")

	flag.Parse()

	// load config from file
	conf := nw.Config{}
	err := conf.LoadFromFile(*configFilePath)
	if err != nil {
		log.Fatalf("initialization: nw_config.toml not found: %s" + err.Error())
	}

	ctx := context.Background()

	// declare and configure logs
	var infoLog, errorLog *slog.Logger
	if conf.General.Debug {
		infoLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		errorLog = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		infoLog = slog.New(slog.NewTextHandler(os.Stdout, nil))
		errorLog = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	// create non-specific app
	app := &cmd.Application{
		Config:   &conf,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	// create http server app
	srvApp := &httpServerApplication{
		Application: app,
		GetOptions:  lys.FillGetOptions(lys.GetOptions{}),   // use defaults
		PostOptions: lys.FillPostOptions(lys.PostOptions{}), // use defaults
	}

	// connect to db and assign pool to srvApp
	srvApp.Db, err = lyspgdb.GetPool(ctx, conf.Db, conf.DbServerUser)
	if err != nil {
		log.Fatalf("initialization: failed to create db connection pool: %s", err.Error())
	}
	defer srvApp.Db.Close()

	// create HTTP server using srvApp's routes and handlers
	srv := &http.Server{
		Addr:    ":" + srvApp.Config.API.Port,
		Handler: srvApp.getRouter(),
	}

	// start the HTTP server process
	startupMsg := fmt.Sprintf("starting nwsrv on port: %s", srvApp.Config.API.Port)
	if conf.General.Debug {
		startupMsg += ", debug: true"
	}
	srvApp.InfoLog.Info(startupMsg)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("initialization: srv.ListenAndServe failed: %s", err.Error())
	}
}
