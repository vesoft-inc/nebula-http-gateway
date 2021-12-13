package types

import "fmt"

type EdgeType = int32

type EdgeRanking = int64

type DataSet interface {
	GetColumnNames() [][]byte
	GetRows() []Row
	String() string
}

type Row interface {
	GetValues() []Value
	String() string
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
	CountSetFieldsValue() int
	IsSetNVal() bool
	IsSetBVal() bool
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
	String() string
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
	String() string
}

type Time interface {
	GetHour() int8
	GetMinute() int8
	GetSec() int8
	GetMicrosec() int32
	String() string
}

type DateTime interface {
	GetYear() int16
	GetMonth() int8
	GetDay() int8
	GetHour() int8
	GetMinute() int8
	GetSec() int8
	GetMicrosec() int32
	String() string
}

type Vertex interface {
	GetVid() Value
	GetTags() []Tag
	IsSetVid() bool
	String() string
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
	String() string
}

type Path interface {
	GetSrc() Vertex
	GetSteps() []Step
	IsSetSrc() bool
	String() string
}

type NList interface {
	GetValues() []Value
	String() string
}

type NMap interface {
	GetKvs() map[string]Value
	String() string
}

type NSet interface {
	GetValues() []Value
	String() string
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
}

type Tag interface {
	GetName() []byte
	GetProps() map[string]Value
	String() string
}

type Step interface {
	GetDst() Vertex
	GetType() EdgeType
	GetName() []byte
	GetRanking() EdgeRanking
	GetProps() map[string]Value
	IsSetDst() bool
	String() string
}

type Point interface {
	GetCoord() Coordinate
	IsSetCoord() bool
	String() string
}

type LineString interface {
	GetCoordList() []Coordinate
	String() string
}

type Polygon interface {
	GetCoordListList() [][]Coordinate
	String() string
}

type Coordinate interface {
	GetX() float64
	GetY() float64
	String() string
}

type PlanDescription interface {
	GetPlanNodeDescs() []PlanNodeDescription
	GetNodeIndexMap() map[int64]int64
	GetFormat() []byte
	GetOptimizeTimeInUs() int32
	String() string
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
}

type Pair interface {
	GetKey() []byte
	GetValue() []byte
	String() string
}

type ProfilingStats interface {
	GetRows() int64
	GetExecDurationInUs() int64
	GetTotalDurationInUs() int64
	GetOtherStats() map[string][]byte
	IsSetOtherStats() bool
	String() string
}

type PlanNodeBranchInfo interface {
	GetIsDoBranch() bool
	GetConditionNodeID() int64
	String() string
}
