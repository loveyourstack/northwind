package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/northwind/internal/stores/core/corecategory"
	"github.com/loveyourstack/northwind/internal/stores/core/coreproduct"
	"github.com/loveyourstack/northwind/internal/stores/core/coresupplier"
)

// getRouter returns a mux providing the HTTP server's routes
func (srvApp *httpServerApplication) getRouter() http.Handler {

	// define env struct needed for route handlers
	apiEnv := lys.Env{
		ErrorLog:    srvApp.ErrorLog,
		Validate:    srvApp.Validate,
		GetOptions:  srvApp.GetOptions,
		PostOptions: srvApp.PostOptions,
	}

	// increase default max # results from GET
	apiEnv.GetOptions.MaxPerPage = 5000

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(lys.NotFound())

	// public routes
	r.HandleFunc("/", lys.Message("Welcome to the "+srvApp.Config.General.AppName+" API. Please log in.")).Methods("GET")

	// put all routes requiring auth behind "/a" for authed
	authedR := r.PathPrefix("/a").Subrouter()

	// define middleware for authed routes
	authedR.Use(srvApp.authenticate)
	authedR.Use(srvApp.logRequest) // must come after authentication
	authedR.Use(secureHeaders)

	// add subroutes into main router
	for _, subRoute := range srvApp.getSubRoutes(apiEnv) {
		subRouter := authedR.PathPrefix(subRoute.Url).Subrouter()
		_ = subRoute.RouteAdder(subRouter)
	}

	// apply CORS middleware to allow access to Vue app
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{srvApp.Config.UI.Url}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Encoding", "Authorization", "Content-Length", "Content-Type", "X-CSRF-Token"}),
		handlers.ExposedHeaders([]string{"Content-Disposition"}),
		handlers.AllowCredentials(),
	)
	return (cors)(r)
}

// getSubRoutes returns all subroutes used by the server
func (srvApp *httpServerApplication) getSubRoutes(apiEnv lys.Env) []lys.SubRoute {

	return []lys.SubRoute{
		{Url: "/core", RouteAdder: srvApp.coreRoutes(apiEnv)},
	}
}

func (srvApp *httpServerApplication) coreRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		endpoint := "/categories"

		categoryStore := corecategory.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, categoryStore)).Methods("GET")

		endpoint = "/products"

		productStore := coreproduct.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, productStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, productStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, productStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, productStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, productStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, productStore)).Methods("DELETE")

		endpoint = "/suppliers"

		supplierStore := coresupplier.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, supplierStore)).Methods("GET") // for dummy login dropdown

		return r
	}
}
