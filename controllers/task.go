package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-http-gateway/service/importer"
)

type TaskController struct {
	beego.Controller
}

type ImportRequest struct {
	ConfigPath string `json:"configPath"`
}

type ImportActionRequest struct {
	TaskID     string `json:"taskID"`
	TaskAction string `json:"taskAction"`
}

func (this *TaskController) Import() {
	var res Response
	var params ImportRequest

	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	tid, err := importer.Import(params.ConfigPath)
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
