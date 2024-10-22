package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/northwind/internal/stores/core/corecategory"
)

func (srvApp *httpServerApplication) coreAddCategoryTest(w http.ResponseWriter, r *http.Request) {

	catStore := corecategory.Store{Db: srvApp.Db}
	input := corecategory.Input{
		Name:        "name2",
		Description: "desc2",
	}
	_, err := catStore.Insert(context.Background(), input)
	if err != nil {
		lys.HandleError(r.Context(), fmt.Errorf("catStore.Insert failed: %w", err), srvApp.ErrorLog, w)
		return
	}

	resp := lys.StdResponse{
		Status: lys.ReqSucceeded,
		Data:   "success",
	}
	lys.JsonResponse(resp, http.StatusOK, w)
}
