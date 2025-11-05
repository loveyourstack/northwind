package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lyspgmon/stores/lyspgbloat"
	"github.com/loveyourstack/lys/lyspgmon/stores/lyspgquery"
	"github.com/loveyourstack/lys/lyspgmon/stores/lyspgsetting"
	"github.com/loveyourstack/lys/lyspgmon/stores/lyspgtablesize"
	"github.com/loveyourstack/lys/lyspgmon/stores/lyspgunusedidx"
	"github.com/loveyourstack/lys/lysstring"
	"github.com/loveyourstack/northwind/internal/enums/sysrole"
	"github.com/loveyourstack/northwind/internal/stores/core/corecategory"
	"github.com/loveyourstack/northwind/internal/stores/core/corecountry"
	"github.com/loveyourstack/northwind/internal/stores/core/coreproduct"
	"github.com/loveyourstack/northwind/internal/stores/core/coresupplier"
	"github.com/loveyourstack/northwind/internal/stores/hr/hremployee"
	"github.com/loveyourstack/northwind/internal/stores/hr/hrmeetingsched"
	"github.com/loveyourstack/northwind/internal/stores/sales/salescustomer"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesorder"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesorderitem"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesordersbysalesman"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesshipper"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesterritory"
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
		{Url: "/ecb", RouteAdder: srvApp.ecbRoutes()},
		{Url: "/hr", RouteAdder: srvApp.hrRoutes(apiEnv)},
		{Url: "/sales", RouteAdder: srvApp.salesRoutes(apiEnv)},
		{Url: "/tech", RouteAdder: srvApp.techRoutes(apiEnv)},
	}
}

func (srvApp *httpServerApplication) coreRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		endpoint := "/categories"

		categoryStore := corecategory.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, categoryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, categoryStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, categoryStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, categoryStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, categoryStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, categoryStore)).Methods("DELETE")

		categoryStoreDU := corecategory.DuStore{Db: srvApp.Db}
		r.HandleFunc(endpoint+"-data-updates", lys.Get(apiEnv, categoryStoreDU)).Methods("GET")

		// test
		r.HandleFunc(endpoint+"/test", srvApp.coreAddCategoryTest).Methods("POST")

		endpoint = "/countries"

		countryStore := corecountry.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, countryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, countryStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, countryStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, countryStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, countryStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, countryStore)).Methods("DELETE")

		endpoint = "/products"

		productStore := coreproduct.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, productStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, productStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, productStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, productStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, productStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, productStore)).Methods("DELETE")

		r.HandleFunc(endpoint+"/distinct/countries", lys.GetSimple(apiEnv, productStore.DistinctSupplierCountries)).Methods("GET")

		endpoint = "/suppliers"

		supplierStore := coresupplier.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, supplierStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, supplierStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, supplierStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, supplierStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, supplierStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, supplierStore)).Methods("DELETE")

		supplierStoreDU := coresupplier.DuStore{Db: srvApp.Db}
		r.HandleFunc(endpoint+"-data-updates", lys.Get(apiEnv, supplierStoreDU)).Methods("GET")

		return r
	}
}

func (srvApp *httpServerApplication) ecbRoutes() lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		endpoint := "/sync-currencies"

		r.HandleFunc(endpoint, srvApp.ecbSyncCurrencies).Methods("POST")

		return r
	}
}

func (srvApp *httpServerApplication) hrRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		// restrict this subroute
		r.Use(lys.AuthorizeRole(lysstring.DeAlias(sysrole.RestrictedActions[:])))

		endpoint := "/employees"

		employeeStore := hremployee.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, employeeStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, employeeStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, employeeStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, employeeStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, employeeStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, employeeStore)).Methods("DELETE")

		endpoint = "/meeting-schedule"

		meetingSchedStore := hrmeetingsched.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, meetingSchedStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, meetingSchedStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, meetingSchedStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, meetingSchedStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, meetingSchedStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, meetingSchedStore)).Methods("DELETE")

		return r
	}
}

