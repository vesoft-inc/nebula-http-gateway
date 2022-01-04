package v3_0

import (
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0/graph"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type authResponseWrapper struct {
	*graph.AuthResponse
}

func newAuthResponseWrapper(authResponse *graph.AuthResponse) types.AuthResponse {
	return authResponseWrapper{authResponse}
}

func (w authResponseWrapper) SessionID() *int64 {
	sid := w.AuthResponse.GetSessionID()
	return &sid
}

func (w authResponseWrapper) GetTimezoneInfo() types.TimezoneInfo {
	timezoneOffset := w.AuthResponse.GetTimeZoneOffsetSeconds()
	timezoneName := w.AuthResponse.GetTimeZoneName()

	timezoneInfo := types.TimezoneInfo{}
	timezoneInfo.SetOffset(timezoneOffset)
	timezoneInfo.SetName(timezoneName)
	return timezoneInfo
}

type executionResponseWrapper struct {
	*graph.ExecutionResponse
}

func newExecutionResponseWrapper(executionResponse *graph.ExecutionResponse) types.ExecutionResponse {
	return executionResponseWrapper{executionResponse}
}

func (w executionResponseWrapper) GetLatencyInUs() int64 {
	return int64(w.ExecutionResponse.GetLatencyInUs())
}

func (w executionResponseWrapper) GetData() types.DataSet {
	return newDataSetWrapper(w.ExecutionResponse.GetData())
}

func (w executionResponseWrapper) GetPlanDesc() types.PlanDescription {
	return newPlanDescriptionWrapper(w.ExecutionResponse.GetPlanDesc())
}

func (w executionResponseWrapper) GetErrorCode() nerrors.ErrorCode {
	return nerrors.ErrorCode(w.ExecutionResponse.GetErrorCode())
}

func (w executionResponseWrapper) String() string {
	return w.ExecutionResponse.String()
}

type dataSetWrapper struct {
	*nthrift.DataSet
}

func newDataSetWrapper(dataSet *nthrift.DataSet) types.DataSet {
	if dataSet == nil {
		return nil
	}
	return dataSetWrapper{dataSet}
}

func (w dataSetWrapper) GetRows() []types.Row {
	return newRowsWrapper(w.DataSet.Rows)
}

func (w dataSetWrapper) Unwrap() interface{} {
	return w.DataSet
}

type rowWrapper struct {
	*nthrift.Row
}

func newRowWrapper(row *nthrift.Row) types.Row {
	if row == nil {
		return nil
	}
	return rowWrapper{row}
}

func newRowsWrapper(rows []*nthrift.Row) []types.Row {
	if rows == nil {
		return nil
	}
	rs := make([]types.Row, len(rows))
	for i := range rows {
		rs[i] = newRowWrapper(rows[i])
	}
	return rs
}

func (w rowWrapper) GetValues() []types.Value {
	return newVaulesWrapper(w.Row.Values)
}

func (w rowWrapper) Unwrap() interface{} {
	return w.Row
}

type valueWrapper struct {
	*nthrift.Value
	builder *valueBuilder
}

func newValueWrapper(value *nthrift.Value) types.Value {
	if value == nil {
		return nil
	}

	v := *value
	return valueWrapper{Value: value, builder: &valueBuilder{v}}
}

func newVaulesWrapper(values []*nthrift.Value) []types.Value {
	if values == nil {
		return nil
	}
	vs := make([]types.Value, len(values))
	for i := range values {
		vs[i] = newValueWrapper(values[i])
	}
	return vs
}

func (w valueWrapper) GetNVal() types.NullType {
	return newNullTypeWrapper(w.Value.NVal)
}

func (w valueWrapper) GetDVal() types.Date {
	return newDateWrapper(w.Value.GetDVal())
}

func (w valueWrapper) GetTVal() types.Time {
	return newTimeWrapper(w.Value.GetTVal())
}

func (w valueWrapper) GetDtVal() types.DateTime {
	return newDateTimeWrapper(w.Value.GetDtVal())
}

func (w valueWrapper) GetVVal() types.Vertex {
	return newVertexWrapper(w.Value.GetVVal())
}

