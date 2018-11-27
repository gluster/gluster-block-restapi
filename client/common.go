package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const defaultClientTimeout = time.Second * 5

func defaultHTTPClient() *http.Client {
	roundTripper := http.DefaultTransport

	if tr, ok := roundTripper.(*http.Transport); ok {
		tr.DisableCompression = true
	}
	return &http.Client{
		Transport: roundTripper,
		Timeout:   defaultClientTimeout,
	}
}

func closeRespBody(response *http.Response) {
	var errs []error

	if response.Body == nil {
		return
	}

	if _, err := io.Copy(ioutil.Discard, response.Body); err != nil {
		errs = append(errs, err)
	}

	if err := response.Body.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return
	}

	for _, err := range errs {
		log.WithError(err).Error("error in closing response body")
	}
}

func checkResponseError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errResp = &ErrorResponse{}
	if err := json.Unmarshal(body, errResp); err != nil {
		return fmt.Errorf("request for url: %s returned status code: %s, body: %s", resp.Request.URL, resp.Status, string(body))
	}

	errResp.HTTPStatusCode = resp.StatusCode
	return errResp
}
