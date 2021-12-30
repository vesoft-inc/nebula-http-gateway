package v3_0

import (
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
	return w.AuthResponse.SessionID
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

func (w executionResponseWrapper) GetSpaceName() []byte {
	return w.ExecutionResponse.GetSpaceName()
}

func (w executionResponseWrapper) GetPlanDesc() types.PlanDescription {
	return newPlanDescriptionWrapper(w.ExecutionResponse.PlanDesc)
}

func (w executionResponseWrapper) GetComment() []byte {
	return w.ExecutionResponse.GetComment()
}

func (w executionResponseWrapper) IsSetData() bool {
	return w.ExecutionResponse.IsSetData()
}

func (w executionResponseWrapper) IsSetSpaceName() bool {
	return w.ExecutionResponse.IsSetSpaceName()
}

func (w executionResponseWrapper) IsSetErrorMsg() bool {
	return w.ExecutionResponse.IsSetErrorMsg()
}

func (w executionResponseWrapper) IsSetPlanDesc() bool {
	return w.ExecutionResponse.IsSetPlanDesc()
}

func (w executionResponseWrapper) IsSetComment() bool {
	return w.ExecutionResponse.IsSetComment()
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
	return newRowsWrapper(w.Rows)
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
	return newVaulesWrapper(w.Values)
}

type valueWrapper struct {
	*nthrift.Value
}

func newValueWrapper(value *nthrift.Value) types.Value {
	if value == nil {
		return nil
	}
	return valueWrapper{value}
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

func (w valueWrapper) GetNVal() *types.NullType {
	return newNullTypeWrapper(w.NVal)
}

func (w valueWrapper) GetDVal() types.Date {
	return newDateWrapper(w.DVal)
}

func (w valueWrapper) GetTVal() types.Time {
	return newTimeWrapper(w.TVal)
}

func (w valueWrapper) GetDtVal() types.DateTime {
	return newDateTimeWrapper(w.DtVal)
}

func (w valueWrapper) GetVVal() types.Vertex {
	return newVertexWrapper(w.VVal)
}

func (w valueWrapper) GetEVal() types.Edge {
	return newEdgeWrapper(w.EVal)
}

func (w valueWrapper) GetPVal() types.Path {
	return newPathWrapper(w.PVal)
}

func (w valueWrapper) GetLVal() types.NList {
	return newNListWrapper(w.LVal)
}

func (w valueWrapper) GetMVal() types.NMap {
	return newNMapWrapper(w.MVal)
}

func (w valueWrapper) GetUVal() types.NSet {
	return newNSetWrapper(w.UVal)
}

func (w valueWrapper) GetGVal() types.DataSet {
	return newDataSetWrapper(w.GVal)
}

func (w valueWrapper) GetGgVal() types.Geography {
	return newGeographyWrapper(w.GgVal)
}

func (w valueWrapper) IsSetGgVal() bool {
	return false
}

func newNullTypeWrapper(nullType *nthrift.NullType) *types.NullType {
	if nullType == nil {
		return nil
	}
	return types.NullTypePtr(types.NullTypeToValue[nullType.String()])
}

type dateWrapper struct {
	*nthrift.Date
}

func newDateWrapper(date *nthrift.Date) types.Date {
	if date == nil {
		return nil
	}
	return dateWrapper{date}
}

type timeWrapper struct {
	*nthrift.Time
}

func newTimeWrapper(time *nthrift.Time) types.Time {
	if time == nil {
		return nil
	}
	return timeWrapper{time}
}

type dateTimeWrapper struct {
	*nthrift.DateTime
}

func newDateTimeWrapper(dateTime *nthrift.DateTime) types.DateTime {
	if dateTime == nil {
		return nil
	}
	return dateTimeWrapper{dateTime}
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
	return newValueWrapper(w.Vid)
}

func (w vertexWrapper) GetTags() []types.Tag {
	return newTagsWrapper(w.Tags)
}

type edgeWrapper struct {
	*nthrift.Edge
}

func newEdgeWrapper(edge *nthrift.Edge) types.Edge {
	if edge == nil {
		return nil
	}
	return edgeWrapper{edge}
}

func (w edgeWrapper) GetSrc() types.Value {
	value := newValueWrapper(w.Src)
	return value
}

func (w edgeWrapper) GetDst() types.Value {
	value := newValueWrapper(w.Dst)
	return value
}

func (w edgeWrapper) GetType() types.EdgeType {
	return newEdgeTypeWrapper(w.Type)
}

func (w edgeWrapper) GetRanking() types.EdgeRanking {
	return newEdgeRankingWrapper(w.Ranking)
}

func (w edgeWrapper) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(w.Props))
	for k, v := range w.Props {
		props[k] = newValueWrapper(v)
	}
	return props
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
	return newVertexWrapper(w.Src)
}
func (w pathWrapper) GetSteps() []types.Step {
	return newStepsWrapper(w.Steps)
}

