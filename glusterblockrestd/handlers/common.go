package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

var glusterBlockCLI = "gluster-block"

// SetGlusterBlockCLI sets the gluster-block CLI path
func SetGlusterBlockCLI(cli string) {
	glusterBlockCLI = cli
}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines the list of routes of our API
type Routes []Route

var routes Routes

func registerRoute(route Route) {
	routes = append(routes, route)
}

// NewRoutes returns all the routes related to gluster block Volumes
func NewRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
