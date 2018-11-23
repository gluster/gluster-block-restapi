package handlers

import (
	"encoding/json"
	stderrors "errors"
	"net/http"

	block "github.com/gluster/gluster-block-restapi/glusterblockrestd/blockvolmanager"
	"github.com/gluster/gluster-block-restapi/pkg/api"
	"github.com/gluster/gluster-block-restapi/pkg/errors"
	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/gorilla/mux"
)

func (gb *GlusterBlockHandler) createBlockVolume(w http.ResponseWriter, r *http.Request) {
	var (
		req        = &api.BlockVolumeCreateReq{}
		resp       = &api.BlockVolumeCreateResp{}
		errResp    = &api.CLIError{}
		pathParams = mux.Vars(r)
		opts       = []block.OptFunc{}
	)

	if err := utils.UnmarshalRequest(r, req); err != nil {
		utils.SendHTTPError(w, http.StatusBadRequest, errors.ErrJSONParsingFailed)
		return
	}

	opts = append(opts,
		block.WithHaCount(req.HaCount),
		block.WithRingBufferSizeInMB(req.RingBufferSizeInMB),
		block.WithStorage(req.Storage),
	)

	if req.AuthEnabled {
		opts = append(opts, block.WithAuthEnabled)
	}
	if req.FullPrealloc {
		opts = append(opts, block.WithFullPrealloc)
	}

	body, err := gb.blockVolManager.CreateBlockVolume(pathParams["hostvolume"], pathParams["blockname"], req.Size, req.Hosts, opts...)

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
