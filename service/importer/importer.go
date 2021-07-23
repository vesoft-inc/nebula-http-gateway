package importer

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/vesoft-inc/nebula-importer/pkg/cmd"
	"github.com/vesoft-inc/nebula-importer/pkg/config"
	importerErrors "github.com/vesoft-inc/nebula-importer/pkg/errors"
)

type ImportResult struct {
	TaskId      string `json:"taskId"`
	TimeCost    string `json:"timeCost"` // Milliseconds
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

func Import(configPath string) (tid string, err error) {

	path := filepath.Join(
		beego.AppConfig.String("uploadspath"),
		configPath,
	)

	conf, err := config.Parse(path)
	if err != nil {
		beego.Error(err.(importerErrors.ImporterError))
		return tid, err.(importerErrors.ImporterError)
	}

	if err := conf.ValidateAndReset(""); err != nil {
		return tid, err
	}

	tid = NewTaskID()

	beego.Debug(fmt.Sprintf("Start a import task: `%s` with config: `%s`", tid, path))

	task := &Task{
		Runner: &cmd.Runner{},
	}
	GetTaskMgr().PutTask(tid, task)

	go func() {
		result := ImportResult{}
		now := time.Now()
		task.Runner.Run(conf)
		task.TimeCost = time.Since(now).Milliseconds()

		result.TaskId = tid
		result.TimeCost = fmt.Sprintf("%dms", task.TimeCost)

		if rerr := task.Runner.Error(); rerr != nil {
			err, _ := rerr.(importerErrors.ImporterError)

			result.ErrorResult.ErrorCode = err.ErrCode
			result.ErrorResult.ErrorMsg = err.ErrMsg.Error()

			beego.Error(fmt.Sprintf("Failed to finish a import task: `%s` with config: `%s`, task result: `%v`", tid, path, result))
		} else {
			result.FailedRows = task.Runner.NumFailed
			GetTaskMgr().DelTask(tid)

			beego.Debug(fmt.Sprintf("Success to finish a import task: `%s` with config: `%s`, task result: `%v`", tid, path, result))
		}
	}()
	return tid, nil
}

func ImportAction(taskID string, taskAction TaskAction) (result ActionResult, err error) {
	beego.Debug(fmt.Sprintf("Start a import task action: `%s` for task: `%s`", taskAction.String(), taskID))

	result = ActionResult{}

	switch taskAction {
	case Stop:
		if ok := GetTaskMgr().StopTask(taskID); !ok {
			tid, _ := strconv.ParseUint(taskID, 0, 64)
			if tid > GetTaskID() {
				result.TaskStatus = "Task not exist"
			} else {
				result.TaskStatus = "Task has stoped"
			}
		} else {
			result.TaskStatus = "Task stop successfully"
		}
	case StopAll:
		tids := GetTaskMgr().GetAllTaskIDs()
		result.TaskIDs = make([]string, 0, len(tids))

		for _, tid := range tids {
			GetTaskMgr().StopTask(tid)
			result.TaskIDs = append(result.TaskIDs, tid)
		}

		result.TaskStatus = "Tasks stop successfully"
	case Query:
		if _, ok := GetTaskMgr().GetTask(taskID); !ok {
			tid, _ := strconv.ParseUint(taskID, 0, 64)
			if tid > GetTaskID() {
				result.TaskStatus = "Task not exist"
			} else {
				result.TaskStatus = "Task has stoped"
			}
		} else {
			result.TaskStatus = "Task is processing"
		}
	case QueryAll:
		tids := GetTaskMgr().GetAllTaskIDs()
		result.TaskIDs = make([]string, 0, len(tids))

		for _, tid := range tids {
			GetTaskMgr().StopTask(tid)
			result.TaskIDs = append(result.TaskIDs, tid)
		}

		result.TaskStatus = "Tasks are processing"
	default:
		err = errors.New("unknown task action")
	}

	beego.Debug(fmt.Sprintf("The import task action: `%s` for task: `%s` finished, action result: `%v`", taskAction.String(), taskID, result))

	return result, err
}
