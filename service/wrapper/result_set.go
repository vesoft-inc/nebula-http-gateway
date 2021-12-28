package wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	cErrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
	"sort"
	"strings"
)

/*
	copy from nebula-go/result_set.go, and fit with ccore/nebula/types
*/

type ResultSet struct {
	resp            types.ExecutionResponse
	columnNames     []string
	colNameIndexMap map[string]int
}

type Record struct {
	columnNames     *[]string
	_record         []*ValueWrapper
	colNameIndexMap *map[string]int
}

type Node struct {
	vertex          types.Vertex
	tags            []string // tag name
	tagNameIndexMap map[string]int
}

type Relationship struct {
	edge types.Edge
}

type segment struct {
	startNode    *Node
	relationship *Relationship
	endNode      *Node
}

type PathWrapper struct {
	path             types.Path
	nodeList         []*Node
	relationshipList []*Relationship
	segments         []segment
}

type TimeWrapper struct {
	time types.Time
}

type DateWrapper struct {
	date types.Date
}

type DateTimeWrapper struct {
	dateTime types.DateTime
}

func GenResultSet(resp types.ExecutionResponse) (*ResultSet, error) {
	var colNames []string
	var colNameIndexMap = make(map[string]int)

	if resp.GetData() == nil { // if resp.Data != nil then resp.Data.row and resp.Data.colNames wont be nil
		return &ResultSet{
			resp:            resp,
			columnNames:     colNames,
			colNameIndexMap: colNameIndexMap,
		}, nil
	}
	for i, name := range resp.GetData().GetColumnNames() {
		colNames = append(colNames, string(name))
		colNameIndexMap[string(name)] = i
	}

	return &ResultSet{
		resp:            resp,
		columnNames:     colNames,
		colNameIndexMap: colNameIndexMap,
	}, nil
}

// Returns a 2D array of strings representing the query result
// If resultSet.resp.data is nil, returns an empty 2D array
func (res ResultSet) AsStringTable() [][]string {
	var resTable [][]string
	colNames := res.GetColNames()
	resTable = append(resTable, colNames)
	rows := res.GetRows()
	for _, row := range rows {
		var tempRow []string
		for _, val := range row.GetValues() {
			tempRow = append(tempRow, ValueWrapper{val}.String())
		}
		resTable = append(resTable, tempRow)
	}
	return resTable
}

// Returns all values in the given column
func (res ResultSet) GetValuesByColName(colName string) ([]*ValueWrapper, error) {
	if !res.hasColName(colName) {
		return nil, fmt.Errorf("failed to get values, given column name '%s' does not exist", colName)
	}
	// Get index
	index := res.colNameIndexMap[colName]
	var valList []*ValueWrapper
	for _, row := range res.resp.GetData().GetRows() {
		valList = append(valList, &ValueWrapper{row.GetValues()[index]})
	}
	return valList, nil
}

// Returns all values in the row at given index
func (res ResultSet) GetRowValuesByIndex(index int) (*Record, error) {
	if err := checkIndex(index, res.resp.GetData().GetRows()); err != nil {
		return nil, err
	}
	valWrap, err := GenValWraps(res.resp.GetData().GetRows()[index])
	if err != nil {
		return nil, err
	}
	return &Record{
		columnNames:     &res.columnNames,
		_record:         valWrap,
		colNameIndexMap: &res.colNameIndexMap,
	}, nil
}

// Returns the number of total rows
func (res ResultSet) GetRowSize() int {
	if res.resp.GetData() == nil {
		return 0
	}
	return len(res.resp.GetData().GetRows())
}

// Returns the number of total columns
func (res ResultSet) GetColSize() int {
	if res.resp.GetData() == nil {
		return 0
	}
	return len(res.resp.GetData().GetColumnNames())
}

// Returns all rows
func (res ResultSet) GetRows() []types.Row {
	if res.resp.GetData() == nil {
		var empty []types.Row
		return empty
	}
	return res.resp.GetData().GetRows()
}

func (res ResultSet) GetColNames() []string {
	return res.columnNames
}

