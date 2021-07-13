package dao

import (
	"encoding/json"
	"errors"
	"strconv"

	"nebula-http-gateway/service/logger"
	"nebula-http-gateway/service/pool"
	taskmgr "nebula-http-gateway/service/taskmgr"
	common "nebula-http-gateway/utils"

	nebula "github.com/vesoft-inc/nebula-go/v2"
	nebulaType "github.com/vesoft-inc/nebula-go/v2/nebula"
	"github.com/vesoft-inc/nebula-importer/pkg/cmd"
	"github.com/vesoft-inc/nebula-importer/pkg/config"
	importerErrors "github.com/vesoft-inc/nebula-importer/pkg/errors"
)

type ExecuteResult struct {
	Headers  []string                `json:"headers"`
	Tables   []map[string]common.Any `json:"tables"`
	TimeCost int32                   `json:"timeCost"`
}
type ImportResult struct {
	TaskId      string `json:"taskId"`
	FailedRows  int64  `json:"failedRows"`
	ErrorResult struct {
		ErrorCode int    `json:"errorCode"`
		ErrorMsg  string `json:"errorMsg"`
	}
}

type ActionResult struct {
	TaskIDs    []string `json:"taskIDs"`
	TaskStatus string   `json:"taskStatus"`
}

type list []common.Any

// Connect return if the nebula connect succeed
func Connect(address string, port int, username string, password string) (nsid string, err error) {
	nsid, err = pool.NewConnection(address, port, username, password)
	if err != nil {
		return "", err
	}
	return nsid, err
}

func Disconnect(nsid string) {
	pool.Disconnect(nsid)
}

func Execute(nsid string, gql string) (result ExecuteResult, err error) {
	result = ExecuteResult{
		Headers: make([]string, 0),
		Tables:  make([]map[string]common.Any, 0),
	}
	connection, err := pool.GetConnection(nsid)
	if err != nil {
		return result, err
	}

	responseChannel := make(chan pool.ChannelResponse)
	connection.RequestChannel <- pool.ChannelRequest{
		Gql:             gql,
		ResponseChannel: responseChannel,
	}
	response := <-responseChannel
	if response.Error != nil {
		return result, response.Error
	}
	resp := response.Result
	if resp.IsSetPlanDesc() {
		format := string(resp.GetPlanDesc().GetFormat())
		if format == "row" {
			result.Headers = []string{"id", "name", "dependencies", "profiling data", "operator info"}
			rows := resp.MakePlanByRow()
			for i := 0; i < len(rows); i++ {
				var rowValue = make(map[string]common.Any)
				rowValue["id"] = rows[i][0]
				rowValue["name"] = rows[i][1]
				rowValue["dependencies"] = rows[i][2]
				rowValue["profiling data"] = rows[i][3]
				rowValue["operator info"] = rows[i][4]
				result.Tables = append(result.Tables, rowValue)
			}
			return result, err
		} else {
			var rowValue = make(map[string]common.Any)
			result.Headers = append(result.Headers, "format")
			if format == "dot" {
				rowValue["format"] = resp.MakeDotGraph()
			} else if format == "dot:struct" {
				rowValue["format"] = resp.MakeDotGraphByStruct()
			}
			result.Tables = append(result.Tables, rowValue)
			return result, err
		}
	}

	if !resp.IsSucceed() {
		logger.Debugf("ErrorCode: %v, ErrorMsg: %s", resp.GetErrorCode(), resp.GetErrorMsg())
		return result, errors.New(string(resp.GetErrorMsg()))
	}
	if !resp.IsEmpty() {
		rowSize := resp.GetRowSize()
		colSize := resp.GetColSize()
		colNames := resp.GetColNames()
		result.Headers = colNames
		for i := 0; i < rowSize; i++ {
			var rowValue = make(map[string]common.Any)
			record, err := resp.GetRowValuesByIndex(i)
			var _verticesParsedList = make(list, 0)
			var _edgesParsedList = make(list, 0)
			var _pathsParsedList = make(list, 0)
			if err != nil {
				return result, err
			}
			for j := 0; j < colSize; j++ {
				rowData, err := record.GetValueByIndex(j)
				if err != nil {
					return result, err
				}
				value, err := getValue(rowData)
				if err != nil {
					return result, err
				}
				rowValue[result.Headers[j]] = value
				valueType := rowData.GetType()
				if valueType == "vertex" {
					var parseValue = make(map[string]common.Any)
					parseValue, err = getVertexInfo(rowData, parseValue)
					parseValue["type"] = "vertex"
					_verticesParsedList = append(_verticesParsedList, parseValue)
				} else if valueType == "edge" {
					var parseValue = make(map[string]common.Any)
					parseValue, err = getEdgeInfo(rowData, parseValue)
					parseValue["type"] = "edge"
					_edgesParsedList = append(_edgesParsedList, parseValue)
				} else if valueType == "path" {
					var parseValue = make(map[string]common.Any)
					parseValue, err = getPathInfo(rowData, parseValue)
					parseValue["type"] = "path"
					_pathsParsedList = append(_pathsParsedList, parseValue)
				} else if valueType == "list" {
					err = getListInfo(rowData, "list", &_verticesParsedList, &_edgesParsedList, &_pathsParsedList)
				} else if valueType == "set" {
					err = getListInfo(rowData, "set", &_verticesParsedList, &_edgesParsedList, &_pathsParsedList)
				} else if valueType == "map" {
					err = getMapInfo(rowData, &_verticesParsedList, &_edgesParsedList, &_pathsParsedList)
				}
				if len(_verticesParsedList) > 0 {
					rowValue["_verticesParsedList"] = _verticesParsedList
				}
				if len(_edgesParsedList) > 0 {
					rowValue["_edgesParsedList"] = _edgesParsedList
				}
				if len(_pathsParsedList) > 0 {
					rowValue["_pathsParsedList"] = _pathsParsedList
				}
				if err != nil {
					return result, err
				}
			}
			result.Tables = append(result.Tables, rowValue)
		}
	}
	result.TimeCost = resp.GetLatency()
	return result, nil
}

