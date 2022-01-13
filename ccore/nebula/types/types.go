package types

import (
	"fmt"
)

type Any = interface{}

type ParameterList []string
type ParameterMap map[string]interface{}

/*
	define the interfaces to fit internal/<ver>/thrift/ttypes and  internal/<ver>/thrift/graph/ttypes
*/

type TimezoneInfo struct {
	offset int32
	name   []byte
}

func (t *TimezoneInfo) GetOffset() int32 {
	return t.offset
}

func (t *TimezoneInfo) GetName() []byte {
	return t.name
}

func (t *TimezoneInfo) SetOffset(offset int32) *TimezoneInfo {
	t.offset = offset
	return t
}

func (t *TimezoneInfo) SetName(name []byte) *TimezoneInfo {
	t.name = name
	return t
}

type EdgeType = int32

type EdgeRanking = int64

type DataSet interface {
	GetColumnNames() [][]byte
	GetRows() []Row
	String() string
	Unwrap() interface{}
}

type Row interface {
	GetValues() []Value
	String() string
	Unwrap() interface{}
}

type Value interface {
	GetNVal() NullType
	GetBVal() bool
	GetIVal() int64
	GetFVal() float64
	GetSVal() []byte
	GetDVal() Date
	GetTVal() Time
	GetDtVal() DateTime
	GetVVal() Vertex
	GetEVal() Edge
	GetPVal() Path
	GetLVal() NList
	GetMVal() NMap
	GetUVal() NSet
	GetGVal() DataSet
	GetGgVal() Geography
	GetDuVal() Duration
	CountSetFieldsValue() int
	IsSetNVal() bool
	IsSetBVal() bool
	IsSetIVal() bool
	IsSetFVal() bool
	IsSetSVal() bool
	IsSetDVal() bool
	IsSetTVal() bool
	IsSetDtVal() bool
	IsSetVVal() bool
	IsSetEVal() bool
	IsSetPVal() bool
	IsSetLVal() bool
	IsSetMVal() bool
	IsSetUVal() bool
	IsSetGVal() bool
	IsSetGgVal() bool
	IsSetDuVal() bool
	SetNVal(*NullType) Value
	SetBVal(*bool) Value
	SetIVal(*int64) Value
	SetFVal(*float64) Value
	SetSVal([]byte) Value
	SetDVal(Date) Value
	SetTVal(Time) Value
	SetDtVal(DateTime) Value
	SetVVal(Vertex) Value
	SetEVal(Edge) Value
	SetPVal(Path) Value
	SetLVal(NList) Value
	SetMVal(NMap) Value
	SetUVal(NSet) Value
	SetGVal(DataSet) Value
	SetGgVal(Geography) Value
	SetDuVal(Duration) Value
	String() string
	Unwrap() interface{}
}

type NullType int64

const (
	NullType___NULL__     NullType = 0
	NullType_NaN          NullType = 1
	NullType_BAD_DATA     NullType = 2
	NullType_BAD_TYPE     NullType = 3
	NullType_ERR_OVERFLOW NullType = 4
	NullType_UNKNOWN_PROP NullType = 5
	NullType_DIV_BY_ZERO  NullType = 6
	NullType_OUT_OF_RANGE NullType = 7
)

var NullTypeToName = map[NullType]string{
	NullType___NULL__:     "__NULL__",
	NullType_NaN:          "NaN",
	NullType_BAD_DATA:     "BAD_DATA",
	NullType_BAD_TYPE:     "BAD_TYPE",
	NullType_ERR_OVERFLOW: "ERR_OVERFLOW",
	NullType_UNKNOWN_PROP: "UNKNOWN_PROP",
	NullType_DIV_BY_ZERO:  "DIV_BY_ZERO",
	NullType_OUT_OF_RANGE: "OUT_OF_RANGE",
}

var NullTypeToValue = map[string]NullType{
	"__NULL__":     NullType___NULL__,
	"NaN":          NullType_NaN,
	"BAD_DATA":     NullType_BAD_DATA,
	"BAD_TYPE":     NullType_BAD_TYPE,
	"ERR_OVERFLOW": NullType_ERR_OVERFLOW,
	"UNKNOWN_PROP": NullType_UNKNOWN_PROP,
	"DIV_BY_ZERO":  NullType_DIV_BY_ZERO,
	"OUT_OF_RANGE": NullType_OUT_OF_RANGE,
}

var NullTypeNames = []string{
	"__NULL__",
	"NaN",
	"BAD_DATA",
	"BAD_TYPE",
	"ERR_OVERFLOW",
	"UNKNOWN_PROP",
	"DIV_BY_ZERO",
	"OUT_OF_RANGE",
}

var NullTypeValues = []NullType{
	NullType___NULL__,
	NullType_NaN,
	NullType_BAD_DATA,
	NullType_BAD_TYPE,
	NullType_ERR_OVERFLOW,
	NullType_UNKNOWN_PROP,
	NullType_DIV_BY_ZERO,
	NullType_OUT_OF_RANGE,
}

func (p NullType) String() string {
	if v, ok := NullTypeToName[p]; ok {
		return v
	}
	return "<UNSET>"
}

func NullTypeFromString(s string) (NullType, error) {
	if v, ok := NullTypeToValue[s]; ok {
		return v, nil
	}
	return NullType(0), fmt.Errorf("not a valid NullType string")
}