func (w valueWrapper) GetEVal() types.Edge {
	return newEdgeWrapper(w.Value.GetEVal())
}

func (w valueWrapper) GetPVal() types.Path {
	return newPathWrapper(w.Value.GetPVal())
}

func (w valueWrapper) GetLVal() types.NList {
	return newNListWrapper(w.Value.GetLVal())
}

func (w valueWrapper) GetMVal() types.NMap {
	return newNMapWrapper(w.Value.GetMVal())
}

func (w valueWrapper) GetUVal() types.NSet {
	return newNSetWrapper(w.Value.GetUVal())
}

func (w valueWrapper) GetGVal() types.DataSet {
	return newDataSetWrapper(w.Value.GetGVal())
}

func (w valueWrapper) GetGgVal() types.Geography {
	return newGeographyWrapper(w.Value.GetGgVal())
}

func (w valueWrapper) GetDuVal() types.Duration {
	return newDurationWrapper(w.Value.GetDuVal())
}

func (w valueWrapper) SetNVal(nval *types.NullType) types.Value {
	w.Value = w.Value.SetNVal((*nthrift.NullType)(nval))
	return w
}

func (w valueWrapper) SetBVal(bval *bool) types.Value {
	w.Value = w.Value.SetBVal(bval)
	return w
}

func (w valueWrapper) SetIVal(ival *int64) types.Value {
	w.Value = w.Value.SetIVal(ival)
	return w
}

func (w valueWrapper) SetFVal(fval *float64) types.Value {
	w.Value = w.Value.SetFVal(fval)
	return w
}

func (w valueWrapper) SetSVal(sval []byte) types.Value {
	w.Value = w.Value.SetSVal(sval)
	return w
}

func (w valueWrapper) SetDVal(dval types.Date) types.Value {
	w.Value = w.Value.SetDVal(dval.Unwrap().(*nthrift.Date))
	return w
}

func (w valueWrapper) SetTVal(tval types.Time) types.Value {
	w.Value = w.Value.SetTVal(tval.Unwrap().(*nthrift.Time))
	return w
}

func (w valueWrapper) SetDtVal(dtval types.DateTime) types.Value {
	w.Value = w.Value.SetDtVal(dtval.Unwrap().(*nthrift.DateTime))
	return w
}

func (w valueWrapper) SetVVal(vval types.Vertex) types.Value {
	w.Value = w.Value.SetVVal(vval.Unwrap().(*nthrift.Vertex))
	return w
}

func (w valueWrapper) SetEVal(eval types.Edge) types.Value {
	w.Value = w.Value.SetEVal(eval.Unwrap().(*nthrift.Edge))
	return w
}

func (w valueWrapper) SetPVal(pval types.Path) types.Value {
	w.Value = w.Value.SetPVal(pval.Unwrap().(*nthrift.Path))
	return w
}

func (w valueWrapper) SetLVal(lval types.NList) types.Value {
	w.Value = w.Value.SetLVal(lval.Unwrap().(*nthrift.NList))
	return w
}

func (w valueWrapper) SetMVal(mval types.NMap) types.Value {
	w.Value = w.Value.SetMVal(mval.Unwrap().(*nthrift.NMap))
	return w
}

func (w valueWrapper) SetUVal(uval types.NSet) types.Value {
	w.Value = w.Value.SetUVal(uval.Unwrap().(*nthrift.NSet))
	return w
}

func (w valueWrapper) SetGVal(gval types.DataSet) types.Value {
	w.Value = w.Value.SetGVal(gval.Unwrap().(*nthrift.DataSet))
	return w
}

func (w valueWrapper) SetGgVal(ggval types.Geography) types.Value {
	w.Value = w.Value.SetGgVal(ggval.Unwrap().(*nthrift.Geography))
	return w
}

func (w valueWrapper) SetDuVal(duval types.Duration) types.Value {
	w.Value = w.Value.SetDuVal(duval.Unwrap().(*nthrift.Duration))
	return w
}

func (w valueWrapper) Unwrap() interface{} {
	return w.Value
}

func (w valueWrapper) Builder() types.ValueBuilder {
	return w.builder
}

