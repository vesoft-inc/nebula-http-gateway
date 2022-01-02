package v2_5

import nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5"

func newValue() *nthrift.Value {
	return nthrift.NewValue()
}