// Returns an integer representing an error type
// 0    ErrorCode_SUCCEEDED
// -1   ErrorCode_E_DISCONNECTED
// -2   ErrorCode_E_FAIL_TO_CONNECT
// -3   ErrorCode_E_RPC_FAILURE
// -4   ErrorCode_E_BAD_USERNAME_PASSWORD
// -5   ErrorCode_E_SESSION_INVALID
// -6   ErrorCode_E_SESSION_TIMEOUT
// -7   ErrorCode_E_SYNTAX_ERROR
// -8   ErrorCode_E_EXECUTION_ERROR
// -9   ErrorCode_E_STATEMENT_EMPTY
// -10  ErrorCode_E_USER_NOT_FOUND
// -11  ErrorCode_E_BAD_PERMISSION
// -12  ErrorCode_E_SEMANTIC_ERROR
func (res ResultSet) GetErrorCode() cErrors.ErrorCode {
	return res.resp.GetErrorCode()
}

func (res ResultSet) GetLatency() int64 {
	return res.resp.GetLatencyInUs()
}

func (res ResultSet) GetSpaceName() string {
	if res.resp.GetSpaceName() == nil {
		return ""
	}
	return string(res.resp.GetSpaceName())
}

func (res ResultSet) GetErrorMsg() string {
	if res.resp.GetErrorMsg() == nil {
		return ""
	}
	return string(res.resp.GetErrorMsg())
}

func (res ResultSet) IsSetPlanDesc() bool {
	return res.resp.GetPlanDesc() != nil
}

func (res ResultSet) GetPlanDesc() types.PlanDescription {
	return res.resp.GetPlanDesc()
}

func (res ResultSet) IsSetComment() bool {
	return res.resp.GetComment() != nil
}

func (res ResultSet) GetComment() string {
	if res.resp.GetComment() == nil {
		return ""
	}
	return string(res.resp.GetComment())
}

func (res ResultSet) IsSetData() bool {
	return res.resp.GetData() != nil
}

func (res ResultSet) IsEmpty() bool {
	if !res.IsSetData() || len(res.resp.GetData().GetRows()) == 0 {
		return true
	}
	return false
}

func (res ResultSet) IsSucceed() bool {
	return res.GetErrorCode() == cErrors.ErrorCode_SUCCEEDED
}

func (res ResultSet) IsPartialSucceed() bool {
	return res.GetErrorCode() == cErrors.ErrorCode_E_PARTIAL_SUCCEEDED
}

func (res ResultSet) hasColName(colName string) bool {
	if _, ok := res.colNameIndexMap[colName]; ok {
		return true
	}
	return false
}