type valueBuilder struct {
	value nthrift.Value
}

func (b valueBuilder) NVal(nval *types.NullType) types.ValueBuilder {
	b.value.SetNVal((*nthrift.NullType)(nval))
	return b
}

func (b valueBuilder) BVal(bval *bool) types.ValueBuilder {
	b.value.SetBVal(bval)
	return b
}

func (b valueBuilder) IVal(ival *int64) types.ValueBuilder {
	b.value.SetIVal(ival)
	return b
}

func (b valueBuilder) FVal(fval *float64) types.ValueBuilder {
	b.value.SetFVal(fval)
	return b
}

func (b valueBuilder) SVal(sval []byte) types.ValueBuilder {
	b.value.SetSVal(sval)
	return b
}

func (b valueBuilder) DVal(dval types.Date) types.ValueBuilder {
	b.value.SetDVal(dval.Unwrap().(*nthrift.Date))
	return b
}

func (b valueBuilder) TVal(tval types.Time) types.ValueBuilder {
	b.value.SetTVal(tval.Unwrap().(*nthrift.Time))
	return b
}

func (b valueBuilder) DtVal(dtval types.DateTime) types.ValueBuilder {
	b.value.SetDtVal(dtval.Unwrap().(*nthrift.DateTime))
	return b
}

func (b valueBuilder) VVal(vval types.Vertex) types.ValueBuilder {
	b.value.SetVVal(vval.Unwrap().(*nthrift.Vertex))
	return b
}

func (b valueBuilder) EVal(eval types.Edge) types.ValueBuilder {
	b.value.SetEVal(eval.Unwrap().(*nthrift.Edge))
	return b
}

func (b valueBuilder) PVal(pval types.Path) types.ValueBuilder {
	b.value.SetPVal(pval.Unwrap().(*nthrift.Path))
	return b
}

func (b valueBuilder) LVal(lval types.NList) types.ValueBuilder {
	b.value.SetLVal(lval.Unwrap().(*nthrift.NList))
	return b
}

func (b valueBuilder) MVal(mval types.NMap) types.ValueBuilder {
	b.value.SetMVal(mval.Unwrap().(*nthrift.NMap))
	return b
}

func (b valueBuilder) UVal(uval types.NSet) types.ValueBuilder {
	b.value.SetUVal(uval.Unwrap().(*nthrift.NSet))
	return b
}

func (b valueBuilder) GVal(gval types.DataSet) types.ValueBuilder {
	b.value.SetGVal(gval.Unwrap().(*nthrift.DataSet))
	return b
}

func (b valueBuilder) GgVal(ggval types.Geography) types.ValueBuilder {
	b.value.SetGgVal(ggval.Unwrap().(*nthrift.Geography))
	return b
}

func (b valueBuilder) DuVal(duval types.Duration) types.ValueBuilder {
	b.value.SetDuVal(duval.Unwrap().(*nthrift.Duration))
	return b
}

func (b valueBuilder) Emit() types.Value {
	return newValueWrapper(&b.value)
}

func newNullTypeWrapper(nullType *nthrift.NullType) types.NullType {
	if nullType == nil {
		return -1
	}
	return types.NullTypeToValue[nullType.String()]
}

type dateWrapper struct {
	*nthrift.Date
	builder *dateBuilder
}

func newDateWrapper(date *nthrift.Date) types.Date {
	if date == nil {
		return nil
	}

	d := *date
	return dateWrapper{date, &dateBuilder{d}}
}

func (w dateWrapper) SetYear(year int16) types.Date {
	w.Date = w.Date.SetYear(year)
	return w
}

func (w dateWrapper) SetMonth(month int8) types.Date {
	w.Date = w.Date.SetMonth(month)
	return w
}

func (w dateWrapper) SetDay(day int8) types.Date {
	w.Date = w.Date.SetDay(day)
	return w
}

func (w dateWrapper) Unwrap() interface{} {
	return w.Date
}

func (w dateWrapper) Builder() types.DateBuilder {
	return w.builder
}

type dateBuilder struct {
	date nthrift.Date
}

