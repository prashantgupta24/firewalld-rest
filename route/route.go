package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

//NewRouter creates a new mux router for application
func NewRouter() *mux.Router {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/m{[0-9]+}").Subrouter().StrictSlash(true)

	router.Use(loggingMiddleware, validateMiddleware)
	for _, route := range routesForApp {
		subrouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return subrouter
}

var routesForApp = routes{
	route{
		"Index Page",
		"GET",
		"/",
		Index,
	},
	route{
		"Add New IP",
		"POST",
		"/ip",
		IPAdd,
	},
	route{
		"Show all IPs present",
		"GET",
		"/ip",
		ShowAllIPs,
	},
	route{
		"Show if particular IP is present",
		"GET",
		"/ip/{ip}",
		IPShow,
	},
	route{
		"Delete IP",
		"DELETE",
		"/ip/{ip}",
		IPDelete,
	},
}
