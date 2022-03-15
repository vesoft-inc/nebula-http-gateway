package errors

import "errors"

var (
	ErrUnsupportedVersion  = errors.New("unsupported version")
	ErrUnsupported         = errors.New("unsupported")
	ErrNoEndpoints         = errors.New("no endpoints")
	ErrNoJobStats          = errors.New("no job stats")
	ErrUnknownMetaEndpoint = errors.New("unknown meta endpoint to update connection")
	ErrNoValidMetaEndpoint = errors.New("no valid meta endpoint to connect")
)
