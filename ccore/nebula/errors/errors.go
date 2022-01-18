package errors

import "errors"

var (
	ErrUnsupportedVersion = errors.New("unsupported version")
	ErrUnsupported        = errors.New("unsupported")
	ErrNoEndpoints        = errors.New("no endpoints")
	ErrNoJobStats         = errors.New("no job stats")
)
