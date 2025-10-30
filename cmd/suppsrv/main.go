package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lyspgdb"
	"github.com/loveyourstack/northwind/cmd"
	"github.com/loveyourstack/northwind/internal/nw"
)

func main() {

	configFileName := "suppsrv_config.toml"

	// mandatory flag if not using default
	configFilePath := flag.String("configFilePath", configFileName, "Path to the config file")

	flag.Parse()

	// load config from file
	conf := nw.Config{}
	err := conf.LoadFromFile(*configFilePath)
	if err != nil {
		log.Fatalf("initialization: %s not found: %s", configFileName, err.Error())
	}

	ctx := context.Background()

	// create non-specific app
	app := cmd.NewApplication(&conf)

	// create http server app
	srvApp := &httpServerApplication{
		Application: app,
		GetOptions:  lys.FillGetOptions(lys.GetOptions{}),   // use defaults
		PostOptions: lys.FillPostOptions(lys.PostOptions{}), // use defaults
	}

	// connect to db and assign pool to srvApp
	srvApp.Db, err = getPoolUsingSupplierId(ctx, conf.Db, conf.DbServerUser, srvApp.Config.General.AppName+" Srv", srvApp.ErrorLog)
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
	startupMsg := fmt.Sprintf("starting suppsrv on port: %s", srvApp.Config.API.Port)
	if conf.General.Debug {
		startupMsg += ", debug: true"
	}
	srvApp.InfoLog.Info(startupMsg)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("initialization: srv.ListenAndServe failed: %s", err.Error())
	}
}

// adapted from https://github.com/jackc/pgx/issues/288
func getPoolUsingSupplierId(ctx context.Context, dbConfig lyspgdb.Database, userConfig lyspgdb.User, appName string, errorLog *slog.Logger) (db *pgxpool.Pool, err error) {

	cfg, err := lyspgdb.GetConfig(dbConfig, userConfig, appName)
	if err != nil {
		return nil, fmt.Errorf("lyspgdb.GetConfig failed: %w", err)
	}

	cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {

		// get supplierId from context (was set in authenticate() middleware)
		supplierId, ok := ctx.Value(supplierIdCtxKey).(int64)
		if !ok {
			errorLog.Error("supplierID not found in ctx")
			return false
		}

		// set supplierId into this connection's setting
		_, err := conn.Exec(ctx, "SELECT set_config('app.supplier_id', $1, false)", fmt.Sprintf("%v", supplierId))
		if err != nil {
			errorLog.Error("conn.Exec (set supplierId) failed: " + err.Error())
			return false
		}

		return true
	}

	cfg.AfterRelease = func(conn *pgx.Conn) bool {

		// unset the supplierId setting before this connection is released to pool
		_, err := conn.Exec(ctx, "SELECT set_config('app.supplier_id', '', false)")
		if err != nil {
			errorLog.Error("conn.Exec (unset supplierId) failed: " + err.Error())
			return false
		}

		return true
	}

	db, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig failed: %w", err)
	}

	return db, nil
}