func Import(path string) (tid string, err error) {

	conf, err := config.Parse(path)
	if err != nil {
		return tid, err.(importerErrors.ImporterError)
	}

	if err := conf.ValidateAndReset(""); err != nil {
		return tid, err
	}

	tid = taskmgr.NewTaskID()

	logger.Debugf("Start a import task: `%s` with config: `%s`", tid, path)

	task := &taskmgr.Task{
		Runner: &cmd.Runner{},
	}
	taskmgr.GetTaskMgr().PutTask(tid, task)

	go func() {
		result := ImportResult{}
		task.Runner.Run(conf)
		if rerr := task.Runner.Error(); rerr != nil {
			err, _ := rerr.(importerErrors.ImporterError)

			result.ErrorResult.ErrorCode = err.ErrCode
			result.ErrorResult.ErrorMsg = err.ErrMsg.Error()
			result.TaskId = tid

			logger.Errorf("Failed to finish a import task: `%s` with config: `%s`", tid, path)
			resultAsBytes, _ := json.Marshal(result)
			logger.Infof("Import task result: `%s`", string(resultAsBytes))
		} else {
			result.FailedRows = task.Runner.NumFailed
			taskmgr.GetTaskMgr().DelTask(tid)
			logger.Debugf("Success to finish a import task: `%s` with config: `%s`", tid, path)
			resultAsBytes, _ := json.Marshal(result)
			logger.Debugf("Import task result: `%s`", string(resultAsBytes))
		}
	}()
	return tid, nil
}