func (b dateBuilder) Year(year int16) types.DateBuilder {
	b.date.SetYear(year)
	return b
}

func (b dateBuilder) Month(month int8) types.DateBuilder {
	b.date.SetMonth(month)
	return b
}

func (b dateBuilder) Day(day int8) types.DateBuilder {
	b.date.SetDay(day)
	return b
}

func (b dateBuilder) Emit() types.Date {
	return newDateWrapper(&b.date)
}

type timeWrapper struct {
	*nthrift.Time
	builder *timeBuilder
}

func newTimeWrapper(time *nthrift.Time) types.Time {
	if time == nil {
		return nil
	}

	t := *time
	return timeWrapper{time, &timeBuilder{t}}
}

func (w timeWrapper) SetHour(hour int8) types.Time {
	w.Time = w.Time.SetHour(hour)
	return w
}

func (w timeWrapper) SetMinute(minute int8) types.Time {
	w.Time = w.Time.SetMinute(minute)
	return w
}

func (w timeWrapper) SetSec(sec int8) types.Time {
	w.Time = w.Time.SetSec(sec)
	return w
}

func (w timeWrapper) SetMicrosec(microsec int32) types.Time {
	w.Time = w.Time.SetMicrosec(microsec)
	return w
}

func (w timeWrapper) Unwrap() interface{} {
	return w.Time
}

func (w timeWrapper) Builder() types.TimeBuilder {
	return w.builder
}

type timeBuilder struct {
	time nthrift.Time
}

func (b timeBuilder) Hour(hour int8) types.TimeBuilder {
	b.time.SetHour(hour)
	return b
}

func (b timeBuilder) Minute(minute int8) types.TimeBuilder {
	b.time.SetMinute(minute)
	return b
}

func (b timeBuilder) Sec(sec int8) types.TimeBuilder {
	b.time.SetSec(sec)
	return b
}

func (b timeBuilder) Microsec(microsec int32) types.TimeBuilder {
	b.time.SetMicrosec(microsec)
	return b
}

func (b timeBuilder) Emit() types.Time {
	return newTimeWrapper(&b.time)
}

type dateTimeWrapper struct {
	*nthrift.DateTime
	builder *dateTimeBuilder
}

func newDateTimeWrapper(dateTime *nthrift.DateTime) types.DateTime {
	if dateTime == nil {
		return nil
	}

	dt := *dateTime
	return dateTimeWrapper{dateTime, &dateTimeBuilder{dt}}
}

func (w dateTimeWrapper) SetYear(year int16) types.DateTime {
	w.DateTime = w.DateTime.SetYear(year)
	return w
}

func (w dateTimeWrapper) SetMonth(month int8) types.DateTime {
	w.DateTime = w.DateTime.SetMonth(month)
	return w
}

func (w dateTimeWrapper) SetDay(day int8) types.DateTime {
	w.DateTime = w.DateTime.SetDay(day)
	return w
}

func (w dateTimeWrapper) SetHour(hour int8) types.DateTime {
	w.DateTime = w.DateTime.SetHour(hour)
	return w
}

func (w dateTimeWrapper) SetMinute(minute int8) types.DateTime {
	w.DateTime = w.DateTime.SetMinute(minute)
	return w
}

func (w dateTimeWrapper) SetSec(sec int8) types.DateTime {
	w.DateTime = w.DateTime.SetSec(sec)
	return w
}

func (w dateTimeWrapper) SetMicrosec(microsec int32) types.DateTime {
	w.DateTime = w.DateTime.SetMicrosec(microsec)
	return w
}

func (w dateTimeWrapper) Unwrap() interface{} {
	return w.DateTime
}

func (w dateTimeWrapper) Builder() types.DateTimeBuilder {
	return w.builder
}

type dateTimeBuilder struct {
	dateTime nthrift.DateTime
}

func (b dateTimeBuilder) Year(year int16) types.DateTimeBuilder {
	b.dateTime.SetYear(year)
	return b
}

func (b dateTimeBuilder) Month(month int8) types.DateTimeBuilder {
	b.dateTime.SetMonth(month)
	return b
}

