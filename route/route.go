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

	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware, validateMiddleware)
	for _, route := range routesForApp {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routesForApp = routes{
	route{
		"Index",
		"GET",
		"/",
		Index,
	},
	route{
		"Add",
		"POST",
		"/ip",
		IPAdd,
	},
	route{
		"ShowAll",
		"GET",
		"/ip",
		ShowAllIPs,
	},
	route{
		"Show",
		"GET",
		"/ip/{ip}",
		IPShow,
	},
	route{
		"Delete",
		"DELETE",
		"/ip/{ip}",
		IPDelete,
	},
}
