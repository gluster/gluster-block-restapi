package utils

import (
	"encoding/json"
	"net/http"

	"github.com/aravindavk/gluster-block-restapi/pkg/api"
	log "github.com/sirupsen/logrus"
)

// UnmarshalRequest unmarshals JSON in `r` into `v`
func UnmarshalRequest(r *http.Request, v interface{}) error {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.WithError(err).WithField("req", v).
				Error("Failed to close the http request body")
		}
	}()
	return json.NewDecoder(r.Body).Decode(v)
}

// SendHTTPResponse sends non-error response to the client.
func SendHTTPResponse(w http.ResponseWriter, statusCode int, resp []byte) {

	if resp != nil {
		// Do not include content-type header for responses such as 204
		// which as per RFC, should not have a response body.
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

	w.WriteHeader(statusCode)

	if resp != nil {
		if _, err := w.Write(resp); err != nil {
			log.WithError(err).Error("Failed to send the response -", resp)
		}
	}
}

// SendHTTPError sends an error response to the client.
func SendHTTPError(w http.ResponseWriter, statusCode int, apierr error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	errMsg := ""
	errCode := -1
	if apierr != nil {
		if v, ok := apierr.(*ExecuteCommandError); ok {
			errCode = v.ExitStatus
		}
		errMsg = apierr.Error()
	}

	resp := api.ErrorResp{
		Code:  errCode,
		Error: errMsg,
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Failed to send the response -", resp)
	}
}
