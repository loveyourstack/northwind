package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/loveyourstack/connectors/apiclients/ecbapi"
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
		log.Fatalf("initialization: nw_config.toml not found: %s", err.Error())
	}

	ctx := context.Background()

	// create non-specific app
	app := cmd.NewApplication(&conf)

	// create http server app
	srvApp := &httpServerApplication{
		Application: app,
		GetOptions:  lys.FillGetOptions(lys.GetOptions{}),   // use defaults
		PostOptions: lys.FillPostOptions(lys.PostOptions{}), // use defaults
		EcbClient:   ecbapi.NewClient(app.InfoLog, app.ErrorLog),
	}

	// connect to db and assign pool to srvApp
	srvApp.Db, err = lyspgdb.GetPool(ctx, conf.Db, conf.DbServerUser, srvApp.Config.General.AppName+" Srv")
	if err != nil {
		log.Fatalf("initialization: failed to create regular db connection pool: %s", err.Error())
	}
	defer srvApp.Db.Close()

	// connect to db using db owner and assign to srvApp
	srvApp.OwnerDb, err = lyspgdb.GetPool(ctx, conf.Db, conf.DbOwnerUser, conf.General.AppName+" Srv")
	if err != nil {
		log.Fatalf("initialization: failed to create owner db connection pool: %s", err.Error())
	}
	defer srvApp.OwnerDb.Close()

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