func (srvApp *httpServerApplication) salesRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {
		schemaName := "sales"

		endpoint := "/customers"

		customerStore := salescustomer.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, customerStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, customerStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, customerStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, customerStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, customerStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, customerStore)).Methods("DELETE")

		endpoint = "/orders"

		orderStore := salesorder.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, orderStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, orderStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, orderStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}/restore", lys.RestoreById(apiEnv, srvApp.Db, orderStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, orderStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, orderStore)).Methods("PATCH")
		//r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, orderStore)).Methods("DELETE")
		r.HandleFunc(endpoint+"/{id}/archive", lys.ArchiveById(apiEnv, srvApp.Db, orderStore)).Methods("DELETE")

		endpoint = "/orders-by-salesman"

		ordersBySalesmanStore := salesordersbysalesman.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, ordersBySalesmanStore, lys.GetOption[salesordersbysalesman.Model]{
			SetFuncUrlParamNames: ordersBySalesmanStore.GetSetFuncUrlParamNames(),
		})).Methods("GET")

		endpoint = "/order-value-latest-weeks"

		r.HandleFunc(endpoint, lys.GetSimple(apiEnv, orderStore.SelectOrderValueLatestWeeks)).Methods("GET")

		endpoint = "/order-items"

		orderItemStore := salesorderitem.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, orderItemStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, orderItemStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, orderItemStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}/restore", lys.RestoreById(apiEnv, srvApp.Db, orderItemStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, orderItemStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, orderItemStore)).Methods("PATCH")
		//r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, orderItemStore)).Methods("DELETE")
		r.HandleFunc(endpoint+"/{id}/archive", lys.ArchiveById(apiEnv, srvApp.Db, orderItemStore)).Methods("DELETE")

		r.HandleFunc("/regions", lys.GetEnumValues(apiEnv, srvApp.Db, schemaName, "region")).Methods("GET")

		endpoint = "/shippers"

		shipperStore := salesshipper.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, shipperStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, shipperStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, shipperStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, shipperStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, shipperStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, shipperStore)).Methods("DELETE")

		endpoint = "/territories"

		territoryStore := salesterritory.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get(apiEnv, territoryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById(apiEnv, territoryStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post(apiEnv, territoryStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put(apiEnv, territoryStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, territoryStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, territoryStore)).Methods("DELETE")

		return r
	}
}

func (srvApp *httpServerApplication) techRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		endpoint := "/long-running-query"

		r.HandleFunc(endpoint, lys.PgSleep(srvApp.Db, srvApp.ErrorLog, 120)).Methods("GET")

		endpoint = "/pg-bloat"

		bloatStore := lyspgbloat.Store{Db: srvApp.OwnerDb} // uses db owner
		r.HandleFunc(endpoint, lys.Get(apiEnv, bloatStore)).Methods("GET")

		endpoint = "/pg-database-size"

		r.HandleFunc(endpoint, srvApp.techGetDatabaseSize).Methods("GET")

		endpoint = "/pg-queries"

		activityStore := lyspgquery.Store{Db: srvApp.OwnerDb} // uses db owner
		r.HandleFunc(endpoint, lys.Get(apiEnv, activityStore)).Methods("GET")

		endpoint = "/pg-settings"

		settingStore := lyspgsetting.Store{Db: srvApp.OwnerDb} // uses db owner
		r.HandleFunc(endpoint, lys.Get(apiEnv, settingStore)).Methods("GET")

		endpoint = "/pg-table-size"

		tableSizeStore := lyspgtablesize.Store{Db: srvApp.OwnerDb} // uses db owner
		r.HandleFunc(endpoint, lys.Get(apiEnv, tableSizeStore)).Methods("GET")

		endpoint = "/pg-unused-indexes"

		unusedIdxStore := lyspgunusedidx.Store{Db: srvApp.OwnerDb} // uses db owner
		r.HandleFunc(endpoint, lys.Get(apiEnv, unusedIdxStore)).Methods("GET")

		endpoint = "/pg-version"

		r.HandleFunc(endpoint, srvApp.techGetPgVersion).Methods("GET")

		return r
	}
}
