package v2_5

import (
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5"
)

func codeErrorIfHappened(code nthrift.ErrorCode, msg []byte) error {
	if code == nthrift.ErrorCode_SUCCEEDED {
		return nil
	}
	// TODO: Align with the code of nerrors
	return nerrors.NewCodeError(nerrors.ErrorCode(code), string(msg))
}