func (b dateTimeBuilder) Day(day int8) types.DateTimeBuilder {
	b.dateTime.SetDay(day)
	return b
}

func (b dateTimeBuilder) Hour(hour int8) types.DateTimeBuilder {
	b.dateTime.SetHour(hour)
	return b
}

func (b dateTimeBuilder) Minute(minute int8) types.DateTimeBuilder {
	b.dateTime.SetMinute(minute)
	return b
}

func (b dateTimeBuilder) Sec(sec int8) types.DateTimeBuilder {
	b.dateTime.SetSec(sec)
	return b
}

func (b dateTimeBuilder) Microsec(microsec int32) types.DateTimeBuilder {
	b.dateTime.SetMicrosec(microsec)
	return b
}

func (b dateTimeBuilder) Emit() types.DateTime {
	return newDateTimeWrapper(&b.dateTime)
}

type vertexWrapper struct {
	*nthrift.Vertex
}

func newVertexWrapper(vertex *nthrift.Vertex) types.Vertex {
	if vertex == nil {
		return nil
	}
	return vertexWrapper{vertex}
}

func (w vertexWrapper) GetVid() types.Value {
	return newValueWrapper(w.Vertex.GetVid())
}

func (w vertexWrapper) GetTags() []types.Tag {
	return newTagsWrapper(w.Vertex.GetTags())
}

func (w vertexWrapper) Unwrap() interface{} {
	return w.Vertex
}

type edgeWrapper struct {
	*nthrift.Edge
	builder *edgeBuilder
}

func newEdgeWrapper(edge *nthrift.Edge) types.Edge {
	if edge == nil {
		return nil
	}

	e := *edge
	return edgeWrapper{edge, &edgeBuilder{e}}
}

func (w edgeWrapper) GetSrc() types.Value {
	value := newValueWrapper(w.Edge.GetSrc())
	return value
}

func (w edgeWrapper) GetDst() types.Value {
	value := newValueWrapper(w.Edge.GetDst())
	return value
}

func (w edgeWrapper) GetType() types.EdgeType {
	return newEdgeTypeWrapper(w.Edge.GetType())
}

func (w edgeWrapper) GetRanking() types.EdgeRanking {
	return newEdgeRankingWrapper(w.Edge.GetRanking())
}

func (w edgeWrapper) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(w.Edge.GetProps()))
	for k, v := range w.Props {
		props[k] = newValueWrapper(v)
	}
	return props
}

func (w edgeWrapper) SetSrc(src types.Value) types.Edge {
	w.Edge = w.Edge.SetSrc(src.Unwrap().(*nthrift.Value))
	return w
}

func (w edgeWrapper) SetDst(dst types.Value) types.Edge {
	w.Edge = w.Edge.SetDst(dst.Unwrap().(*nthrift.Value))
	return w
}

func (w edgeWrapper) SetType(edgeType types.EdgeType) types.Edge {
	w.Edge = w.Edge.SetType(edgeType)
	return w
}

func (w edgeWrapper) SetName(name []byte) types.Edge {
	w.Edge = w.Edge.SetName(name)
	return w
}

func (w edgeWrapper) SetRanking(edgeRanking types.EdgeRanking) types.Edge {
	w.Edge = w.Edge.SetRanking(edgeRanking)
	return w
}

func (w edgeWrapper) SetProps(props map[string]types.Value) types.Edge {
	_props := make(map[string]*nthrift.Value, len(props))
	for k, v := range props {
		_props[k] = v.Unwrap().(*nthrift.Value)
	}
	w.Edge = w.Edge.SetProps(_props)
	return w
}

func (w edgeWrapper) Unwrap() interface{} {
	return w.Edge
}

func (w edgeWrapper) Builder() types.EdgeBuilder {
	return w.builder
}

type edgeBuilder struct {
	edge nthrift.Edge
}

func (b edgeBuilder) Src(src types.Value) types.EdgeBuilder {
	b.edge.SetSrc(src.Unwrap().(*nthrift.Value))
	return b
}

func (b edgeBuilder) Dst(dst types.Value) types.EdgeBuilder {
	b.edge.SetDst(dst.Unwrap().(*nthrift.Value))
	return b
}

