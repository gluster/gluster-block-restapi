package handlers

import (
	"net/http"

	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/gorilla/mux"
)

func (gb *GlusterBlockHandler) blockVolumeInfo(w http.ResponseWriter, r *http.Request) {
	var (
		p           = mux.Vars(r)
		hostVolume  = p["hostvolume"]
		blockVolume = p["blockname"]
	)

	out, err := gb.blockVolManager.BlockVolumeInfo(hostVolume, blockVolume)
	if err != nil {
		utils.SendHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendHTTPResponse(w, http.StatusOK, out)
}
