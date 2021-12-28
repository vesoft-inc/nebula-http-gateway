package errors

import "errors"

var (
	ErrUnsupportedVersion    = errors.New("unsupported version")
	ErrUnsupported           = errors.New("unsupported")
	ErrNoEndpoints           = errors.New("no endpoints")
	ErrVersionEstimateFailed = errors.New("version infer failed")
)
