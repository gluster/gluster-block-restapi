package errors

import (
	"errors"
)

var (
	// ErrJSONParsingFailed error is returned when JSON parsing is failed
	ErrJSONParsingFailed = errors.New("unable to parse the request")
)
