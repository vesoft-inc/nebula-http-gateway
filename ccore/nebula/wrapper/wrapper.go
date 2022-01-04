package wrapper

import (
	"time"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

func WrapTime(localTime time.Time, factory types.FactoryDriver) types.Time {
	tb := factory.NewTimeBuilder()
	t := tb.Build()
	t.SetHour(int8(localTime.Hour()))
	t.SetMinute(int8(localTime.Minute()))
	t.SetSec(int8(localTime.Second()))
	t.SetMicrosec(int32(localTime.Nanosecond() / 1000))

	return t
}

func WrapDateTime(localDT time.Time, factory types.FactoryDriver) types.DateTime {
	dtb := factory.NewDateTimeBuilder()
	datetime := dtb.Build()
	datetime.SetYear(int16(localDT.Year()))
	datetime.SetMonth(int8(localDT.Month()))
	datetime.SetDay(int8(localDT.Day()))
	datetime.SetHour(int8(localDT.Hour()))
	datetime.SetMinute(int8(localDT.Minute()))
	datetime.SetSec(int8(localDT.Second()))
	datetime.SetMicrosec(int32(localDT.Nanosecond() / 1000))

	return datetime
}