func (b edgeBuilder) Type(edgeType types.EdgeType) types.EdgeBuilder {
	b.edge.SetType(edgeType)
	return b
}

func (b edgeBuilder) Name(name []byte) types.EdgeBuilder {
	b.edge.SetName(name)
	return b
}

func (b edgeBuilder) Ranking(edgeRanking types.EdgeRanking) types.EdgeBuilder {
	b.edge.SetRanking(edgeRanking)
	return b
}

func (b edgeBuilder) Props(props map[string]types.Value) types.EdgeBuilder {
	_props := make(map[string]*nthrift.Value, len(props))
	for k, v := range props {
		_props[k] = v.Unwrap().(*nthrift.Value)
	}
	b.edge.SetProps(_props)
	return b
}

func (b edgeBuilder) Emit() types.Edge {
	return newEdgeWrapper(&b.edge)
}

func newEdgeTypeWrapper(edgeType nthrift.EdgeType) types.EdgeType {
	return edgeType
}

func newEdgeRankingWrapper(edgeRanking nthrift.EdgeRanking) types.EdgeRanking {
	return edgeRanking
}

type pathWrapper struct {
	*nthrift.Path
}

func newPathWrapper(path *nthrift.Path) types.Path {
	if path == nil {
		return nil
	}
	return pathWrapper{path}
}

func (w pathWrapper) GetSrc() types.Vertex {
	return newVertexWrapper(w.Path.GetSrc())
}
func (w pathWrapper) GetSteps() []types.Step {
	return newStepsWrapper(w.Path.GetSteps())
}

func (w pathWrapper) Unwrap() interface{} {
	return w.Path
}

type nListWrapper struct {
	*nthrift.NList
}

func (w nListWrapper) GetValues() []types.Value {
	return newVaulesWrapper(w.NList.GetValues())
}

func newNListWrapper(nList *nthrift.NList) types.NList {
	if nList == nil {
		return nil
	}
	return nListWrapper{nList}
}

func (w nListWrapper) Unwrap() interface{} {
	return w.NList
}

type nMapWrapper struct {
	*nthrift.NMap
}

func newNMapWrapper(nMap *nthrift.NMap) types.NMap {
	if nMap == nil {
		return nil
	}
	return nMapWrapper{nMap}
}

func (w nMapWrapper) GetKvs() map[string]types.Value {
	kvs := make(map[string]types.Value, len(w.NMap.GetKvs()))
	for k, v := range w.Kvs {
		kvs[k] = newValueWrapper(v)
	}
	return kvs
}

func (w nMapWrapper) Unwrap() interface{} {
	return w.NMap
}

type nSetWraooer struct {
	*nthrift.NSet
}

func newNSetWrapper(nSet *nthrift.NSet) types.NSet {
	if nSet == nil {
		return nil
	}
	return nSetWraooer{nSet}
}

func (w nSetWraooer) GetValues() []types.Value {
	return newVaulesWrapper(w.NSet.GetValues())
}

func (w nSetWraooer) Unwrap() interface{} {
	return w.NSet
}

type geographyWrapper struct {
	*nthrift.Geography
}

func newGeographyWrapper(geography *nthrift.Geography) types.Geography {
	if geography == nil {
		return nil
	}
	return geographyWrapper{geography}
}

func (w geographyWrapper) GetPtVal() types.Point {
	return newPointWrapper(w.Geography.GetPtVal())
}
func (w geographyWrapper) GetLsVal() types.LineString {
	return newLineStringWrapper(w.Geography.GetLsVal())
}
func (w geographyWrapper) GetPgVal() types.Polygon {
	return newPolygonWrapper(w.Geography.GetPgVal())
}

func (w geographyWrapper) Unwrap() interface{} {
	return w.Geography
}

type tagWrapper struct {
	*nthrift.Tag
}

func newTagWrapper(tag *nthrift.Tag) types.Tag {
	if tag == nil {
		return nil
	}
	return tagWrapper{tag}
}

