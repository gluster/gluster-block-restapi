package client

import "github.com/gluster/gluster-block-restapi/pkg/api"

// GlusterBlockClient is an interface for all block operations method
type GlusterBlockClient interface {
	CreateBlockVolume(hostvolume string, blockname string, req *api.BlockVolumeCreateReq) (*api.BlockVolumeCreateResp, error)
	DeleteBlockVolume(hostvolume string, blockname string, req *api.BlockVolumeDeleteReq) error
	ListBlockVolumes(hostvolume string) (*api.BlockVolumeListResponse, error)
	BlockVolumeInfo(hostvolume string, blockname string) (*api.BlockVolumeInfoResponse, error)
}
