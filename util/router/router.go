package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	StrictSlash bool
	Handler     http.Handler
}

func New(routes []Route) *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		router.
			StrictSlash(route.StrictSlash).
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)

		router.
			StrictSlash(route.StrictSlash).
			Methods(http.MethodOptions).
			Path(route.Pattern).
			Name(route.Name).
			Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
	}

	return router
}
