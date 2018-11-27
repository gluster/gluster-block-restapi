package api

// ErrorResp contains an error code and corresponding text which briefly
// describes the error in short.
type ErrorResp struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// BlockVolumeCreateReq represents gluster block
// volume create request
type BlockVolumeCreateReq struct {
	HaCount            int      `json:"ha"`
	AuthEnabled        bool     `json:"auth-enabled"`
	FullPrealloc       bool     `json:"full-prealloc"`
	Size               uint64   `json:"size"`
	Storage            string   `json:"storage"`
	RingBufferSizeInMB uint64   `json:"ring-buffer-size-mb"`
	Hosts              []string `json:"hosts"`
}

// BlockVolumeCreateResp represents  gluster block
// volume create resp
// Note: same as gluster-block cli output
type BlockVolumeCreateResp struct {
	IQN      string   `json:"IQN"`
	Portals  []string `json:"PORTAL(S)"`
	Username string   `json:"USERNAME"`
	Password string   `json:"PASSWORD"`
	Result   string   `json:"RESULT"`
}

// CLIError respresents error from gluster-block cli
type CLIError struct {
	Result  string `json:"RESULT"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

// BlockVolumeDeleteReq represents gluster block volume
// delete request
type BlockVolumeDeleteReq struct {
	UnlinkStorage bool `json:"unlink-storage"`
	Force         bool `json:"force"`
}

// BlockVolumeListResponse represents gluster block volume
// list resp
type BlockVolumeListResponse struct {
	Blocks []string `json:"blocks"`
	Result string   `json:"RESULT"`
}

// BlockVolumeInfoResponse represents gluster block volume
// info resp
type BlockVolumeInfoResponse struct {
	Name       string   `json:"NAME"`
	Volume     string   `json:"VOLUME"`
	GBID       string   `json:"GBID"`
	Size       string   `json:"SIZE"`
	Ha         int      `json:"HA"`
	Password   string   `json:"PASSWORD"`
	ExportedOn []string `json:"EXPORTED ON"`
}
