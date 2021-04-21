package dao

import (
	"errors"
	"log"

	"nebula-http-gateway/service/pool"
	common "nebula-http-gateway/utils"

	nebula "github.com/vesoft-inc/nebula-go"
	nebulaType "github.com/vesoft-inc/nebula-go/nebula"
)

type ExecuteResult struct {
	Headers  []string                `json:"headers"`
	Tables   []map[string]common.Any `json:"tables"`
	TimeCost int32                   `json:"timeCost"`
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

func getListInfo(valWarp *nebula.ValueWrapper, listType string, _verticesParsedList []common.Any, _edgesParsedList []common.Any, _pathsParsedList []common.Any) ([]common.Any, []common.Any, []common.Any, error) {
	var valueList []nebula.ValueWrapper
	var err error
	if listType == "list" {
		valueList, err = valWarp.AsList()
	} else if listType == "set" {
		valueList, err = valWarp.AsDedupList()
	}
	if err != nil {
		return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
	}
	for _, v := range valueList {
		var props = make(map[string]common.Any)
		vType := v.GetType()
		props["type"] = vType
		if vType == "vertex" {
			props, err = getVertexInfo(&v, props)
			if err == nil {
				_verticesParsedList = append(_verticesParsedList, props)
			} else {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "edge" {
			props, err = getEdgeInfo(&v, props)
			if err == nil {
				_edgesParsedList = append(_edgesParsedList, props)
			} else {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "path" {
			props, err = getPathInfo(&v, props)
			if err == nil {
				_pathsParsedList = append(_pathsParsedList, props)
			} else {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "list" {
			_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getListInfo(&v, "list", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "map" {
			_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getMapInfo(&v, _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "set" {
			_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getListInfo(&v, "set", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else {
			// no need to parse basic value now
		}
	}
	return _verticesParsedList, _edgesParsedList, _pathsParsedList, nil
}

func getMapInfo(valWarp *nebula.ValueWrapper, _verticesParsedList []common.Any, _edgesParsedList []common.Any, _pathsParsedList []common.Any) ([]common.Any, []common.Any, []common.Any, error) {
	valueMap, err := valWarp.AsMap()
	if err != nil {
		return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
	}
	for _, v := range valueMap {
		vType := v.GetType()
		if vType == "vertex" {
			var _props map[string]common.Any
			_props, err = getVertexInfo(&v, _props)
			if err == nil {
				_verticesParsedList = append(_verticesParsedList, _props)
			} else {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "edge" {
			var _props map[string]common.Any
			_props, err = getEdgeInfo(&v, _props)
			if err == nil {
				_edgesParsedList = append(_edgesParsedList, _props)
			} else {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "path" {
			var _props map[string]common.Any
			_props, err = getPathInfo(&v, _props)
			if err == nil {
				_pathsParsedList = append(_pathsParsedList, _props)
			} else {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "list" {
			_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getListInfo(&v, "list", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "map" {
			_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getMapInfo(&v, _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else if vType == "set" {
			_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getListInfo(&v, "set", _verticesParsedList, _edgesParsedList, _pathsParsedList)
			if err != nil {
				return _verticesParsedList, _edgesParsedList, _pathsParsedList, err
			}
		} else {
			// no need to parse basic value now
		}
	}
	return _verticesParsedList, _edgesParsedList, _pathsParsedList, nil
}

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
		var rowValue = make(map[string]common.Any)
		if format == "row" {
			result.Headers = []string{"id", "name", "dependencies", "profiling data", "operator info"}
			rows := resp.MakePlanByRow()
			for i := 0; i < len(rows); i++ {
				rowValue["id"] = rows[i][0]
				rowValue["name"] = rows[i][1]
				rowValue["dependencies"] = rows[i][2]
				rowValue["profiling data"] = rows[i][3]
				rowValue["operator info"] = rows[i][4]
				result.Tables = append(result.Tables, rowValue)
			}
			return result, err
		} else if format == "dot" {
			result.Headers = append(result.Headers, "format")
			rowValue["format"] = resp.MakeDotGraph()
			result.Tables = append(result.Tables, rowValue)
			return result, err
		} else if format == "dot:struct" {
			result.Headers = append(result.Headers, "format")
			rowValue["format"] = resp.MakeDotGraphByStruct()
			result.Tables = append(result.Tables, rowValue)
			return result, err
		}
	}

	if !resp.IsSucceed() {
		log.Printf("ErrorCode: %v, ErrorMsg: %s", resp.GetErrorCode(), resp.GetErrorMsg())
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
			var _verticesParsedList = make([]common.Any, 0)
			var _edgesParsedList = make([]common.Any, 0)
			var _pathsParsedList = make([]common.Any, 0)
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
					_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getListInfo(rowData, "list", _verticesParsedList, _edgesParsedList, _pathsParsedList)
				} else if valueType == "set" {
					_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getListInfo(rowData, "set", _verticesParsedList, _edgesParsedList, _pathsParsedList)
				} else if valueType == "map" {
					_verticesParsedList, _edgesParsedList, _pathsParsedList, err = getMapInfo(rowData, _verticesParsedList, _edgesParsedList, _pathsParsedList)
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
