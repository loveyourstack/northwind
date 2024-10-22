package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/loveyourstack/connectors/apiclients/ecbapi"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/northwind/cmd"
	"github.com/loveyourstack/northwind/internal/enums/sysrole"
)

type httpServerApplication struct {
	*cmd.Application
	GetOptions  lys.GetOptions
	PostOptions lys.PostOptions
	EcbClient   ecbapi.Client
}

// authenticate is middleware that authenticates the user and adds his information to request context
func (srvApp *httpServerApplication) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx context.Context

		// no auth: just assume dummy user
		reqUserInfo := lys.ReqUserInfo{
			Roles:    []string{sysrole.Tech.String()},
			UserId:   1,
			UserName: "Unauthed Dev",
		}

		// add the dummy user info to the request context, then continue
		ctx = context.WithValue(r.Context(), lys.ReqUserInfoCtxKey, reqUserInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// logRequest is middleware that logs requests in the app's info log
func (srvApp *httpServerApplication) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get user from request (assumes user is authed)
		userInfo, ok := r.Context().Value(lys.ReqUserInfoCtxKey).(lys.ReqUserInfo)
		if !ok {
			lys.HandleInternalError(r.Context(), fmt.Errorf("logRequest: user not authenticated"), srvApp.ErrorLog, w)
			return
		}

		// get remote ip
		remoteHostIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			lys.HandleInternalError(r.Context(), fmt.Errorf("logRequest: net.SplitHostPort failed: %w", err), srvApp.ErrorLog, w)
		}

		srvApp.InfoLog.Info(fmt.Sprintf("%s - %s - %s %s %s", remoteHostIP, userInfo.UserName, r.Proto, r.Method, r.URL.RequestURI()))

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
