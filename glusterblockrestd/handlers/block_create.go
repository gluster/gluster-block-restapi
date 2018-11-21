package handlers

import (
	"encoding/json"
	stderrors "errors"
	"net/http"

	"github.com/gluster/gluster-block-restapi/glusterblockrestd/blockvolmanager"
	"github.com/gluster/gluster-block-restapi/pkg/api"
	"github.com/gluster/gluster-block-restapi/pkg/errors"
	"github.com/gluster/gluster-block-restapi/pkg/utils"
)

func blockVolumeCreateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req     = &api.BlockVolumeCreateReq{}
		resp    = &api.BlockVolumeCreateResp{}
		errResp = &api.CLIError{}
	)

	if err := utils.UnmarshalRequest(r, req); err != nil {
		utils.SendHTTPError(w, http.StatusBadRequest, errors.ErrJSONParsingFailed)
		return
	}

	blockManager := blockvolmanager.NewBlockVolumeCLI(blockvolmanager.WithCLIPath(glusterBlockCLI))
	body, err := blockManager.CreateBlockVolume(req)

	if err != nil {
		utils.SendHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(body, errResp); err == nil && errResp.ErrMsg != "" {
		utils.SendHTTPError(w, http.StatusInternalServerError, stderrors.New(errResp.ErrMsg))
		return
	}

	if err := json.Unmarshal(body, resp); err != nil {
		utils.SendHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendHTTPResponse(w, http.StatusCreated, body)

}

func init() {
	registerRoute(Route{
		"BlockVolumeCreate",
		"POST",
		"/v1/blockvolumes",
		blockVolumeCreateHandler,
	})
}
