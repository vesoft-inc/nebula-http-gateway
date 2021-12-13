package v2_5_0

import (
	nebula "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5_0"
	graph "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5_0/graph"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type defaultDataSet struct {
	*nebula.DataSet
}

func (d defaultDataSet) GetRows() []types.Row {
	return rowsWrapper(d.Rows)
}

func dataSetWrapper(ds *nebula.DataSet) types.DataSet {
	return defaultDataSet{ds}
}

type defaultRow struct {
	*nebula.Row
}

func (d defaultRow) GetValues() []types.Value {
	return vaulesWrapper(d.Values)
}

func rowWrapper(r *nebula.Row) types.Row {
	return defaultRow{r}
}

func rowsWrapper(rs []*nebula.Row) []types.Row {
	rows := make([]types.Row, len(rs))
	for i := range rows {
		row := rowWrapper(rs[i])
		rows[i] = row
	}
	return rows
}

type defaultValue struct {
	*nebula.Value
}

func (d defaultValue) GetNVal() types.NullType {
	return nullTypeWrapper(d.NVal)
}

func (d defaultValue) GetDVal() types.Date {
	date := dateWrapper(d.DVal)
	return date
}

func (d defaultValue) GetTVal() types.Time {
	time := timeWrapper(d.TVal)
	return time
}

func (d defaultValue) GetDtVal() types.DateTime {
	dateTime := dateTimeWrapper(d.DtVal)
	return dateTime
}

func (d defaultValue) GetVVal() types.Vertex {
	vertex := vertexWrapper(d.VVal)
	return vertex
}

func (d defaultValue) GetEVal() types.Edge {
	edge := edgeWrapper(d.EVal)
	return edge
}

func (d defaultValue) GetPVal() types.Path {
	path := pathWrapper(d.PVal)
	return path
}

func (d defaultValue) GetLVal() types.NList {
	nlist := nlistWrapper(d.LVal)
	return nlist
}

func (d defaultValue) GetMVal() types.NMap {
	nmap := nmapWrapper(d.MVal)
	return nmap
}

func (d defaultValue) GetUVal() types.NSet {
	nset := nsetWrapper(d.UVal)
	return nset
}

func (d defaultValue) GetGVal() types.DataSet {
	dataset := dataSetWrapper(d.GVal)
	return dataset
}

func (d defaultValue) GetGgVal() types.Geography {
	panic("method not support")
}

func (d defaultValue) IsSetGgVal() bool {
	return false
}

func valueWrapper(v *nebula.Value) types.Value {
	return defaultValue{v}
}

func vaulesWrapper(vs []*nebula.Value) []types.Value {
	values := make([]types.Value, len(vs))
	for i := range values {
		value := valueWrapper(vs[i])
		values[i] = value
	}
	return values
}

func nullTypeWrapper(nt *nebula.NullType) types.NullType {
	return types.NullTypeToValue[nt.String()]
}

type defaultDate struct {
	*nebula.Date
}

func dateWrapper(d *nebula.Date) types.Date {
	return defaultDate{d}
}

type defaultTime struct {
	*nebula.Time
}

func timeWrapper(t *nebula.Time) types.Time {
	return defaultTime{t}
}

type defaultDateTime struct {
	*nebula.DateTime
}

func dateTimeWrapper(dt *nebula.DateTime) types.DateTime {
	return defaultDateTime{dt}
}

type defaultVertex struct {
	*nebula.Vertex
}

func (d defaultVertex) GetVid() types.Value {
	value := valueWrapper(d.Vid)
	return value
}

func (d defaultVertex) GetTags() []types.Tag {
	return tagsWrapper(d.Tags)
}

func vertexWrapper(v *nebula.Vertex) types.Vertex {
	return defaultVertex{v}
}

type defaultEdge struct {
	*nebula.Edge
}

func (d defaultEdge) GetSrc() types.Value {
	value := valueWrapper(d.Src)
	return value
}

func (d defaultEdge) GetDst() types.Value {
	value := valueWrapper(d.Dst)
	return value
}

func (d defaultEdge) GetType() types.EdgeType {
	return edgeTypeWrapper(d.Type)
}

func (d defaultEdge) GetRanking() types.EdgeRanking {
	return edgeRankingWrapper(d.Ranking)
}

func (d defaultEdge) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(d.Props))
	for k, v := range d.Props {
		value := valueWrapper(v)
		props[k] = value
	}
	return props
}

func edgeWrapper(e *nebula.Edge) types.Edge {
	return defaultEdge{e}
}

func edgeTypeWrapper(et nebula.EdgeType) types.EdgeType {
	return types.EdgeType(int32(et))
}

func edgeRankingWrapper(et nebula.EdgeRanking) types.EdgeRanking {
	return types.EdgeRanking(int64(et))
}

type defaultPath struct {
	*nebula.Path
}

func (d defaultPath) GetSrc() types.Vertex {
	src := vertexWrapper(d.Src)
	return src
}
func (d defaultPath) GetSteps() []types.Step {
	return stepsWrapper(d.Steps)
}