// explain/profile format="dot"
func (res ResultSet) MakeDotGraph() string {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var builder strings.Builder
	builder.WriteString("digraph exec_plan {\n")
	builder.WriteString("\trankdir=BT;\n")
	for _, planNodeDesc := range planNodeDescs {
		planNodeName := name(planNodeDesc)
		switch strings.ToLower(string(planNodeDesc.GetName())) {
		case "select":
			builder.WriteString(conditionalNodeString(planNodeName))
			dep := nodeById(p, planNodeDesc.GetDependencies()[0])
			// then branch
			thenNodeId := findBranchEndNode(p, planNodeDesc.GetId(), true)
			builder.WriteString(edgeString(name(nodeById(p, thenNodeId)), name(dep)))
			thenStartId := findFirstStartNodeFrom(p, thenNodeId)
			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, thenStartId)), "Y"))
			// else branch
			elseNodeId := findBranchEndNode(p, planNodeDesc.GetId(), false)
			builder.WriteString(edgeString(name(nodeById(p, elseNodeId)), name(dep)))
			elseStartId := findFirstStartNodeFrom(p, elseNodeId)
			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, elseStartId)), "N"))
			// dep
			builder.WriteString(edgeString(name(dep), planNodeName))
		case "loop":
			builder.WriteString(conditionalNodeString(planNodeName))
			dep := nodeById(p, planNodeDesc.GetDependencies()[0])
			// do branch
			doNodeId := findBranchEndNode(p, planNodeDesc.GetId(), true)
			builder.WriteString(edgeString(name(nodeById(p, doNodeId)), name(planNodeDesc)))
			doStartId := findFirstStartNodeFrom(p, doNodeId)
			builder.WriteString(conditionalEdgeString(name(planNodeDesc), name(nodeById(p, doStartId)), "Do"))
			// dep
			builder.WriteString(edgeString(name(dep), planNodeName))
		default:
			builder.WriteString(nodeString(planNodeDesc, planNodeName))
			if planNodeDesc.IsSetDependencies() {
				for _, depId := range planNodeDesc.GetDependencies() {
					builder.WriteString(edgeString(name(nodeById(p, depId)), planNodeName))
				}
			}
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// explain/profile format="dot:struct"
func (res ResultSet) MakeDotGraphByStruct() string {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var builder strings.Builder
	builder.WriteString("digraph exec_plan {\n")
	builder.WriteString("\trankdir=BT;\n")
	for _, planNodeDesc := range planNodeDescs {
		planNodeName := name(planNodeDesc)
		switch strings.ToLower(string(planNodeDesc.GetName())) {
		case "select":
			builder.WriteString(conditionalNodeString(planNodeName))
		case "loop":
			builder.WriteString(conditionalNodeString(planNodeName))
		default:
			builder.WriteString(nodeString(planNodeDesc, planNodeName))
		}

		if planNodeDesc.IsSetDependencies() {
			for _, depId := range planNodeDesc.GetDependencies() {
				dep := nodeById(p, depId)
				builder.WriteString(edgeString(name(dep), planNodeName))
			}
		}

		if planNodeDesc.IsSetBranchInfo() {
			branchInfo := planNodeDesc.GetBranchInfo()
			condNode := nodeById(p, branchInfo.GetConditionNodeID())
			label := condEdgeLabel(condNode, branchInfo.GetIsDoBranch())
			builder.WriteString(conditionalEdgeString(planNodeName, name(condNode), label))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// explain/profile format="row"
func (res ResultSet) MakePlanByRow() [][]interface{} {
	p := res.GetPlanDesc()
	planNodeDescs := p.GetPlanNodeDescs()
	var rows [][]interface{}
	for _, planNodeDesc := range planNodeDescs {
		var row []interface{}
		row = append(row, planNodeDesc.GetId(), string(planNodeDesc.GetName()))

		if planNodeDesc.IsSetDependencies() {
			var deps []string
			for _, dep := range planNodeDesc.GetDependencies() {
				deps = append(deps, fmt.Sprintf("%d", dep))
			}
			row = append(row, strings.Join(deps, ","))
		} else {
			row = append(row, "")
		}

		if planNodeDesc.IsSetProfiles() {
			var strArr []string
			for i, profile := range planNodeDesc.GetProfiles() {
				otherStats := profile.GetOtherStats()
				if otherStats != nil {
					strArr = append(strArr, "{")
				}
				s := fmt.Sprintf("ver: %d, rows: %d, execTime: %dus, totalTime: %dus",
					i, profile.GetRows(), profile.GetExecDurationInUs(), profile.GetTotalDurationInUs())
				strArr = append(strArr, s)

				for k, v := range otherStats {
					strArr = append(strArr, fmt.Sprintf("%s: %s", k, v))
				}
				if otherStats != nil {
					strArr = append(strArr, "}")
				}
			}
			row = append(row, strings.Join(strArr, "\n"))
		} else {
			row = append(row, "")
		}

		var columnInfo []string
		if planNodeDesc.IsSetBranchInfo() {
			branchInfo := planNodeDesc.GetBranchInfo()
			columnInfo = append(columnInfo, fmt.Sprintf("branch: %t, nodeId: %d\n",
				branchInfo.GetIsDoBranch(), branchInfo.GetConditionNodeID()))
		}

		outputVar := fmt.Sprintf("outputVar: %s", prettyFormatJsonString(planNodeDesc.GetOutputVar()))
		columnInfo = append(columnInfo, outputVar)

		if planNodeDesc.IsSetDescription() {
			desc := planNodeDesc.GetDescription()
			for _, pair := range desc {
				value := prettyFormatJsonString(pair.GetValue())
				columnInfo = append(columnInfo, fmt.Sprintf("%s: %s", string(pair.GetKey()), value))
			}
		}
		row = append(row, strings.Join(columnInfo, "\n"))
		rows = append(rows, row)
	}
	return rows
}

// Returns value in the record at given column index
func (record Record) GetValueByIndex(index int) (*ValueWrapper, error) {
	if err := checkIndex(index, record._record); err != nil {
		return nil, err
	}
	return record._record[index], nil
}

// Returns value in the record at given column name
func (record Record) GetValueByColName(colName string) (*ValueWrapper, error) {
	if !record.hasColName(colName) {
		return nil, fmt.Errorf("failed to get values, given column name '%s' does not exist", colName)
	}
	// Get index
	index := (*record.colNameIndexMap)[colName]
	return record._record[index], nil
}

func (record Record) String() string {
	var strList []string
	for _, val := range record._record {
		strList = append(strList, val.String())
	}
	return strings.Join(strList, ", ")
}

func (record Record) hasColName(colName string) bool {
	if _, ok := (*record.colNameIndexMap)[colName]; ok {
		return true
	}
	return false
}

func GenNode(vertex types.Vertex) (*Node, error) {
	if vertex == nil {
		return nil, fmt.Errorf("failed to generate Node: invalid vertex")
	}
	var tags []string
	nameIndex := make(map[string]int)

	// Iterate through all tags of the vertex
	for i, tag := range vertex.GetTags() {
		name := string(tag.GetName())
		// Get tags
		tags = append(tags, name)
		nameIndex[name] = i
	}

	return &Node{
		vertex:          vertex,
		tags:            tags,
		tagNameIndexMap: nameIndex,
	}, nil
}

// getRawID returns a list of row vid
func (node Node) getRawID() types.Value {
	return node.vertex.GetVid()
}

// GetID returns a list of vid of node
func (node Node) GetID() ValueWrapper {
	return ValueWrapper{node.vertex.GetVid()}
}

// GetTags returns a list of tag names of node
func (node Node) GetTags() []string {
	return node.tags
}

// HasTag checks if node contains given label
func (node Node) HasTag(label string) bool {
	if _, ok := node.tagNameIndexMap[label]; ok {
		return true
	}
	return false
}

// Properties returns all properties of a tag
func (node Node) Properties(tagName string) (map[string]*ValueWrapper, error) {
	kvMap := make(map[string]*ValueWrapper)
	// Check if label exists
	if !node.HasTag(tagName) {
		return nil, fmt.Errorf("failed to get properties: Tag name %s does not exsist in the Node", tagName)
	}
	index := node.tagNameIndexMap[tagName]
	for k, v := range node.vertex.GetTags()[index].GetProps() {
		kvMap[k] = &ValueWrapper{v}
	}
	return kvMap, nil
}

// Keys returns all prop names of the given tag name
func (node Node) Keys(tagName string) ([]string, error) {
	if !node.HasTag(tagName) {
		return nil, fmt.Errorf("failed to get properties: Tag name %s does not exsist in the Node", tagName)
	}
	var propNameList []string
	index := node.tagNameIndexMap[tagName]
	for k := range node.vertex.GetTags()[index].GetProps() {
		propNameList = append(propNameList, k)
	}
	return propNameList, nil
}

// Values returns all prop values of the given tag name
func (node Node) Values(tagName string) ([]*ValueWrapper, error) {
	if !node.HasTag(tagName) {
		return nil, fmt.Errorf("failed to get properties: Tag name %s does not exsist in the Node", tagName)
	}
	var propValList []*ValueWrapper
	index := node.tagNameIndexMap[tagName]
	for _, v := range node.vertex.GetTags()[index].GetProps() {
		propValList = append(propValList, &ValueWrapper{v})
	}
	return propValList, nil
}

// String returns a string representing node
// Node format: ("VertexID" :tag1{k0: v0,k1: v1}:tag2{k2: v2})
func (node Node) String() string {
	var keyList []string
	var kvStr []string
	var tagStr []string
	vertex := node.vertex
	vid := vertex.GetVid()
	for _, tag := range vertex.GetTags() {
		kvs := tag.GetProps()
		tagName := tag.GetName()
		for k := range kvs {
			keyList = append(keyList, k)
		}
		sort.Strings(keyList)
		for _, k := range keyList {
			kvTemp := fmt.Sprintf("%s: %s", k, ValueWrapper{kvs[k]}.String())
			kvStr = append(kvStr, kvTemp)
		}
		tagStr = append(tagStr, fmt.Sprintf("%s{%s}", tagName, strings.Join(kvStr, ", ")))
		keyList = nil
		kvStr = nil
	}
	if len(tagStr) == 0 { // No tag
		return fmt.Sprintf("(%s)", ValueWrapper{vid}.String())
	}
	return fmt.Sprintf("(%s :%s)", ValueWrapper{vid}.String(), strings.Join(tagStr, " :"))
}

// Returns true if two nodes have same vid
func (n1 Node) IsEqualTo(n2 *Node) bool {
	if n1.GetID().IsString() && n2.GetID().IsString() {
		s1, _ := n1.GetID().AsString()
		s2, _ := n2.GetID().AsString()
		return s1 == s2
	} else if n1.GetID().IsInt() && n2.GetID().IsInt() {
		s1, _ := n1.GetID().AsInt()
		s2, _ := n2.GetID().AsInt()
		return s1 == s2
	}
	return false
}

func GenRelationship(edge types.Edge) (*Relationship, error) {
	if edge == nil {
		return nil, fmt.Errorf("failed to generate Relationship: invalid edge")
	}
	return &Relationship{
		edge: edge,
	}, nil
}

func (relationship Relationship) GetSrcVertexID() ValueWrapper {
	if relationship.edge.GetType() > 0 {
		return ValueWrapper{relationship.edge.GetSrc()}
	}
	return ValueWrapper{relationship.edge.GetDst()}
}

func (relationship Relationship) GetDstVertexID() ValueWrapper {
	if relationship.edge.GetType() > 0 {
		return ValueWrapper{relationship.edge.GetDst()}
	}
	return ValueWrapper{relationship.edge.GetSrc()}
}

func (relationship Relationship) GetEdgeName() string {
	return string(relationship.edge.GetName())
}

func (relationship Relationship) GetRanking() int64 {
	return int64(relationship.edge.GetRanking())
}

// Properties returns a map where the key is property name and the value is property name
func (relationship Relationship) Properties() map[string]*ValueWrapper {
	kvMap := make(map[string]*ValueWrapper)
	var (
		keyList   []string
		valueList []*ValueWrapper
	)
	for k, v := range relationship.edge.GetProps() {
		keyList = append(keyList, k)
		valueList = append(valueList, &ValueWrapper{v})
	}

	for i := 0; i < len(keyList); i++ {
		kvMap[keyList[i]] = valueList[i]
	}
	return kvMap
}

// Keys returns a list of keys
func (relationship Relationship) Keys() []string {
	var keys []string
	for key := range relationship.edge.GetProps() {
		keys = append(keys, key)
	}
	return keys
}

// Values returns a list of values wrapped as ValueWrappers
func (relationship Relationship) Values() []*ValueWrapper {
	var values []*ValueWrapper
	for _, value := range relationship.edge.GetProps() {
		values = append(values, &ValueWrapper{value})
	}
	return values
}

// String returns a string representing relationship
// Relationship format: [:edge src->dst @ranking {props}]
func (relationship Relationship) String() string {
	edge := relationship.edge
	var keyList []string
	var kvStr []string
	var src string
	var dst string
	for k := range edge.GetProps() {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		kvTemp := fmt.Sprintf("%s: %s", k, ValueWrapper{edge.GetProps()[k]}.String())
		kvStr = append(kvStr, kvTemp)
	}
	if relationship.edge.GetType() > 0 {
		src = ValueWrapper{edge.GetSrc()}.String()
		dst = ValueWrapper{edge.GetDst()}.String()
	} else {
		src = ValueWrapper{edge.GetDst()}.String()
		dst = ValueWrapper{edge.GetSrc()}.String()
	}
	return fmt.Sprintf(`[:%s %s->%s @%d {%s}]`,
		string(edge.GetName()), src, dst, edge.GetRanking(), fmt.Sprintf("%s", strings.Join(kvStr, ", ")))
}

func (r1 Relationship) IsEqualTo(r2 *Relationship) bool {
	if r1.edge.GetSrc().IsSetSVal() && r2.edge.GetSrc().IsSetSVal() &&
		r1.edge.GetDst().IsSetSVal() && r2.edge.GetDst().IsSetSVal() {
		s1, _ := ValueWrapper{r1.edge.GetSrc()}.AsString()
		s2, _ := ValueWrapper{r2.edge.GetSrc()}.AsString()
		return s1 == s2 && string(r1.edge.GetName()) == string(r2.edge.GetName()) && r1.edge.GetRanking() == r2.edge.GetRanking()
	} else if r1.edge.GetSrc().IsSetIVal() && r2.edge.GetSrc().IsSetIVal() &&
		r1.edge.GetDst().IsSetIVal() && r2.edge.GetDst().IsSetIVal() {
		s1, _ := ValueWrapper{r1.edge.GetSrc()}.AsInt()
		s2, _ := ValueWrapper{r2.edge.GetSrc()}.AsInt()
		return s1 == s2 && string(r1.edge.GetName()) == string(r2.edge.GetName()) && r1.edge.GetRanking() == r2.edge.GetRanking()
	}
	return false
}

func GenPathWrapper(path types.Path) (*PathWrapper, error) {
	if path == nil {
		return nil, fmt.Errorf("failed to generate Path Wrapper: invalid path")
	}
	var (
		nodeList         []*Node
		relationshipList []*Relationship
		segList          []segment
		edge             types.Edge
		segStartNode     *Node
		segEndNode       *Node
		segType          types.EdgeType
	)
	src, err := GenNode(path.GetSrc())
	if err != nil {
		return nil, err
	}
	nodeList = append(nodeList, src)

	for _, step := range path.GetSteps() {
		dst, err := GenNode(step.GetDst())
		if err != nil {
			return nil, err
		}
		nodeList = append(nodeList, dst)
		// determine direction
		stepType := step.GetType()
		if stepType > 0 {
			segStartNode = src
			segEndNode = dst
			segType = stepType
		} else {
			segStartNode = dst // switch with src
			segEndNode = src
			segType = -stepType
		}
		edge = NewEdge(
			segStartNode.getRawID(),
			segEndNode.getRawID(),
			segType,
			step.GetName(),
			step.GetRanking(),
			step.GetProps(),
		)
		relationship, err := GenRelationship(edge)
		if err != nil {
			return nil, err
		}
		relationshipList = append(relationshipList, relationship)

		// Check segments
		if len(segList) > 0 {
			prevStart := segList[len(segList)-1].startNode.GetID()
			prevEnd := segList[len(segList)-1].endNode.GetID()
			nextStart := segStartNode.GetID()
			nextEnd := segEndNode.GetID()
			if prevStart.String() != nextStart.String() && prevStart.String() != nextEnd.String() &&
				prevEnd.String() != nextStart.String() && prevEnd.String() != nextEnd.String() {
				return nil, fmt.Errorf("failed to generate PathWrapper, Path received is invalid")
			}
		}
		segList = append(segList, segment{
			startNode:    segStartNode,
			relationship: relationship,
			endNode:      segEndNode,
		})
		src = dst
	}
	return &PathWrapper{
		path:             path,
		nodeList:         nodeList,
		relationshipList: relationshipList,
		segments:         segList,
	}, nil
}

func (path *PathWrapper) GetPathLength() int {
	return len(path.segments)
}

func (path *PathWrapper) GetNodes() []*Node {
	return path.nodeList
}

func (path *PathWrapper) GetRelationships() []*Relationship {
	return path.relationshipList
}

func (path *PathWrapper) GetSegments() []segment {
	return path.segments
}

func (path *PathWrapper) ContainsNode(node Node) bool {
	for _, n := range path.nodeList {
		if n.IsEqualTo(&node) {
			return true
		}
	}
	return false
}

func (path *PathWrapper) ContainsRelationship(relationship *Relationship) bool {
	for _, r := range path.relationshipList {
		if r.IsEqualTo(relationship) {
			return true
		}
	}
	return false
}

func (path *PathWrapper) GetStartNode() (*Node, error) {
	if len(path.segments) == 0 {
		return nil, fmt.Errorf("failed to get start node, no node in the path")
	}
	return path.segments[0].startNode, nil
}

func (path *PathWrapper) GetEndNode() (*Node, error) {
	if len(path.segments) == 0 {
		return nil, fmt.Errorf("failed to get end node, no node in the path")
	}
	return path.segments[len(path.segments)-1].endNode, nil
}

// Path format: <("VertexID" :tag1{k0: v0,k1: v1})
// -[:TypeName@ranking {edgeProps}]->
// ("VertexID2" :tag1{k0: v0,k1: v1} :tag2{k2: v2})
// -[:TypeName@ranking {edgeProps}]->
// ("VertexID3" :tag1{k0: v0,k1: v1})>
func (pathWrap *PathWrapper) String() string {
	path := pathWrap.path
	src := path.GetSrc()
	steps := path.GetSteps()
	resStr := ValueWrapper{NewValue(src)}.String()
	for _, step := range steps {
		var keyList []string
		var kvStr []string
		for k := range step.GetProps() {
			keyList = append(keyList, k)
		}
		sort.Strings(keyList)
		for _, k := range keyList {
			kvTemp := fmt.Sprintf("%s: %s", k, ValueWrapper{step.GetProps()[k]}.String())
			kvStr = append(kvStr, kvTemp)
		}
		var dirChar1 string
		var dirChar2 string
		if step.GetType() > 0 {
			dirChar1 = "-"
			dirChar2 = "->"
		} else {
			dirChar1 = "<-"
			dirChar2 = "-"
		}
		resStr = resStr + fmt.Sprintf("%s[:%s@%d {%s}]%s%s",
			dirChar1,
			string(step.GetName()),
			step.GetRanking(),
			fmt.Sprintf("%s", strings.Join(kvStr, ", ")),
			dirChar2,
			ValueWrapper{NewValue(step.GetDst())}.String())
	}
	return "<" + resStr + ">"
}

func (p1 *PathWrapper) IsEqualTo(p2 *PathWrapper) bool {
	// Check length
	if len(p1.nodeList) != len(p2.nodeList) || len(p1.relationshipList) != len(p2.relationshipList) ||
		len(p1.segments) != len(p2.segments) {
		return false
	}
	// Check nodes
	for i := 0; i < len(p1.nodeList); i++ {
		if !p1.nodeList[i].IsEqualTo(p2.nodeList[i]) {
			return false
		}
	}
	// Check relationships
	for i := 0; i < len(p1.relationshipList); i++ {
		if !p1.relationshipList[i].IsEqualTo(p2.relationshipList[i]) {
			return false
		}
	}
	// Check segments
	for i := 0; i < len(p1.segments); i++ {
		if !p1.segments[i].startNode.IsEqualTo(p2.segments[i].startNode) ||
			!p1.segments[i].endNode.IsEqualTo(p2.segments[i].endNode) ||
			!p1.segments[i].relationship.IsEqualTo(p2.segments[i].relationship) {
			return false
		}
	}
	return true
}

func GenTimeWrapper(time types.Time) (*TimeWrapper, error) {
	if time == nil {
		return nil, fmt.Errorf("failed to generate Time: invalid Time")
	}

	return &TimeWrapper{
		time: time,
	}, nil
}

// getHour returns the hour in UTC
func (t TimeWrapper) getHour() int8 {
	return t.time.GetHour()
}

// getHour returns the minute in UTC
func (t TimeWrapper) getMinute() int8 {
	return t.time.GetMinute()
}

// getHour returns the second in UTC
func (t TimeWrapper) getSecond() int8 {
	return t.time.GetSec()
}

func (t TimeWrapper) getMicrosec() int32 {
	return t.time.GetMicrosec()
}

// getRawTime returns a types.Time object in UTC.
func (t TimeWrapper) getRawTime() types.Time {
	return t.getRawTime()
}

func (t1 TimeWrapper) IsEqualTo(t2 TimeWrapper) bool {
	return t1.getHour() == t2.getHour() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getSecond() == t2.getSecond() &&
		t1.getMicrosec() == t2.getMicrosec()
}

func GenDateWrapper(date types.Date) (*DateWrapper, error) {
	if date == nil {
		return nil, fmt.Errorf("failed to generate date: invalid date")
	}
	return &DateWrapper{
		date: date,
	}, nil
}

func (d DateWrapper) getYear() int16 {
	return d.date.GetYear()
}

func (d DateWrapper) getMonth() int8 {
	return d.date.GetMonth()
}

func (d DateWrapper) getDay() int8 {
	return d.date.GetDay()
}

// getRawDate returns a types.Date object in UTC.
func (d DateWrapper) getRawDate() types.Date {
	return d.getRawDate()
}

func (d1 DateWrapper) IsEqualTo(d2 DateWrapper) bool {
	return d1.getYear() == d2.getYear() &&
		d1.getMonth() == d2.getMonth() &&
		d1.getDay() == d2.getDay()
}

func GenDateTimeWrapper(datetime types.DateTime) (*DateTimeWrapper, error) {
	if datetime == nil {
		return nil, fmt.Errorf("failed to generate datetime: invalid datetime")
	}
	return &DateTimeWrapper{
		dateTime: datetime,
	}, nil
}

func (dt DateTimeWrapper) getYear() int16 {
	return dt.dateTime.GetYear()
}

func (dt DateTimeWrapper) getMonth() int8 {
	return dt.dateTime.GetMonth()
}

func (dt DateTimeWrapper) getDay() int8 {
	return dt.dateTime.GetDay()
}

func (dt DateTimeWrapper) getHour() int8 {
	return dt.dateTime.GetHour()
}

func (dt DateTimeWrapper) getMinute() int8 {
	return dt.dateTime.GetMinute()
}

func (dt DateTimeWrapper) getSecond() int8 {
	return dt.dateTime.GetSec()
}

func (dt DateTimeWrapper) getMicrosec() int32 {
	return dt.dateTime.GetMicrosec()
}

func (dt1 DateTimeWrapper) IsEqualTo(dt2 DateTimeWrapper) bool {
	return dt1.getYear() == dt2.getYear() &&
		dt1.getMonth() == dt2.getMonth() &&
		dt1.getDay() == dt2.getDay() &&
		dt1.getHour() == dt2.getHour() &&
		dt1.getSecond() == dt2.getSecond() &&
		dt1.getSecond() == dt2.getSecond() &&
		dt1.getMicrosec() == dt2.getMicrosec()
}

// getRawDateTime returns a types.DateTime object representing local dateTime in UTC.
func (dt DateTimeWrapper) getRawDateTime() types.DateTime {
	return dt.dateTime
}

func GenValWraps(row types.Row) ([]*ValueWrapper, error) {
	if row == nil {
		return nil, fmt.Errorf("failed to generate valueWrapper: invalid row")
	}
	var valWraps []*ValueWrapper
	for _, val := range row.GetValues() {
		if val == nil {
			return nil, fmt.Errorf("failed to generate valueWrapper: value is nil")
		}
		valWraps = append(valWraps, &ValueWrapper{val})
	}
	return valWraps, nil
}

func checkIndex(index int, list interface{}) error {
	if _, ok := list.([]types.Row); ok {
		if index < 0 || index >= len(list.([]types.Row)) {
			return fmt.Errorf("failed to get Value, the index is out of range")
		}
		return nil
	} else if _, ok := list.([]*ValueWrapper); ok {
		if index < 0 || index >= len(list.([]*ValueWrapper)) {
			return fmt.Errorf("failed to get Value, the index is out of range")
		}
		return nil
	}
	return fmt.Errorf("given list type is invalid")
}

func graphvizString(s string) string {
	s = strings.Replace(s, "{", "\\{", -1)
	s = strings.Replace(s, "}", "\\}", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	return s
}

func prettyFormatJsonString(value []byte) string {
	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, value, "", "  "); err != nil {
		return string(value)
	}
	return prettyJson.String()
}

func name(planNodeDesc types.PlanNodeDescription) string {
	return fmt.Sprintf("%s_%d", planNodeDesc.GetName(), planNodeDesc.GetId())
}

func condEdgeLabel(condNode types.PlanNodeDescription, doBranch bool) string {
	name := strings.ToLower(string(condNode.GetName()))
	if strings.HasPrefix(name, "select") {
		if doBranch {
			return "Y"
		}
		return "N"
	}
	if strings.HasPrefix(name, "loop") {
		if doBranch {
			return "Do"
		}
	}
	return ""
}

func nodeString(planNodeDesc types.PlanNodeDescription, planNodeName string) string {
	var outputVar = graphvizString(string(planNodeDesc.GetOutputVar()))
	var inputVar string
	if planNodeDesc.IsSetDescription() {
		desc := planNodeDesc.GetDescription()
		for _, pair := range desc {
			key := string(pair.GetKey())
			if key == "inputVar" {
				inputVar = graphvizString(string(pair.GetValue()))
			}
		}
	}
	return fmt.Sprintf("\t\"%s\"[label=\"{%s|outputVar: %s|inputVar: %s}\", shape=Mrecord];\n",
		planNodeName, planNodeName, outputVar, inputVar)
}

func edgeString(start, end string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\";\n", start, end)
}

func conditionalEdgeString(start, end, label string) string {
	return fmt.Sprintf("\t\"%s\"->\"%s\"[label=\"%s\", style=dashed];\n", start, end, label)
}

func conditionalNodeString(name string) string {
	return fmt.Sprintf("\t\"%s\"[shape=diamond];\n", name)
}

func nodeById(p types.PlanDescription, nodeId int64) types.PlanNodeDescription {
	line := p.GetNodeIndexMap()[nodeId]
	return p.GetPlanNodeDescs()[line]
}

func findBranchEndNode(p types.PlanDescription, condNodeId int64, isDoBranch bool) int64 {
	for _, node := range p.GetPlanNodeDescs() {
		if node.IsSetBranchInfo() {
			bInfo := node.GetBranchInfo()
			if bInfo.GetConditionNodeID() == condNodeId && bInfo.GetIsDoBranch() == isDoBranch {
				return node.GetId()
			}
		}
	}
	return -1
}

func findFirstStartNodeFrom(p types.PlanDescription, nodeId int64) int64 {
	node := nodeById(p, nodeId)
	for {
		deps := node.GetDependencies()
		if len(deps) == 0 {
			if strings.ToLower(string(node.GetName())) != "start" {
				return -1
			}
			return node.GetId()
		}
		node = nodeById(p, deps[0])
	}
}
