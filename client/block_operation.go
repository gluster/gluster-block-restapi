package client

import (
	"fmt"
	"net/http"

	"github.com/gluster/gluster-block-restapi/pkg/api"
)

// CreateBlockVolume creates a block volume
func (c *Client) CreateBlockVolume(hostvolume, blockname string, req *api.BlockVolumeCreateReq) (*api.BlockVolumeCreateResp, error) {
	var (
		resp    = &api.BlockVolumeCreateResp{}
		reqPath = fmt.Sprintf("/blockvolumes/%s/%s", hostvolume, blockname)
	)

	err := c.Verb(http.MethodPost).APIVersion("v1").Path(reqPath).Do(req, resp)
	return resp, err
}

// DeleteBlockVolume deletes a block volume
func (c *Client) DeleteBlockVolume(hostvolume, blockname string, req *api.BlockVolumeDeleteReq) error {
	var (
		reqPath = fmt.Sprintf("/blockvolumes/%s/%s", hostvolume, blockname)
	)
	return c.Verb(http.MethodDelete).APIVersion("v1").Path(reqPath).Do(req, nil)
}

// ListBlockVolumes will list all block volumes present in given hosting volume
func (c *Client) ListBlockVolumes(hostvolume string) (*api.BlockVolumeListResponse, error) {
	var (
		resp    = &api.BlockVolumeListResponse{}
		reqPath = fmt.Sprintf("/blockvolumes/%s", hostvolume)
	)

	err := c.Verb(http.MethodGet).APIVersion("v1").Path(reqPath).Do(nil, resp)
	return resp, err
}

// BlockVolumeInfo gives info about a given block volume
func (c *Client) BlockVolumeInfo(hostvolume, blockname string) (*api.BlockVolumeInfoResponse, error) {
	var (
		resp    = &api.BlockVolumeInfoResponse{}
		reqPath = fmt.Sprintf("/blockvolumes/%s/%s", hostvolume, blockname)
	)

	err := c.Verb(http.MethodGet).APIVersion("v1").Path(reqPath).Do(nil, resp)
	return resp, err
}
