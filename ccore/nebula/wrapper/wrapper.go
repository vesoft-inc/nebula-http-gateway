package wrapper

import (
	"fmt"
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

// construct go-type to nebula.Value
func WrapValue(any interface{}, factory types.FactoryDriver) (types.Value, error) {
	var err error
	valueBuilder := factory.NewValueBuilder()
	value := valueBuilder.Build()
	if v, ok := any.(bool); ok {
		value.SetBVal(&v)
	} else if v, ok := any.(int); ok {
		ival := int64(v)
		value.SetIVal(&ival)
	} else if v, ok := any.(float64); ok {
		if v == float64(int64(v)) {
			iv := int64(v)
			value.SetIVal(&iv)
		} else {
			value.SetFVal(&v)
		}
	} else if v, ok := any.(float32); ok {
		if v == float32(int64(v)) {
			iv := int64(v)
			value.SetIVal(&iv)
		} else {
			fval := float64(v)
			value.SetFVal(&fval)
		}
	} else if v, ok := any.(string); ok {
		value.SetSVal([]byte(v))
	} else if any == nil {
		nval := types.NullType___NULL__
		value.SetNVal(&nval)
	} else if v, ok := any.([]interface{}); ok {
		nv, er := Slice2Nlist([]interface{}(v), factory)
		if er != nil {
			err = er
		}
		value.SetLVal(nv)
	} else if v, ok := any.(map[string]interface{}); ok {
		nv, er := Map2Nmap(v, factory)
		if er != nil {
			err = er
		}
		value.SetMVal(nv)
	} else if v, ok := any.(types.Value); ok {
		value = v
	} else if v, ok := any.(types.Date); ok {
		value.SetDVal(v)
	} else if v, ok := any.(types.DateTime); ok {
		value.SetDtVal(v)
	} else if v, ok := any.(types.Duration); ok {
		value.SetDuVal(v)
	} else if v, ok := any.(types.Time); ok {
		value.SetTVal(v)
	} else if v, ok := any.(types.Geography); ok {
		value.SetGgVal(v)
	} else {
		// unsupport other Value type, use this function carefully
		err = fmt.Errorf("Only support convert boolean/float/int/string/map/list to nebula.Value but %T", any)
	}
	return value, err
}

// construct Slice to nebula.NList
func Slice2Nlist(list []interface{}, factory types.FactoryDriver) (types.NList, error) {
	sv := make([]types.Value, 0, len(list))
	nListBuilder := factory.NewNListBuilder()
	for _, item := range list {
		nv, er := WrapValue(item, factory)
		if er != nil {
			return nil, er
		}
		sv = append(sv, nv)
	}
	nListBuilder.Values(sv)
	return nListBuilder.Build(), nil
}

// construct map to nebula.NMap
func Map2Nmap(m map[string]interface{}, factory types.FactoryDriver) (types.NMap, error) {
	nMapBuilder := factory.NewNMapBuilder()
	kvs := map[string]types.Value{}
	for k, v := range m {
		nv, err := WrapValue(v, factory)
		if err != nil {
			return nil, err
		}
		kvs[k] = nv
	}
	nMapBuilder.Kvs(kvs)
	return nMapBuilder.Build(), nil
}
