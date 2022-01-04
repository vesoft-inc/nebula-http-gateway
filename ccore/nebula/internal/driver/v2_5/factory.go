package v2_5

import nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5"

func newValue() *nthrift.Value {
	return nthrift.NewValue()
}

func newDate() *nthrift.Date {
	return nthrift.NewDate()
}

func newTime() *nthrift.Time {
	return nthrift.NewTime()
}

func newDateTime() *nthrift.DateTime {
	return nthrift.NewDateTime()
}

func newEdge() *nthrift.Edge {
	return nthrift.NewEdge()
}
