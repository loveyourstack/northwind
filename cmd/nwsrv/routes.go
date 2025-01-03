package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/loveyourstack/lys"
	"github.com/loveyourstack/lys/lysstring"
	"github.com/loveyourstack/northwind/internal/enums/sysrole"
	"github.com/loveyourstack/northwind/internal/stores/common/commoncountry"
	"github.com/loveyourstack/northwind/internal/stores/core/corecategory"
	"github.com/loveyourstack/northwind/internal/stores/core/coreproduct"
	"github.com/loveyourstack/northwind/internal/stores/core/coresupplier"
	"github.com/loveyourstack/northwind/internal/stores/hr/hremployee"
	"github.com/loveyourstack/northwind/internal/stores/sales/salescustomer"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesorder"
	"github.com/loveyourstack/northwind/internal/stores/sales/salesorderitem"
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
func (srvApp *httpServerApplication) getSubRoutes(apiEnv lys.Env) (subRoutes []lys.SubRoute) {

	subRoutes = append(subRoutes, lys.SubRoute{Url: "/common", RouteAdder: srvApp.commonRoutes(apiEnv)})
	subRoutes = append(subRoutes, lys.SubRoute{Url: "/core", RouteAdder: srvApp.coreRoutes(apiEnv)})
	subRoutes = append(subRoutes, lys.SubRoute{Url: "/ecb", RouteAdder: srvApp.ecbRoutes()})
	subRoutes = append(subRoutes, lys.SubRoute{Url: "/hr", RouteAdder: srvApp.hrRoutes(apiEnv)})
	subRoutes = append(subRoutes, lys.SubRoute{Url: "/sales", RouteAdder: srvApp.salesRoutes(apiEnv)})

	return subRoutes
}

func (srvApp *httpServerApplication) commonRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		endpoint := "/countries"

		countryStore := commoncountry.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[commoncountry.Model](apiEnv, countryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[commoncountry.Model](apiEnv, countryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, countryStore)).Methods("PATCH")

		return r
	}
}

func (srvApp *httpServerApplication) coreRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {

		endpoint := "/categories"

		categoryStore := corecategory.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[corecategory.Model](apiEnv, categoryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[corecategory.Model](apiEnv, categoryStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[corecategory.Input, int64](apiEnv, categoryStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[corecategory.Input](apiEnv, categoryStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, categoryStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, categoryStore)).Methods("DELETE")

		categoryStoreDU := corecategory.DuStore{Db: srvApp.Db}
		r.HandleFunc(endpoint+"-data-updates", lys.Get[corecategory.DuModel](apiEnv, categoryStoreDU)).Methods("GET")

		// test
		r.HandleFunc(endpoint+"/test", srvApp.coreAddCategoryTest).Methods("POST")

		endpoint = "/products"

		productStore := coreproduct.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[coreproduct.Model](apiEnv, productStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[coreproduct.Model](apiEnv, productStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[coreproduct.Input, int64](apiEnv, productStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[coreproduct.Input](apiEnv, productStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, productStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, productStore)).Methods("DELETE")

		r.HandleFunc(endpoint+"/distinct/countries", lys.GetSimple(apiEnv, productStore.DistinctSupplierCommonCountries)).Methods("GET")

		endpoint = "/suppliers"

		supplierStore := coresupplier.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[coresupplier.Model](apiEnv, supplierStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[coresupplier.Model](apiEnv, supplierStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[coresupplier.Input, int64](apiEnv, supplierStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[coresupplier.Input](apiEnv, supplierStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, supplierStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, supplierStore)).Methods("DELETE")

		supplierStoreDU := coresupplier.DuStore{Db: srvApp.Db}
		r.HandleFunc(endpoint+"-data-updates", lys.Get[coresupplier.DuModel](apiEnv, supplierStoreDU)).Methods("GET")

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
		r.HandleFunc(endpoint, lys.Get[hremployee.Model](apiEnv, employeeStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[hremployee.Model](apiEnv, employeeStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[hremployee.Input, int64](apiEnv, employeeStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[hremployee.Input](apiEnv, employeeStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, employeeStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, employeeStore)).Methods("DELETE")

		return r
	}
}

func (srvApp *httpServerApplication) salesRoutes(apiEnv lys.Env) lys.RouteAdderFunc {

	return func(r *mux.Router) *mux.Router {
		schemaName := "sales"

		endpoint := "/customers"

		customerStore := salescustomer.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[salescustomer.Model](apiEnv, customerStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[salescustomer.Model](apiEnv, customerStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[salescustomer.Input, int64](apiEnv, customerStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[salescustomer.Input](apiEnv, customerStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, customerStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, customerStore)).Methods("DELETE")

		endpoint = "/orders"

		orderStore := salesorder.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[salesorder.Model](apiEnv, orderStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[salesorder.Model](apiEnv, orderStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[salesorder.Input, int64](apiEnv, orderStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}/restore", lys.RestoreById(apiEnv, srvApp.Db, orderStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[salesorder.Input](apiEnv, orderStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, orderStore)).Methods("PATCH")
		//r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, orderStore)).Methods("DELETE")
		r.HandleFunc(endpoint+"/{id}/archive", lys.ArchiveById(apiEnv, srvApp.Db, orderStore)).Methods("DELETE")

		endpoint = "/order-value-latest-weeks"

		r.HandleFunc(endpoint, lys.GetSimple[salesorder.OrderValueLatestWeeksModel](apiEnv, orderStore.SelectOrderValueLatestWeeks)).Methods("GET")

		endpoint = "/order-items"

		orderItemStore := salesorderitem.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[salesorderitem.Model](apiEnv, orderItemStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[salesorderitem.Model](apiEnv, orderItemStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[salesorderitem.Input, int64](apiEnv, orderItemStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}/restore", lys.RestoreById(apiEnv, srvApp.Db, orderItemStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[salesorderitem.Input](apiEnv, orderItemStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, orderItemStore)).Methods("PATCH")
		//r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, orderItemStore)).Methods("DELETE")
		r.HandleFunc(endpoint+"/{id}/archive", lys.ArchiveById(apiEnv, srvApp.Db, orderItemStore)).Methods("DELETE")

		r.HandleFunc("/regions", lys.GetEnumValues(apiEnv, srvApp.Db, schemaName, "region")).Methods("GET")

		endpoint = "/shippers"

		shipperStore := salesshipper.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[salesshipper.Model](apiEnv, shipperStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[salesshipper.Model](apiEnv, shipperStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[salesshipper.Input, int64](apiEnv, shipperStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[salesshipper.Input](apiEnv, shipperStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, shipperStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, shipperStore)).Methods("DELETE")

		endpoint = "/territories"

		territoryStore := salesterritory.Store{Db: srvApp.Db}
		r.HandleFunc(endpoint, lys.Get[salesterritory.Model](apiEnv, territoryStore)).Methods("GET")
		r.HandleFunc(endpoint+"/{id}", lys.GetById[salesterritory.Model](apiEnv, territoryStore)).Methods("GET")
		r.HandleFunc(endpoint, lys.Post[salesterritory.Input, int64](apiEnv, territoryStore)).Methods("POST")
		r.HandleFunc(endpoint+"/{id}", lys.Put[salesterritory.Input](apiEnv, territoryStore)).Methods("PUT")
		r.HandleFunc(endpoint+"/{id}", lys.Patch(apiEnv, territoryStore)).Methods("PATCH")
		r.HandleFunc(endpoint+"/{id}", lys.Delete(apiEnv, territoryStore)).Methods("DELETE")

		return r
	}
}
