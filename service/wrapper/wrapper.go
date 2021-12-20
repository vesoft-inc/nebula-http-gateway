package wrapper

import (
	"fmt"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

/*
	copy from ccore/nebula/internal/thrift/ttypes.go, and fit with ccore/nebula/types
*/

type defaultValue struct {
	NVal  types.NullType  `thrift:"nVal,1" db:"nVal" json:"nVal,omitempty"`
	BVal  *bool           `thrift:"bVal,2" db:"bVal" json:"bVal,omitempty"`
	IVal  *int64          `thrift:"iVal,3" db:"iVal" json:"iVal,omitempty"`
	FVal  *float64        `thrift:"fVal,4" db:"fVal" json:"fVal,omitempty"`
	SVal  []byte          `thrift:"sVal,5" db:"sVal" json:"sVal,omitempty"`
	DVal  types.Date      `thrift:"dVal,6" db:"dVal" json:"dVal,omitempty"`
	TVal  types.Time      `thrift:"tVal,7" db:"tVal" json:"tVal,omitempty"`
	DtVal types.DateTime  `thrift:"dtVal,8" db:"dtVal" json:"dtVal,omitempty"`
	VVal  types.Vertex    `thrift:"vVal,9" db:"vVal" json:"vVal,omitempty"`
	EVal  types.Edge      `thrift:"eVal,10" db:"eVal" json:"eVal,omitempty"`
	PVal  types.Path      `thrift:"pVal,11" db:"pVal" json:"pVal,omitempty"`
	LVal  types.NList     `thrift:"lVal,12" db:"lVal" json:"lVal,omitempty"`
	MVal  types.NMap      `thrift:"mVal,13" db:"mVal" json:"mVal,omitempty"`
	UVal  types.NSet      `thrift:"uVal,14" db:"uVal" json:"uVal,omitempty"`
	GVal  types.DataSet   `thrift:"gVal,15" db:"gVal" json:"gVal,omitempty"`
	GgVal types.Geography `thrift:"ggVal,16" db:"ggVal" json:"ggVal,omitempty"`
}

func NewValue(vval types.Vertex) types.Value {
	return &defaultValue{VVal: vval}
}

var Value_NVal_DEFAULT types.NullType

func (p *defaultValue) GetNVal() types.NullType {
	if !p.IsSetNVal() {
		return Value_NVal_DEFAULT
	}
	return p.NVal
}

var Value_BVal_DEFAULT bool

func (p *defaultValue) GetBVal() bool {
	if !p.IsSetBVal() {
		return Value_BVal_DEFAULT
	}
	return *p.BVal
}

var Value_IVal_DEFAULT int64

func (p *defaultValue) GetIVal() int64 {
	if !p.IsSetIVal() {
		return Value_IVal_DEFAULT
	}
	return *p.IVal
}

var Value_FVal_DEFAULT float64

func (p *defaultValue) GetFVal() float64 {
	if !p.IsSetFVal() {
		return Value_FVal_DEFAULT
	}
	return *p.FVal
}

var Value_SVal_DEFAULT []byte

func (p *defaultValue) GetSVal() []byte {
	return p.SVal
}

var Value_DVal_DEFAULT types.Date

func (p *defaultValue) GetDVal() types.Date {
	if !p.IsSetDVal() {
		return Value_DVal_DEFAULT
	}
	return p.DVal
}

var Value_TVal_DEFAULT types.Time

func (p *defaultValue) GetTVal() types.Time {
	if !p.IsSetTVal() {
		return Value_TVal_DEFAULT
	}
	return p.TVal
}

var Value_DtVal_DEFAULT types.DateTime

func (p *defaultValue) GetDtVal() types.DateTime {
	if !p.IsSetDtVal() {
		return Value_DtVal_DEFAULT
	}
	return p.DtVal
}

var Value_VVal_DEFAULT types.Vertex

func (p *defaultValue) GetVVal() types.Vertex {
	if !p.IsSetVVal() {
		return Value_VVal_DEFAULT
	}
	return p.VVal
}

var Value_EVal_DEFAULT types.Edge

func (p *defaultValue) GetEVal() types.Edge {
	if !p.IsSetEVal() {
		return Value_EVal_DEFAULT
	}
	return p.EVal
}

var Value_PVal_DEFAULT types.Path

func (p *defaultValue) GetPVal() types.Path {
	if !p.IsSetPVal() {
		return Value_PVal_DEFAULT
	}
	return p.PVal
}

var Value_LVal_DEFAULT types.NList

func (p *defaultValue) GetLVal() types.NList {
	if !p.IsSetLVal() {
		return Value_LVal_DEFAULT
	}
	return p.LVal
}

var Value_MVal_DEFAULT types.NMap

func (p *defaultValue) GetMVal() types.NMap {
	if !p.IsSetMVal() {
		return Value_MVal_DEFAULT
	}
	return p.MVal
}

var Value_UVal_DEFAULT types.NSet

func (p *defaultValue) GetUVal() types.NSet {
	if !p.IsSetUVal() {
		return Value_UVal_DEFAULT
	}
	return p.UVal
}

var Value_GVal_DEFAULT types.DataSet

func (p *defaultValue) GetGVal() types.DataSet {
	if !p.IsSetGVal() {
		return Value_GVal_DEFAULT
	}
	return p.GVal
}

var Value_GgVal_DEFAULT types.Geography

func (p *defaultValue) GetGgVal() types.Geography {
	if !p.IsSetGgVal() {
		return Value_GgVal_DEFAULT
	}
	return p.GgVal
}
func (p *defaultValue) CountSetFieldsValue() int {
	count := 0
	if p.IsSetNVal() {
		count++
	}
	if p.IsSetBVal() {
		count++
	}
	if p.IsSetIVal() {
		count++
	}
	if p.IsSetFVal() {
		count++
	}
	if p.IsSetSVal() {
		count++
	}
	if p.IsSetDVal() {
		count++
	}
	if p.IsSetTVal() {
		count++
	}
	if p.IsSetDtVal() {
		count++
	}
	if p.IsSetVVal() {
		count++
	}
	if p.IsSetEVal() {
		count++
	}
	if p.IsSetPVal() {
		count++
	}
	if p.IsSetLVal() {
		count++
	}
	if p.IsSetMVal() {
		count++
	}
	if p.IsSetUVal() {
		count++
	}
	if p.IsSetGVal() {
		count++
	}
	if p.IsSetGgVal() {
		count++
	}
	return count

}

func (p *defaultValue) IsSetNVal() bool {
	return p != nil && p.NVal != types.NullTypeToValue["__NULL__"]
}

func (p *defaultValue) IsSetBVal() bool {
	return p != nil && p.BVal != nil
}

func (p *defaultValue) IsSetIVal() bool {
	return p != nil && p.IVal != nil
}

func (p *defaultValue) IsSetFVal() bool {
	return p != nil && p.FVal != nil
}

func (p *defaultValue) IsSetSVal() bool {
	return p != nil && p.SVal != nil
}

func (p *defaultValue) IsSetDVal() bool {
	return p != nil && p.DVal != nil
}

func (p *defaultValue) IsSetTVal() bool {
	return p != nil && p.TVal != nil
}

func (p *defaultValue) IsSetDtVal() bool {
	return p != nil && p.DtVal != nil
}

func (p *defaultValue) IsSetVVal() bool {
	return p != nil && p.VVal != nil
}

func (p *defaultValue) IsSetEVal() bool {
	return p != nil && p.EVal != nil
}

func (p *defaultValue) IsSetPVal() bool {
	return p != nil && p.PVal != nil
}

func (p *defaultValue) IsSetLVal() bool {
	return p != nil && p.LVal != nil
}

func (p *defaultValue) IsSetMVal() bool {
	return p != nil && p.MVal != nil
}

func (p *defaultValue) IsSetUVal() bool {
	return p != nil && p.UVal != nil
}

func (p *defaultValue) IsSetGVal() bool {
	return p != nil && p.GVal != nil
}

func (p *defaultValue) IsSetGgVal() bool {
	return p != nil && p.GgVal != nil
}

func (p *defaultValue) String() string {
	if p == nil {
		return "<nil>"
	}

	var nValVal string
	if p.NVal == types.NullTypeToValue["__NULL__"] {
		nValVal = "<nil>"
	} else {
		nValVal = fmt.Sprintf("%v", p.NVal)
	}
	var bValVal string
	if p.BVal == nil {
		bValVal = "<nil>"
	} else {
		bValVal = fmt.Sprintf("%v", *p.BVal)
	}
	var iValVal string
	if p.IVal == nil {
		iValVal = "<nil>"
	} else {
		iValVal = fmt.Sprintf("%v", *p.IVal)
	}
	var fValVal string
	if p.FVal == nil {
		fValVal = "<nil>"
	} else {
		fValVal = fmt.Sprintf("%v", *p.FVal)
	}
	sValVal := fmt.Sprintf("%v", p.SVal)
	var dValVal string
	if p.DVal == nil {
		dValVal = "<nil>"
	} else {
		dValVal = fmt.Sprintf("%v", p.DVal)
	}
	var tValVal string
	if p.TVal == nil {
		tValVal = "<nil>"
	} else {
		tValVal = fmt.Sprintf("%v", p.TVal)
	}
	var dtValVal string
	if p.DtVal == nil {
		dtValVal = "<nil>"
	} else {
		dtValVal = fmt.Sprintf("%v", p.DtVal)
	}
	var vValVal string
	if p.VVal == nil {
		vValVal = "<nil>"
	} else {
		vValVal = fmt.Sprintf("%v", p.VVal)
	}
	var eValVal string
	if p.EVal == nil {
		eValVal = "<nil>"
	} else {
		eValVal = fmt.Sprintf("%v", p.EVal)
	}
	var pValVal string
	if p.PVal == nil {
		pValVal = "<nil>"
	} else {
		pValVal = fmt.Sprintf("%v", p.PVal)
	}
	var lValVal string
	if p.LVal == nil {
		lValVal = "<nil>"
	} else {
		lValVal = fmt.Sprintf("%v", p.LVal)
	}
	var mValVal string
	if p.MVal == nil {
		mValVal = "<nil>"
	} else {
		mValVal = fmt.Sprintf("%v", p.MVal)
	}
	var uValVal string
	if p.UVal == nil {
		uValVal = "<nil>"
	} else {
		uValVal = fmt.Sprintf("%v", p.UVal)
	}
	var gValVal string
	if p.GVal == nil {
		gValVal = "<nil>"
	} else {
		gValVal = fmt.Sprintf("%v", p.GVal)
	}
	var ggValVal string
	if p.GgVal == nil {
		ggValVal = "<nil>"
	} else {
		ggValVal = fmt.Sprintf("%v", p.GgVal)
	}
	return fmt.Sprintf("Value({NVal:%s BVal:%s IVal:%s FVal:%s SVal:%s DVal:%s TVal:%s DtVal:%s VVal:%s EVal:%s PVal:%s LVal:%s MVal:%s UVal:%s GVal:%s GgVal:%s})", nValVal, bValVal, iValVal, fValVal, sValVal, dValVal, tValVal, dtValVal, vValVal, eValVal, pValVal, lValVal, mValVal, uValVal, gValVal, ggValVal)
}

type defaultEdge struct {
	Src     types.Value            `thrift:"src,1" db:"src" json:"src"`
	Dst     types.Value            `thrift:"dst,2" db:"dst" json:"dst"`
	Type    types.EdgeType         `thrift:"type,3" db:"type" json:"type"`
	Name    []byte                 `thrift:"name,4" db:"name" json:"name"`
	Ranking types.EdgeRanking      `thrift:"ranking,5" db:"ranking" json:"ranking"`
	Props   map[string]types.Value `thrift:"props,6" db:"props" json:"props"`
}

func NewEdge(src, dst types.Value, t types.EdgeType, name []byte, ranking types.EdgeRanking, props map[string]types.Value) types.Edge {
	edge := &defaultEdge{src, dst, t, name, ranking, props}
	return edge
}

var Edge_Src_DEFAULT types.Value

func (p *defaultEdge) GetSrc() types.Value {
	if !p.IsSetSrc() {
		return Edge_Src_DEFAULT
	}
	return p.Src
}

var Edge_Dst_DEFAULT types.Value

func (p *defaultEdge) GetDst() types.Value {
	if !p.IsSetDst() {
		return Edge_Dst_DEFAULT
	}
	return p.Dst
}

func (p *defaultEdge) GetType() types.EdgeType {
	return p.Type
}

func (p *defaultEdge) GetName() []byte {
	return p.Name
}

func (p *defaultEdge) GetRanking() types.EdgeRanking {
	return p.Ranking
}

func (p *defaultEdge) GetProps() map[string]types.Value {
	return p.Props
}
func (p *defaultEdge) IsSetSrc() bool {
	return p != nil && p.Src != nil
}

func (p *defaultEdge) IsSetDst() bool {
	return p != nil && p.Dst != nil
}

func (p *defaultEdge) String() string {
	if p == nil {
		return "<nil>"
	}

	var srcVal string
	if p.Src == nil {
		srcVal = "<nil>"
	} else {
		srcVal = fmt.Sprintf("%v", p.Src)
	}
	var dstVal string
	if p.Dst == nil {
		dstVal = "<nil>"
	} else {
		dstVal = fmt.Sprintf("%v", p.Dst)
	}
	typeVal := fmt.Sprintf("%v", p.Type)
	nameVal := fmt.Sprintf("%v", p.Name)
	rankingVal := fmt.Sprintf("%v", p.Ranking)
	propsVal := fmt.Sprintf("%v", p.Props)
	return fmt.Sprintf("Edge({Src:%s Dst:%s Type:%s Name:%s Ranking:%s Props:%s})", srcVal, dstVal, typeVal, nameVal, rankingVal, propsVal)
}
