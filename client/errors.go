package client

import (
	"encoding/json"
	"fmt"

	"github.com/gluster/gluster-block-restapi/pkg/api"
)

// ErrorResponse embeds server error response schema with returned http status code
// It implements `error` interface
type ErrorResponse struct {
	*api.ErrorResp
	HTTPStatusCode int `json:"http_status_code"`
}

// Error returns json representation of ErrorResponse
func (e *ErrorResponse) Error() string {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("http_status_code: %d, error_msg: %s, error_code: %d", e.HTTPStatusCode, e.ErrorResp.Error, e.Code)
	}
	return string(body)
}
