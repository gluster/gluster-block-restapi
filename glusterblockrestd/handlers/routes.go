package handlers

import (
	"net/http"

	"github.com/gluster/gluster-block-restapi/glusterblockrestd/blockvolmanager"
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

// NewRouter returns a router which dispatches all incoming http Request to a matched handler.
func NewRouter() *mux.Router {
	var (
		router              = mux.NewRouter().StrictSlash(true)
		blockManager        = blockvolmanager.NewBlockVolumeCLI(glusterBlockCLI)
		glusterBlockHandler = NewGlusterBlockHandler(blockManager)
	)
	// register all routes with a given router instance
	glusterBlockHandler.RegisterRoutes(router)
	return router
}
