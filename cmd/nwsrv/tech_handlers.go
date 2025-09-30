package main

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/loveyourstack/lys"
)

func (srvApp *httpServerApplication) techGetDatabaseSize(w http.ResponseWriter, r *http.Request) {

	stmt := "SELECT pg_size_pretty(pg_database_size(current_database()));"
	rows, _ := srvApp.Db.Query(r.Context(), stmt)
	sizePretty, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	if err != nil {
		lys.HandleError(r.Context(), fmt.Errorf("pgx.CollectExactlyOneRow failed: %w", err), srvApp.ErrorLog, w)
		return
	}

	// return success
	resp := lys.StdResponse{
		Status: lys.ReqSucceeded,
		Data:   sizePretty,
	}
	lys.JsonResponse(resp, http.StatusOK, w)
}

func (srvApp *httpServerApplication) techGetPgVersion(w http.ResponseWriter, r *http.Request) {

	stmt := "SELECT version();"
	rows, _ := srvApp.Db.Query(r.Context(), stmt)
	vers, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	if err != nil {
		lys.HandleError(r.Context(), fmt.Errorf("pgx.CollectExactlyOneRow failed: %w", err), srvApp.ErrorLog, w)
		return
	}

	// return success
	resp := lys.StdResponse{
		Status: lys.ReqSucceeded,
		Data:   vers,
	}
	lys.JsonResponse(resp, http.StatusOK, w)
}