func newTagsWrapper(tags []*nthrift.Tag) []types.Tag {
	if tags == nil {
		return nil
	}
	ts := make([]types.Tag, len(tags))
	for i := range tags {
		ts[i] = newTagWrapper(tags[i])
	}
	return ts
}

func (w tagWrapper) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(w.Tag.GetProps()))
	for k, v := range w.Props {
		value := newValueWrapper(v)
		props[k] = value
	}
	return props
}

func (w tagWrapper) Unwrap() interface{} {
	return w.Tag
}

type stepWrapper struct {
	*nthrift.Step
}

func newStepWrapper(step *nthrift.Step) types.Step {
	if step == nil {
		return nil
	}
	return stepWrapper{step}
}

func newStepsWrapper(steps []*nthrift.Step) []types.Step {
	if steps == nil {
		return nil
	}
	ss := make([]types.Step, len(steps))
	for i := range steps {
		ss[i] = newStepWrapper(steps[i])
	}
	return ss
}

func (w stepWrapper) GetDst() types.Vertex {
	return newVertexWrapper(w.Step.GetDst())
}

func (w stepWrapper) GetType() types.EdgeType {
	return newEdgeTypeWrapper(w.Step.GetType())
}

func (w stepWrapper) GetRanking() types.EdgeRanking {
	return newEdgeRankingWrapper(w.Step.GetRanking())
}

func (w stepWrapper) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(w.Step.GetProps()))
	for k, v := range w.Props {
		props[k] = newValueWrapper(v)
	}
	return props
}

func (w stepWrapper) Unwrap() interface{} {
	return w.Step
}

type pointWrapper struct {
	*nthrift.Point
}

func newPointWrapper(point *nthrift.Point) types.Point {
	if point == nil {
		return nil
	}
	return pointWrapper{point}
}

func (w pointWrapper) GetCoord() types.Coordinate {
	return newCoordinateWrapper(w.Point.GetCoord())
}

func (w pointWrapper) Unwrap() interface{} {
	return w.Point
}

type lineStringWrapper struct {
	*nthrift.LineString
}

func newLineStringWrapper(lineString *nthrift.LineString) types.LineString {
	if lineString == nil {
		return nil
	}
	return lineStringWrapper{lineString}
}

func (w lineStringWrapper) GetCoordList() []types.Coordinate {
	return newCoordinatesWrapper(w.LineString.GetCoordList())
}

func (w lineStringWrapper) Unwrap() interface{} {
	return w.LineString
}

type polygonWrapper struct {
	*nthrift.Polygon
}

func newPolygonWrapper(polygon *nthrift.Polygon) types.Polygon {
	if polygon == nil {
		return nil
	}
	return polygonWrapper{polygon}
}

func (w polygonWrapper) GetCoordListList() [][]types.Coordinate {
	return newCoordinatesSliceWrapper(w.Polygon.GetCoordListList())
}

func (w polygonWrapper) Unwrap() interface{} {
	return w.Polygon
}

type coordinateWrapper struct {
	*nthrift.Coordinate
}

func newCoordinateWrapper(coordinate *nthrift.Coordinate) types.Coordinate {
	if coordinate == nil {
		return nil
	}
	return coordinateWrapper{coordinate}
}

func newCoordinatesWrapper(cs []*nthrift.Coordinate) []types.Coordinate {
	if cs == nil {
		return nil
	}
	coords := make([]types.Coordinate, len(cs))
	for i := range cs {
		coords[i] = newCoordinateWrapper(cs[i])
	}
	return coords
}

func newCoordinatesSliceWrapper(coordinatesSlice [][]*nthrift.Coordinate) [][]types.Coordinate {
	if coordinatesSlice == nil {
		return nil
	}
	coordsSlice := make([][]types.Coordinate, len(coordinatesSlice))
	for i := range coordinatesSlice {
		coordsSlice[i] = newCoordinatesWrapper(coordinatesSlice[i])
	}
	return coordsSlice
}

func (w coordinateWrapper) Unwrap() interface{} {
	return w.Coordinate
}

type durationWrapper struct {
	*nthrift.Duration
}

func newDurationWrapper(duration *nthrift.Duration) types.Duration {
	if duration == nil {
		return nil
	}
	return durationWrapper{duration}
}

