package v2_6

import nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_6"

func newValue() *nthrift.Value {
	return nthrift.NewValue()
}
