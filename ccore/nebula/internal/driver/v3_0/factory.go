package v3_0

import nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0"

func newValue() *nthrift.Value {
	return nthrift.NewValue()
}
