package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lyserr"
	"github.com/loveyourstack/northwind/cmd"
)

type httpServerApplication struct {
	*cmd.Application
	GetOptions  lys.GetOptions
	PostOptions lys.PostOptions
}

type contextKey string

const supplierIdCtxKey contextKey = "SupplierId"

// authenticate is middleware that authenticates the user and adds his information to request context
func (srvApp *httpServerApplication) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx context.Context

		// get supplier ID from auth header
		supplierId, err := getSupplierIdFromReq(r)
		if err != nil {
			lys.HandleInternalError(r.Context(), fmt.Errorf("authenticate: getSupplierIdFromReq failed"), srvApp.ErrorLog, w)
			return
		}

		// add the dummy user info to the request context, then continue
		ctx = context.WithValue(r.Context(), supplierIdCtxKey, supplierId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// logRequest is middleware that logs requests in the app's info log
func (srvApp *httpServerApplication) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get supplier ID from request (assumes user is authed)
		supplierId, ok := r.Context().Value(supplierIdCtxKey).(int64)
		if !ok {
			lys.HandleInternalError(r.Context(), fmt.Errorf("logRequest: user not authenticated"), srvApp.ErrorLog, w)
			return
		}

		// get remote ip
		remoteHostIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			lys.HandleInternalError(r.Context(), fmt.Errorf("logRequest: net.SplitHostPort failed: %w", err), srvApp.ErrorLog, w)
		}

		srvApp.InfoLog.Info(fmt.Sprintf("%s - %v - %s %s %s", remoteHostIP, supplierId, r.Proto, r.Method, r.URL.RequestURI()))

		next.ServeHTTP(w, r)
	})
}

// secureHeaders is middleware that add XSS protection headers to responses
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// getSupplierIdFromReq returns the supplier ID from the request's Authorization header
func getSupplierIdFromReq(r *http.Request) (supplierId int64, err error) {

	authHeaderFound := false

	for name, vals := range r.Header {

		if name != "Authorization" {
			continue
		}

		authHeaderFound = true
		supplierIdStr := vals[0]
		supplierId, err = strconv.ParseInt(supplierIdStr, 10, 64)
		if err != nil {
			return 0, lyserr.User{Message: "supplier ID must be an integer", StatusCode: http.StatusForbidden}
		}
		break
	}

	if !authHeaderFound {
		return 0, lyserr.User{Message: "Authorization header missing from request", StatusCode: http.StatusForbidden}
	}
	if supplierId == 0 {
		return 0, lyserr.User{Message: "supplier ID missing from Authorization header", StatusCode: http.StatusForbidden}
	}

	return supplierId, nil
}
