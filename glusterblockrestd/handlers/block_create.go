package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gluster/gluster-block-restapi/pkg/api"
	"github.com/gluster/gluster-block-restapi/pkg/errors"
	"github.com/gluster/gluster-block-restapi/pkg/utils"

	"github.com/gorilla/mux"
)

func blockVolumeCreateHandler(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	hostVolume := p["hostvolume"]
	blockName := p["blockname"]

	var req api.BlockVolumeCreateReq
	if err := utils.UnmarshalRequest(r, &req); err != nil {
		utils.SendHTTPError(w, http.StatusBadRequest, errors.ErrJSONParsingFailed)
		return
	}

	// create  <volname/blockname> [ha <count>]
	//                               [auth <enable|disable>]
	//                               [prealloc <full|no>]
	//                               [storage <filename>]
	//                               [ring-buffer <size-in-MB-units>]
	//                               <host1[,host2,...]> [size]
	//         create block device [defaults: ha 1, auth disable, prealloc no, size in bytes,
	// 	                     ring-buffer default size dependends on kernel]
	cmdArgs := []string{"create", hostVolume + "/" + blockName}
	if req.HaCount > 0 {
		cmdArgs = append(cmdArgs, "ha", strconv.Itoa(req.HaCount))
	}
	if req.FullPrealloc {
		cmdArgs = append(cmdArgs, "prealloc", "full")
	}
	if req.AuthEnabled {
		cmdArgs = append(cmdArgs, "auth", "enable")
	}
	if req.RingBufferSizeInMB > 0 {
		cmdArgs = append(cmdArgs, "ring-buffer", strconv.FormatUint(req.RingBufferSizeInMB, 10))
	}

	if req.Storage != "" {
		cmdArgs = append(cmdArgs, "storage", req.Storage)
	}

	cmdArgs = append(cmdArgs, strings.Join(req.Hosts, ","), strconv.FormatUint(req.Size, 10), "--json")

	out, err := utils.ExecuteCommandOutput(glusterBlockCLI, cmdArgs...)
	if err != nil {
		utils.SendHTTPError(w, http.StatusInternalServerError, err)
		return
	}

	utils.SendHTTPResponse(w, http.StatusCreated, out)
}

func init() {
	registerRoute(Route{
		"BlockVolumeCreate",
		"POST",
		"/v1/blockvolumes/{hostvolume}/{blockname}",
		blockVolumeCreateHandler,
	})
}