type nListWrapper struct {
	*nthrift.NList
}

func (w nListWrapper) GetValues() []types.Value {
	return newVaulesWrapper(w.Values)
}

func newNListWrapper(nList *nthrift.NList) types.NList {
	if nList == nil {
		return nil
	}
	return nListWrapper{nList}
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
	kvs := make(map[string]types.Value, len(w.Kvs))
	for k, v := range w.Kvs {
		kvs[k] = newValueWrapper(v)
	}
	return kvs
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
	return newVaulesWrapper(w.Values)
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
	return newPointWrapper(w.PtVal)
}
func (w geographyWrapper) GetLsVal() types.LineString {
	return newLineStringWrapper(w.LsVal)
}
func (w geographyWrapper) GetPgVal() types.Polygon {
	return newPolygonWrapper(w.PgVal)
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
	props := make(map[string]types.Value, len(w.Props))
	for k, v := range w.Props {
		value := newValueWrapper(v)
		props[k] = value
	}
	return props
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
	return newVertexWrapper(w.Dst)
}

func (w stepWrapper) GetType() types.EdgeType {
	return newEdgeTypeWrapper(w.Type)
}

func (w stepWrapper) GetRanking() types.EdgeRanking {
	return newEdgeRankingWrapper(w.Ranking)
}

func (w stepWrapper) GetProps() map[string]types.Value {
	props := make(map[string]types.Value, len(w.Props))
	for k, v := range w.Props {
		props[k] = newValueWrapper(v)
	}
	return props
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
	return newCoordinateWrapper(w.Coord)
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
	return newCoordinatesWrapper(w.CoordList)
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
	return newCoordinatesSliceWrapper(w.CoordListList)
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
	return planNodeDescriptionsWrapper(w.PlanNodeDescs)
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

func (w planNodeDescriptionWrapper) GetDescription() []types.Pair {
	return newPairsWrapper(w.Description)
}

func (w planNodeDescriptionWrapper) GetProfiles() []types.ProfilingStats {
	return newProfilingStatssWrapper(w.Profiles)
}

func (w planNodeDescriptionWrapper) GetBranchInfo() types.PlanNodeBranchInfo {
	return newPlanNodeBranchInfoWrapper(w.BranchInfo)
}

func planNodeDescriptionsWrapper(planNodeDescriptions []*graph.PlanNodeDescription) []types.PlanNodeDescription {
	if planNodeDescriptions == nil {
		return nil
	}
	descriptions := make([]types.PlanNodeDescription, len(planNodeDescriptions))
	for i := range planNodeDescriptions {
		descriptions[i] = newPlanNodeDescriptionWrapper(planNodeDescriptions[i])
	}
	return descriptions
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

type planNodeBranchInfoWarpper struct {
	*graph.PlanNodeBranchInfo
}

func newPlanNodeBranchInfoWrapper(planNodeBranchInfo *graph.PlanNodeBranchInfo) types.PlanNodeBranchInfo {
	if planNodeBranchInfo == nil {
		return nil
	}
	return planNodeBranchInfoWarpper{planNodeBranchInfo}
}
