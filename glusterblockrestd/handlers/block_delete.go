package handlers

import (
	"net/http"

	"github.com/gluster/gluster-block-restapi/pkg/api"
	"github.com/gluster/gluster-block-restapi/pkg/errors"
	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/gorilla/mux"
)

func (gb *GlusterBlockHandler) deleteBlockVolume(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	hostVolume := p["hostvolume"]
	blockName := p["blockname"]

	var req api.BlockVolumeDeleteReq
	if err := utils.UnmarshalRequest(r, &req); err != nil {
		utils.SendHTTPError(w, http.StatusBadRequest, errors.ErrJSONParsingFailed)
		return
	}

	//delete  <volname/blockname> [unlink-storage <yes|no>] [force]
	cmdArgs := []string{"delete", hostVolume + "/" + blockName}
	if req.UnlinkStorage {
		cmdArgs = append(cmdArgs, "unlink-storage", "yes")
	}
	if req.Force {
		cmdArgs = append(cmdArgs, "force")
	}

	err := utils.ExecuteCommandRun(glusterBlockCLI, cmdArgs...)
	if err != nil {
		utils.SendHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendHTTPResponse(w, http.StatusNoContent, nil)
}
