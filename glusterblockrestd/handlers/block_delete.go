package handlers

import (
	"net/http"

	block "github.com/gluster/gluster-block-restapi/glusterblockrestd/blockvolmanager"
	"github.com/gluster/gluster-block-restapi/pkg/api"
	"github.com/gluster/gluster-block-restapi/pkg/errors"
	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/gorilla/mux"
)

func (gb *GlusterBlockHandler) deleteBlockVolume(w http.ResponseWriter, r *http.Request) {
	var (
		p          = mux.Vars(r)
		hostVolume = p["hostvolume"]
		blockName  = p["blockname"]
		opts       = []block.OptFunc{}
		req        = api.BlockVolumeDeleteReq{}
	)

	if err := utils.UnmarshalRequest(r, &req); err != nil {
		utils.SendHTTPError(w, http.StatusBadRequest, errors.ErrJSONParsingFailed)
		return
	}

	if req.UnlinkStorage {
		opts = append(opts, block.WithUnlinkStorage)
	}

	if req.Force {
		opts = append(opts, block.WithForceDelete)
	}

	if err := gb.blockVolManager.DeleteBlockVolume(hostVolume, blockName, opts...); err != nil {
		utils.SendHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendHTTPResponse(w, http.StatusNoContent, nil)
}