func NullTypePtr(v NullType) *NullType { return &v }

type Date interface {
	GetYear() int16
	GetMonth() int8
	GetDay() int8
	SetYear(int16) Date
	SetMonth(int8) Date
	SetDay(int8) Date
	String() string
	Unwrap() interface{}
}

type Time interface {
	GetHour() int8
	GetMinute() int8
	GetSec() int8
	GetMicrosec() int32
	SetHour(int8) Time
	SetMinute(int8) Time
	SetSec(int8) Time
	SetMicrosec(int32) Time
	String() string
	Unwrap() interface{}
}

type DateTime interface {
	GetYear() int16
	GetMonth() int8
	GetDay() int8
	GetHour() int8
	GetMinute() int8
	GetSec() int8
	GetMicrosec() int32
	SetYear(int16) DateTime
	SetMonth(int8) DateTime
	SetDay(int8) DateTime
	SetHour(int8) DateTime
	SetMinute(int8) DateTime
	SetSec(int8) DateTime
	SetMicrosec(int32) DateTime
	String() string
	Unwrap() interface{}
}

type Vertex interface {
	GetVid() Value
	GetTags() []Tag
	IsSetVid() bool
	String() string
	Unwrap() interface{}
}

type Edge interface {
	GetSrc() Value
	GetDst() Value
	GetType() EdgeType
	GetName() []byte
	GetRanking() EdgeRanking
	GetProps() map[string]Value
	IsSetSrc() bool
	IsSetDst() bool
	SetSrc(Value) Edge
	SetDst(Value) Edge
	SetType(EdgeType) Edge
	SetName([]byte) Edge
	SetRanking(EdgeRanking) Edge
	SetProps(map[string]Value) Edge
	String() string
	Unwrap() interface{}
}

type Path interface {
	GetSrc() Vertex
	GetSteps() []Step
	IsSetSrc() bool
	String() string
	Unwrap() interface{}
}

type NList interface {
	GetValues() []Value
	SetValues([]Value) NList
	String() string
	Unwrap() interface{}
}

type NMap interface {
	GetKvs() map[string]Value
	SetKvs(map[string]Value) NMap
	String() string
	Unwrap() interface{}
}

type NSet interface {
	GetValues() []Value
	String() string
	Unwrap() interface{}
}

type Geography interface {
	GetPtVal() Point
	GetLsVal() LineString
	GetPgVal() Polygon
	CountSetFieldsGeography() int
	IsSetPtVal() bool
	IsSetLsVal() bool
	IsSetPgVal() bool
	String() string
	Unwrap() interface{}
}

type Tag interface {
	GetName() []byte
	GetProps() map[string]Value
	String() string
	Unwrap() interface{}
}

type Step interface {
	GetDst() Vertex
	GetType() EdgeType
	GetName() []byte
	GetRanking() EdgeRanking
	GetProps() map[string]Value
	IsSetDst() bool
	String() string
	Unwrap() interface{}
}

type Point interface {
	GetCoord() Coordinate
	IsSetCoord() bool
	String() string
	Unwrap() interface{}
}

type LineString interface {
	GetCoordList() []Coordinate
	String() string
	Unwrap() interface{}
}

type Polygon interface {
	GetCoordListList() [][]Coordinate
	String() string
	Unwrap() interface{}
}

type Coordinate interface {
	GetX() float64
	GetY() float64
	String() string
	Unwrap() interface{}
}

type Duration interface {
	GetSeconds() int64
	GetMicroseconds() int32
	GetMonths() int32
	String() string
	Unwrap() interface{}
}

type PlanDescription interface {
	GetPlanNodeDescs() []PlanNodeDescription
	GetNodeIndexMap() map[int64]int64
	GetFormat() []byte
	GetOptimizeTimeInUs() int32
	String() string
	Unwrap() interface{}
}

type PlanNodeDescription interface {
	GetName() []byte
	GetId() int64
	GetOutputVar() []byte
	GetDescription() []Pair
	GetProfiles() []ProfilingStats
	GetBranchInfo() PlanNodeBranchInfo
	GetDependencies() []int64
	IsSetDescription() bool
	IsSetProfiles() bool
	IsSetBranchInfo() bool
	IsSetDependencies() bool
	String() string
	Unwrap() interface{}
}

type Pair interface {
	GetKey() []byte
	GetValue() []byte
	String() string
	Unwrap() interface{}
}

type ProfilingStats interface {
	GetRows() int64
	GetExecDurationInUs() int64
	GetTotalDurationInUs() int64
	GetOtherStats() map[string][]byte
	IsSetOtherStats() bool
	String() string
	Unwrap() interface{}
}

type PlanNodeBranchInfo interface {
	GetIsDoBranch() bool
	GetConditionNodeID() int64
	String() string
	Unwrap() interface{}
}

type (
	BalanceCmd   string
	BalanceStats string
)

const (
	BalanceData       = BalanceCmd("Balance Data")
	BalanceDataRemove = BalanceCmd("Balance Data Remove")
	BalanceLeader     = BalanceCmd("Balance Leader")

	Balanced   = BalanceStats("Balanced")
	ImBalanced = BalanceStats("ImBalanced")
	Balancing  = BalanceStats("Balancing")
)

type BalanceReq struct {
	Cmd           BalanceCmd
	Space         string
	HostsToRemove []string
}
