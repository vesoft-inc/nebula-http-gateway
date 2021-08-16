package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-http-gateway/service/importer"
	"github.com/vesoft-inc/nebula-importer/pkg/config"

	importerErrors "github.com/vesoft-inc/nebula-importer/pkg/errors"
)

type TaskController struct {
	beego.Controller
}

type ImportRequest struct {
	ConfigPath string            `json:"configPath"`
	ConfigBody config.YAMLConfig `json:"configBody"`
}

type ImportActionRequest struct {
	TaskID     string `json:"taskID"`
	TaskAction string `json:"taskAction"`
}

func (this *TaskController) Import() {
	var (
		res    Response
		params ImportRequest
		tid    string
		err    error
	)

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	if err == nil {
		tid, err = importer.Import(params.ConfigPath, &params.ConfigBody)
	} else {
		err = importerErrors.Wrap(importerErrors.InvalidConfigPathOrFormat, err)
	}

	if err == nil {
		res.Code = 0
		res.Message = fmt.Sprintf("Import task %s submit successfully", tid)
	} else {
		res.Code = -1
		res.Message = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *TaskController) ImportAction() {
	var res Response
	var params ImportActionRequest

	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	result, err := importer.ImportAction(params.TaskID, importer.NewTaskAction(params.TaskAction))
	if err == nil {
		res.Code = 0
		res.Data = result
		res.Message = "Processing a task action successfully"
	} else {
		res.Code = -1
		res.Message = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