func (w durationWrapper) Unwrap() interface{} {
	return w.Duration
}

type planDescriptionWrapper struct {
	*graph.PlanDescription
}

func newPlanDescriptionWrapper(planDescription *graph.PlanDescription) types.PlanDescription {
	if planDescription == nil {
		return nil
	}
	return planDescriptionWrapper{planDescription}
}

func (w planDescriptionWrapper) GetPlanNodeDescs() []types.PlanNodeDescription {
	return newPlanNodeDescriptionsWrapper(w.PlanDescription.GetPlanNodeDescs())
}

func (w planDescriptionWrapper) Unwrap() interface{} {
	return w.PlanDescription
}

type planNodeDescriptionWrapper struct {
	*graph.PlanNodeDescription
}

func newPlanNodeDescriptionWrapper(planNodeDescription *graph.PlanNodeDescription) types.PlanNodeDescription {
	if planNodeDescription == nil {
		return nil
	}
	return planNodeDescriptionWrapper{planNodeDescription}
}

func newPlanNodeDescriptionsWrapper(planNodeDescriptions []*graph.PlanNodeDescription) []types.PlanNodeDescription {
	if planNodeDescriptions == nil {
		return nil
	}
	descriptions := make([]types.PlanNodeDescription, len(planNodeDescriptions))
	for i := range planNodeDescriptions {
		descriptions[i] = newPlanNodeDescriptionWrapper(planNodeDescriptions[i])
	}
	return descriptions
}

func (w planNodeDescriptionWrapper) GetDescription() []types.Pair {
	return newPairsWrapper(w.PlanNodeDescription.GetDescription())
}

func (w planNodeDescriptionWrapper) GetProfiles() []types.ProfilingStats {
	return newProfilingStatssWrapper(w.PlanNodeDescription.GetProfiles())
}

func (w planNodeDescriptionWrapper) GetBranchInfo() types.PlanNodeBranchInfo {
	return newPlanNodeBranchInfoWrapper(w.PlanNodeDescription.GetBranchInfo())
}

func (w planNodeDescriptionWrapper) Unwrap() interface{} {
	return w.PlanNodeDescription
}

type pairWrapper struct {
	*graph.Pair
}

func newPairWrapper(pair *graph.Pair) types.Pair {
	if pair == nil {
		return nil
	}
	return pairWrapper{pair}
}

func newPairsWrapper(pairs []*graph.Pair) []types.Pair {
	if pairs == nil {
		return nil
	}
	ps := make([]types.Pair, len(pairs))
	for i := range pairs {
		ps[i] = newPairWrapper(pairs[i])
	}
	return ps
}

func (w pairWrapper) Unwrap() interface{} {
	return w.Pair
}

type profilingStatsWrapper struct {
	*graph.ProfilingStats
}

func newProfilingStatsWrapper(profilingStats *graph.ProfilingStats) types.ProfilingStats {
	if profilingStats == nil {
		return nil
	}
	return profilingStatsWrapper{profilingStats}
}

func newProfilingStatssWrapper(profilingStatsSlice []*graph.ProfilingStats) []types.ProfilingStats {
	if profilingStatsSlice == nil {
		return nil
	}
	statsSlice := make([]types.ProfilingStats, len(profilingStatsSlice))
	for i := range profilingStatsSlice {
		statsSlice[i] = newProfilingStatsWrapper(profilingStatsSlice[i])
	}
	return statsSlice
}

func (w profilingStatsWrapper) Unwrap() interface{} {
	return w.ProfilingStats
}

type planNodeBranchInfoWrapper struct {
	*graph.PlanNodeBranchInfo
}

func newPlanNodeBranchInfoWrapper(planNodeBranchInfo *graph.PlanNodeBranchInfo) types.PlanNodeBranchInfo {
	if planNodeBranchInfo == nil {
		return nil
	}
	return planNodeBranchInfoWrapper{planNodeBranchInfo}
}

func (w planNodeBranchInfoWrapper) Unwrap() interface{} {
	return w.PlanNodeBranchInfo
}
