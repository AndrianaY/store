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
	}
	makeSwaggerHandler(router)
	return router
}

func makeSwaggerHandler(r *mux.Router) {
	const docsPath = "/docs"

	r.StrictSlash(false).Path(docsPath).Handler(http.RedirectHandler(docsPath+"/", http.StatusMovedPermanently))
	r.StrictSlash(true).PathPrefix(docsPath + "/").Handler(
		http.StripPrefix(docsPath+"/", http.FileServer(http.Dir("./swagger/"))),
	)

	r.Path("/api-docs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./openapi.yml")
	})
}
