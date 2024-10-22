package main

import (
	"fmt"
	"net/http"

	"github.com/loveyourstack/connectors/csyncdb"
	"github.com/loveyourstack/lys"
)

func (srvApp *httpServerApplication) ecbSyncCurrencies(w http.ResponseWriter, r *http.Request) {

	err := csyncdb.EcbCurrencies(r.Context(), srvApp.Db, srvApp.EcbClient)
	if err != nil {
		lys.HandleError(r.Context(), fmt.Errorf("csyncdb.EcbCurrencies failed: %w", err), srvApp.ErrorLog, w)
		return
	}

	// return success
	resp := lys.StdResponse{
		Status: lys.ReqSucceeded,
		Data:   "synced",
	}
	lys.JsonResponse(resp, http.StatusOK, w)
}
