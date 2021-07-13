package controllers

import (
	"encoding/json"
	"fmt"
	dao "nebula-http-gateway/service/dao"
	"nebula-http-gateway/service/taskmgr"

	"github.com/astaxie/beego"
)

type TaskController struct {
	beego.Controller
}

type ImportRequest struct {
	ConfigPath string `json:"configPath"`
}

type ActionRequest struct {
	TaskID     string `json:"taskID"`
	TaskAction string `json:"taskAction"`
}

func (this *TaskController) Import() {
	var res Response
	var params ImportRequest

	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	tid, err := dao.Import(params.ConfigPath)
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

func (this *TaskController) Action() {
	var res Response
	var params ActionRequest

	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	result, err := dao.Action(params.TaskID, taskmgr.NewTaskAction(params.TaskAction))
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
