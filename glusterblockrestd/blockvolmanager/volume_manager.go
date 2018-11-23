package blockvolmanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gluster/gluster-block-restapi/pkg/utils"
)

// BlockVolumeManager defines set of methods for various gluster block operation
type BlockVolumeManager interface {
	CreateBlockVolume(hostVolName string, blockName string, size uint64, hosts []string, optFuncs ...OptFunc) ([]byte, error)
	DeleteBlockVolume(hostVolName string, blockName string, optFuncs ...OptFunc) error
	BlockVolumeInfo(hostVolName string, blockName string) ([]byte, error)
	ListBlockVolume(hostVolName string) ([]byte, error)
}

type blockVolumeCLI struct {
	cliPath string
}

// NewBlockVolumeCLI returns a concrete instance implementing the BlockVolumeManager
// interface. If the `clipath`  param is empty then it will use the default gluster-block
// cli path
func NewBlockVolumeCLI(cliPath string) BlockVolumeManager {
	bm := &blockVolumeCLI{
		cliPath: cliPath,
	}

	if bm.cliPath != "" {
		return bm
	}

	if defaultCliPath, err := exec.LookPath("gluster-block"); err == nil {
		bm.cliPath = defaultCliPath
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
func (bm *blockVolumeCLI) CreateBlockVolume(hostVolName string, blockName string, size uint64, hosts []string, optFuncs ...OptFunc) ([]byte, error) {
	var (
		cmdArgs = []string{"create", hostVolName + "/" + blockName}
		opt     = &Options{}
	)

	opt.applyOpts(optFuncs...)
	optinalArgs := opt.prepareArgs()

	cmdArgs = append(cmdArgs, optinalArgs...)
	cmdArgs = append(cmdArgs, strings.Join(hosts, ","), strconv.FormatUint(size, 10), "--json")

	out, err := utils.ExecuteCommandOutput(bm.cliPath, cmdArgs...)

	if err != nil {
		return nil, err
	}

	return truncateCliOutput(out)
}

// DeleteBlockVolume deletes a block volume
// command to delete a block volume:
// delete  <volname/blockname> [unlink-storage <yes|no>] [force]
func (bm *blockVolumeCLI) DeleteBlockVolume(hostVolName string, blockName string, optFuncs ...OptFunc) error {
	var (
		cmdArgs = []string{"delete", hostVolName + "/" + blockName}
		opt     = &Options{}
	)

	opt.applyOpts(optFuncs...)
	optinalArgs := opt.prepareArgs()

	cmdArgs = append(cmdArgs, optinalArgs...)

	return utils.ExecuteCommandRun(bm.cliPath, cmdArgs...)
}

// BlockVolumeInfo returns details about block device.
func (bm *blockVolumeCLI) BlockVolumeInfo(hostVolName string, blockName string) ([]byte, error) {
	var cmdArgs = []string{"info", hostVolName + "/" + blockName, "--json"}
	return utils.ExecuteCommandOutput(bm.cliPath, cmdArgs...)
}

// ListBlockVolume will list available block devices for a given hosting volume
func (bm *blockVolumeCLI) ListBlockVolume(hostVolName string) ([]byte, error) {
	var cmdArgs = []string{"list", hostVolName, "--json"}
	return utils.ExecuteCommandOutput(bm.cliPath, cmdArgs...)
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