func Action(taskID string, taskAction taskmgr.TaskAction) (result ActionResult, err error) {
	logger.Debugf("Start a import task action: `%s` for task: `%s`", taskAction.String(), taskID)

	result = ActionResult{}

	switch taskAction {
	case taskmgr.Stop:
		if ok := taskmgr.GetTaskMgr().StopTask(taskID); !ok {
			tid, _ := strconv.ParseUint(taskID, 0, 64)
			if tid > taskmgr.GetTaskID() {
				result.TaskStatus = "Task not exist"
			} else {
				result.TaskStatus = "Task has stoped"
			}
		} else {
			result.TaskStatus = "Task stop successfully"
		}
	case taskmgr.StopAll:
		tids := taskmgr.GetTaskMgr().GetAllTaskIDs()
		result.TaskIDs = make([]string, 0, len(tids))

		for _, tid := range tids {
			taskmgr.GetTaskMgr().StopTask(tid)
			result.TaskIDs = append(result.TaskIDs, tid)
		}

		result.TaskStatus = "Task stop successfully"
	case taskmgr.Query:
		if _, ok := taskmgr.GetTaskMgr().GetTask(taskID); !ok {
			tid, _ := strconv.ParseUint(taskID, 0, 64)
			if tid > taskmgr.GetTaskID() {
				result.TaskStatus = "Task not exist"
			} else {
				result.TaskStatus = "Task has stoped"
			}
		} else {
			result.TaskStatus = "Task is processing"
		}
	case taskmgr.QueryAll:
		tids := taskmgr.GetTaskMgr().GetAllTaskIDs()
		result.TaskIDs = make([]string, 0, len(tids))

		for _, tid := range tids {
			taskmgr.GetTaskMgr().StopTask(tid)
			result.TaskIDs = append(result.TaskIDs, tid)
		}

		result.TaskStatus = "Task is processing"
	default:
		err = errors.New("unknown task action")
	}

	logger.Debugf("The import task action: `%s` for task: `%s` finished", taskAction.String(), taskID)
	resultAsBytes, _ := json.Marshal(result)
	logger.Debugf("Import task action result: `%s`", string(resultAsBytes))

	return result, err
}

func getID(idWarp nebula.ValueWrapper) common.Any {
	idType := idWarp.GetType()
	var vid common.Any
	if idType == "string" {
		vid, _ = idWarp.AsString()
	} else if idType == "int" {
		vid, _ = idWarp.AsInt()
	}
	return vid
}

func getValue(valWarp *nebula.ValueWrapper) (common.Any, error) {
	var valType = valWarp.GetType()
	if valType == "vertex" {
		return valWarp.String(), nil
	} else if valType == "edge" {
		return valWarp.String(), nil
	} else if valType == "path" {
		return valWarp.String(), nil
	} else if valType == "list" {
		return valWarp.String(), nil
	} else if valType == "map" {
		return valWarp.String(), nil
	} else if valType == "set" {
		return valWarp.String(), nil
	} else {
		return getBasicValue(valWarp)
	}
}

func getBasicValue(valWarp *nebula.ValueWrapper) (common.Any, error) {
	var valType = valWarp.GetType()
	if valType == "null" {
		value, err := valWarp.AsNull()
		switch value {
		case nebulaType.NullType___NULL__:
			return "NULL", err
		case nebulaType.NullType_NaN:
			return "NaN", err
		case nebulaType.NullType_BAD_DATA:
			return "BAD_DATA", err
		case nebulaType.NullType_BAD_TYPE:
			return "BAD_TYPE", err
		}
		return "NULL", err
	} else if valType == "bool" {
		return valWarp.AsBool()
	} else if valType == "int" {
		return valWarp.AsInt()
	} else if valType == "float" {
		return valWarp.AsFloat()
	} else if valType == "string" {
		return valWarp.AsString()
	} else if valType == "date" {
		return valWarp.String(), nil
	} else if valType == "time" {
		return valWarp.String(), nil
	} else if valType == "datetime" {
		return valWarp.String(), nil
	} else if valType == "empty" {
		return "_EMPTY_", nil
	}
	return "", nil
}

func getVertexInfo(valWarp *nebula.ValueWrapper, data map[string]common.Any) (map[string]common.Any, error) {
	node, err := valWarp.AsNode()
	if err != nil {
		return nil, err
	}
	id := node.GetID()
	data["vid"] = getID(id)
	tags := make([]string, 0)
	properties := make(map[string]map[string]common.Any)
	for _, tagName := range node.GetTags() {
		tags = append(tags, tagName)
		props, err := node.Properties(tagName)
		if err != nil {
			return nil, err
		}
		_props := make(map[string]common.Any)
		for k, v := range props {
			value, err := getValue(v)
			if err != nil {
				return nil, err
			}
			_props[k] = value
		}
		properties[tagName] = _props
	}
	data["tags"] = tags
	data["properties"] = properties
	return data, nil
}

