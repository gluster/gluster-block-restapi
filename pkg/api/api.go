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

// BlockVolumeDeleteReq represents gluster block volume
// delete request
type BlockVolumeDeleteReq struct {
	UnlinkStorage bool `json:"unlink-storage"`
	Force         bool `json:"force"`
}
