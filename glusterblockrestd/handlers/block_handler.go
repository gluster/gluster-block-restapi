package handlers

import (
	"net/http"

	"github.com/gluster/gluster-block-restapi/glusterblockrestd/blockvolmanager"

	"github.com/gorilla/mux"
)

// GlusterBlockHandler defines http Handlers for all gluster-block rest api.
type GlusterBlockHandler struct {
	blockVolManager blockvolmanager.BlockVolumeManager
}

// NewGlusterBlockHandler returns a GlusterBlockHandler
func NewGlusterBlockHandler(manager blockvolmanager.BlockVolumeManager) *GlusterBlockHandler {
	return &GlusterBlockHandler{manager}
}

// getAllRoutes returns all http Routes which are registered with GlusterBlockHandler
func (gb *GlusterBlockHandler) getAllRoutes() Routes {
	return Routes{
		{
			Name:        "BlockVolumeCreate",
			Method:      http.MethodPost,
			Pattern:     "/v1/blockvolumes/{hostvolume}/{blockname}",
			HandlerFunc: gb.createBlockVolume,
		},
		{
			Name:        "BlockVolumeDelete",
			Method:      http.MethodDelete,
			Pattern:     "/v1/blockvolumes/{hostvolume}/{blockname}",
			HandlerFunc: gb.deleteBlockVolume,
		},
		{
			Name:        "BlockVolumesInfo",
			Method:      http.MethodGet,
			Pattern:     "/v1/blockvolumes/{hostvolume}/{blockname}",
			HandlerFunc: gb.blockVolumeInfo,
		},
		{
			Name:        "BlockVolumesList",
			Method:      http.MethodGet,
			Pattern:     "/v1/blockvolumes/{hostvolume}",
			HandlerFunc: gb.listBlockVolumes,
		},
	}
}

// RegisterRoutes will register all gluster-block http route with provided Router.
func (gb *GlusterBlockHandler) RegisterRoutes(mux *mux.Router) {
	routes := gb.getAllRoutes()
	for _, route := range routes {
		mux.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}