func getEdgeInfo(valWarp *nebula.ValueWrapper, data map[string]common.Any) (map[string]common.Any, error) {
	relationship, err := valWarp.AsRelationship()
	if err != nil {
		return nil, err
	}
	srcID := relationship.GetSrcVertexID()
	data["srcID"] = getID(srcID)
	dstID := relationship.GetDstVertexID()
	data["dstID"] = getID(dstID)
	edgeName := relationship.GetEdgeName()
	data["edgeName"] = edgeName
	rank := relationship.GetRanking()
	data["rank"] = rank
	properties := make(map[string]common.Any)
	props := relationship.Properties()
	for k, v := range props {
		value, err := getValue(v)
		if err != nil {
			return nil, err
		}
		properties[k] = value
	}
	data["properties"] = properties
	return data, nil
}

func getPathInfo(valWarp *nebula.ValueWrapper, data map[string]common.Any) (map[string]common.Any, error) {
	path, err := valWarp.AsPath()
	if err != nil {
		return nil, err
	}
	relationships := path.GetRelationships()
	var _relationships []common.Any
	for _, relation := range relationships {
		_relation := make(map[string]common.Any)
		srcID := relation.GetSrcVertexID()
		_relation["srcID"] = getID(srcID)
		dstID := relation.GetDstVertexID()
		_relation["dstID"] = getID(dstID)
		edgeName := relation.GetEdgeName()
		_relation["edgeName"] = edgeName
		rank := relation.GetRanking()
		_relation["rank"] = rank
		_relationships = append(_relationships, _relation)
	}
	data["relationships"] = _relationships
	return data, nil
}

func getListInfo(valWarp *nebula.ValueWrapper, listType string, _verticesParsedList *list, _edgesParsedList *list, _pathsParsedList *list) error {
	var valueList []nebula.ValueWrapper
	var err error
	if listType == "list" {
		valueList, err = valWarp.AsList()
	} else if listType == "set" {
		valueList, err = valWarp.AsDedupList()
	}
	if err != nil {
		return err
	}
	for _, v := range valueList {
		var props = make(map[string]common.Any)
		vType := v.GetType()
		props["type"] = vType
		if vType == "vertex" {
			props, err = getVertexInfo(&v, props)
			if err == nil {
				*_verticesParsedList = append(*_verticesParsedList, props)
			} else {
				return err
			}
		} else if vType == "edge" {
			props, err = getEdgeInfo(&v, props)
			if err == nil {
				*_edgesParsedList = append(*_edgesParsedList, props)
			} else {
				return err
			}
		} else if vType == "path" {
			props, err = getPathInfo(&v, props)
			if err == nil {
				*_pathsParsedList = append(*_pathsParsedList, props)
			} else {
				return err
			}
		} else if vType == "list" {
			err = getListInfo(&v, "list", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return err
			}
		} else if vType == "map" {
			err = getMapInfo(&v, _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return err
			}
		} else if vType == "set" {
			err = getListInfo(&v, "set", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return err
			}
		} else {
			// no need to parse basic value now
		}
	}
	return nil
}

func getMapInfo(valWarp *nebula.ValueWrapper, _verticesParsedList *list, _edgesParsedList *list, _pathsParsedList *list) error {
	valueMap, err := valWarp.AsMap()
	if err != nil {
		return err
	}
	for _, v := range valueMap {
		vType := v.GetType()
		if vType == "vertex" {
			var _props map[string]common.Any
			_props, err = getVertexInfo(&v, _props)
			if err == nil {
				*_verticesParsedList = append(*_verticesParsedList, _props)
			} else {
				return err
			}
		} else if vType == "edge" {
			var _props map[string]common.Any
			_props, err = getEdgeInfo(&v, _props)
			if err == nil {
				*_edgesParsedList = append(*_edgesParsedList, _props)
			} else {
				return err
			}
		} else if vType == "path" {
			var _props map[string]common.Any
			_props, err = getPathInfo(&v, _props)
			if err == nil {
				*_pathsParsedList = append(*_pathsParsedList, _props)
			} else {
				return err
			}
		} else if vType == "list" {
			err = getListInfo(&v, "list", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return err
			}
		} else if vType == "map" {
			err = getMapInfo(&v, _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return err
			}
		} else if vType == "set" {
			err = getListInfo(&v, "set", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return err
			}
		} else {
			// no need to parse basic value now
		}
	}
	return nil
}
