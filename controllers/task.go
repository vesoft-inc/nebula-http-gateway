package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
		taskID string = importer.NewTaskID()
		err    error
	)

	task := importer.NewTask(taskID)
	importer.GetTaskMgr().PutTask(taskID, &task)

	err = json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	if err != nil {
		err = importerErrors.Wrap(importerErrors.InvalidConfigPathOrFormat, err)
	} else {
		err = importer.Import(taskID, params.ConfigPath, &params.ConfigBody)
	}

	if err != nil {
		// task err: import task not start err handle
		task.TaskStatus = importer.StatusAborted.String()
		logs.Error(fmt.Sprintf("Failed to start a import task: `%s`, task result: `%v`", taskID, err))

		res.Code = -1
		res.Message = err.Error()
	} else {
		res.Code = 0
		res.Data = []string{taskID}
		res.Message = fmt.Sprintf("Import task %s submit successfully", taskID)
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