func pathWrapper(p *nebula.Path) types.Path {
	return defaultPath{p}
}

type defaultNList struct {
	*nebula.NList
}

func (d defaultNList) GetValues() []types.Value {
	return vaulesWrapper(d.Values)
}

func nlistWrapper(nl *nebula.NList) types.NList {
	return defaultNList{nl}
}

type defaultNMap struct {
	*nebula.NMap
}

func (d defaultNMap) GetKvs() map[string]types.Value {
	kvs := make(map[string]types.Value, len(d.Kvs))
	for k, v := range d.Kvs {
		value := valueWrapper(v)
		kvs[k] = value
	}
	return kvs
}

func nmapWrapper(nm *nebula.NMap) types.NMap {
	return defaultNMap{nm}
}

type defaultNSet struct {
	*nebula.NSet
}

func (d defaultNSet) GetValues() []types.Value {
	return vaulesWrapper(d.Values)
}

func nsetWrapper(ns *nebula.NSet) types.NSet {
	return defaultNSet{ns}
}

type defaultTag struct {
	*nebula.Tag
}

func (d defaultTag) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(d.Props))
	for k, v := range d.Props {
		value := valueWrapper(v)
		props[k] = value
	}
	return props
}

func tagWrapper(t *nebula.Tag) types.Tag {
	return defaultTag{t}
}

func tagsWrapper(ts []*nebula.Tag) []types.Tag {
	tags := make([]types.Tag, len(ts))
	for i := range ts {
		tag := tagWrapper(ts[i])
		tags[i] = tag
	}
	return tags
}

type defaultStep struct {
	*nebula.Step
}

func (d defaultStep) GetDst() types.Vertex {
	dst := vertexWrapper(d.Dst)
	return dst
}

func (d defaultStep) GetType() types.EdgeType {
	return edgeTypeWrapper(d.Type)
}

func (d defaultStep) GetRanking() types.EdgeRanking {
	return edgeRankingWrapper(d.Ranking)
}

func (d defaultStep) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(d.Props))
	for k, v := range d.Props {
		value := valueWrapper(v)
		props[k] = value
	}
	return props
}

func stepWrapper(s *nebula.Step) types.Step {
	return defaultStep{s}
}

func stepsWrapper(ss []*nebula.Step) []types.Step {
	steps := make([]types.Step, len(ss))
	for i := range ss {
		step := stepWrapper(ss[i])
		steps[i] = step
	}
	return steps
}

type defaultPlanDescription struct {
	*graph.PlanDescription
}

func (d defaultPlanDescription) GetPlanNodeDescs() []types.PlanNodeDescription {
	return planNodeDescriptionsWrapper(d.PlanNodeDescs)
}

func planDescriptionWrapper(pd *graph.PlanDescription) types.PlanDescription {
	return defaultPlanDescription{pd}
}

type defaultPlanNodeDescription struct {
	*graph.PlanNodeDescription
}

func (d defaultPlanNodeDescription) GetDescription() []types.Pair {
	return pairsWrapper(d.Description)
}

func (d defaultPlanNodeDescription) GetProfiles() []types.ProfilingStats {
	return profilingStatssWrapper(d.Profiles)
}

func (d defaultPlanNodeDescription) GetBranchInfo() types.PlanNodeBranchInfo {
	return planNodeBranchInfoWrapper(d.BranchInfo)
}

func planNodeDescriptionWrapper(pnd *graph.PlanNodeDescription) types.PlanNodeDescription {
	return defaultPlanNodeDescription{pnd}
}

func planNodeDescriptionsWrapper(pnds []*graph.PlanNodeDescription) []types.PlanNodeDescription {
	planNodeDescriptions := make([]types.PlanNodeDescription, len(pnds))
	for i := range pnds {
		planNodeDescriptions[i] = planNodeDescriptionWrapper(pnds[i])
	}
	return planNodeDescriptions
}

type defaultPair struct {
	*graph.Pair
}

func pairWrapper(p *graph.Pair) types.Pair {
	return defaultPair{p}
}

func pairsWrapper(ps []*graph.Pair) []types.Pair {
	pairs := make([]types.Pair, len(ps))
	for i := range ps {
		pairs[i] = pairWrapper(ps[i])
	}
	return pairs
}

type defaultProfilingStats struct {
	*graph.ProfilingStats
}

func profilingStatsWrapper(ps *graph.ProfilingStats) types.ProfilingStats {
	return defaultProfilingStats{ps}
}

func profilingStatssWrapper(pss []*graph.ProfilingStats) []types.ProfilingStats {
	profilingStatss := make([]types.ProfilingStats, len(pss))
	for i := range pss {
		profilingStatss[i] = profilingStatsWrapper(pss[i])
	}
	return profilingStatss
}

type defaultPlanNodeBranchInfo struct {
	*graph.PlanNodeBranchInfo
}

func planNodeBranchInfoWrapper(pnb *graph.PlanNodeBranchInfo) types.PlanNodeBranchInfo {
	return defaultPlanNodeBranchInfo{pnb}
}
