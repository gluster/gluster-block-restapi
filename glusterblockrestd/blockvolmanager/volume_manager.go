package blockvolmanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gluster/gluster-block-restapi/pkg/api"
	"github.com/gluster/gluster-block-restapi/pkg/utils"
)

// BlockVolumeManager defines set of methods for various gluster block operation
type BlockVolumeManager interface {
	CreateBlockVolume(req *api.BlockVolumeCreateReq) ([]byte, error)
}

// BlockVolumeCliOptFuncs receives a blockVolumeCLI object and overrides its members
type BlockVolumeCliOptFuncs func(cli *blockVolumeCLI)

type blockVolumeCLI struct {
	cliPath string
}

// WithCLIPath overrides cliPath of gluster-block
func WithCLIPath(path string) BlockVolumeCliOptFuncs {
	return func(cli *blockVolumeCLI) {
		cli.cliPath = path
	}
}

// NewBlockVolumeCLI returns a concrete instance implementing the BlockVolumeManager
// interface.
func NewBlockVolumeCLI(optFuncs ...BlockVolumeCliOptFuncs) BlockVolumeManager {
	bm := &blockVolumeCLI{}

	if defaultCliPath, err := exec.LookPath("gluster-block"); err == nil {
		bm.cliPath = defaultCliPath
	}

	for _, optFunc := range optFuncs {
		optFunc(bm)
	}
	return bm
}

// CreateBlockVolume creates a gluster block volume using gluster-block cli
// command to create a block volume:
// create  <volname/blockname>   [ha <count>]
//                               [auth <enable|disable>]
//                               [prealloc <full|no>]
//                               [storage <filename>]
//                               [ring-buffer <size-in-MB-units>]
//                               <host1[,host2,...]> [size]
// create block device [defaults: ha 1, auth disable, prealloc no, size in bytes,
// ring-buffer default size dependends on kernel]
func (bm *blockVolumeCLI) CreateBlockVolume(req *api.BlockVolumeCreateReq) ([]byte, error) {
	cmdArgs := []string{"create", req.HostingVolume + "/" + req.Name}

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

	out, err := utils.ExecuteCommandOutput(bm.cliPath, cmdArgs...)

	if err != nil {
		return nil, err
	}

	return truncateCliOutput(out)
}

// truncateCliOutput will check if the given `body` is a valid
// json.If it is a valid json return as it is.
// In case body contains some extra msg from cli along with a
// json body, then it will truncate the extra msg and returns
// the json body only
func truncateCliOutput(body []byte) ([]byte, error) {
	// check if a given body is a valid json. If it is a valid
	// json return the body as it is.
	if err := json.Unmarshal(body, &json.RawMessage{}); err == nil {
		return body, nil
	}

	// if given body contains extra msg along with json body ,
	// e.g. if the cli output is like
	// The size 10000000 will align to sector size 512 bytes
	//{ "IQN": "iqn.2016-12.org.gluster-block:45de1338-52f8-4a92-b82e-b4f0c1429299",
	// "PORTAL(S)": [ "192.168.122.16:3260", "192.168.122.131:3260", "192.168.122.247:3260" ], "RESULT": "SUCCESS" }

	// In this case we need to truncate the extra msg (The size 10000000 will align to sector size 512 bytes)
	// from cli output before sending the resp to client

	index := bytes.IndexRune(body, '{')
	if index == -1 {
		return nil, errors.New("invalid json body")
	}

	// truncate the extra msg
	body = body[index:]

	if err := json.Unmarshal(body, &json.RawMessage{}); err != nil {
		return nil, err
	}

	return body, nil

}
