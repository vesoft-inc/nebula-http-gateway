package importer

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/astaxie/beego"
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
	Results []Task `json:"results"`
	Msg     string `json:"msg"`
}

func Import(taskID string, configPath string, configBody *config.YAMLConfig) (err error) {
	var conf *config.YAMLConfig

	if configPath != "" {
		conf, err = config.Parse(
			filepath.Join(
				beego.AppConfig.String("uploadspath"),
				configPath,
			),
		)

		if err != nil {
			beego.Error(err.(importerErrors.ImporterError))
			return err.(importerErrors.ImporterError)
		}
	} else {
		conf = configBody
	}

	if err := conf.ValidateAndReset(""); err != nil {
		return err
	}

	task := NewTask(taskID)
	GetTaskMgr().PutTask(taskID, &task)

	beego.Debug(fmt.Sprintf("Start a import task: `%s`", taskID))

	go func() {
		result := ImportResult{}

		now := time.Now()
		task.GetRunner().Run(conf)
		timeCost := time.Since(now).Milliseconds()

		result.TaskId = taskID
		result.TimeCost = fmt.Sprintf("%dms", timeCost)

		if rerr := task.GetRunner().Error(); rerr != nil {
			task.TaskStatus = StatusAborted.String()

			err, _ := rerr.(importerErrors.ImporterError)

			result.ErrorResult.ErrorCode = err.ErrCode
			result.ErrorResult.ErrorMsg = err.ErrMsg.Error()

			beego.Error(fmt.Sprintf("Failed to finish a import task: `%s`, task result: `%v`", taskID, result))
		} else {
			task.TaskStatus = StatusStoped.String()

			result.FailedRows = task.GetRunner().NumFailed
			GetTaskMgr().DelTask(taskID)

			beego.Debug(fmt.Sprintf("Success to finish a import task: `%s`, task result: `%v`", taskID, result))
		}
	}()
	return nil
}

func ImportAction(taskID string, taskAction TaskAction) (result ActionResult, err error) {
	beego.Debug(fmt.Sprintf("Start a import task action: `%s` for task: `%s`", taskAction.String(), taskID))

	result = ActionResult{}

	switch taskAction {
	case ActionQuery:
		result.Msg = actionQuery(taskID, &result)
	case ActionQueryAll:
		result.Msg = actionQueryAll(&result)
	case ActionStop:
		result.Msg = actionStop(taskID, &result)
	case ActionStopAll:
		result.Msg = actionStopAll(&result)
	default:
		err = errors.New("unknown task action")
	}

	beego.Debug(fmt.Sprintf("The import task action: `%s` for task: `%s` finished, action result: `%v`", taskAction.String(), taskID, result))

	return result, err
}

func actionQuery(taskID string, result *ActionResult) (msg string) {
	task := Task{}

	tid, _ := strconv.ParseUint(taskID, 0, 64)

	log.Println(tid)

	if tid > GetTaskID() {
		task.TaskID = taskID
		task.TaskStatus = StatusNotExisted.String()
		result.Results = append(result.Results, task)
		return "Task not existed"
	}

	if t, ok := GetTaskMgr().GetTask(taskID); !ok {
		task.TaskID = taskID
		task.TaskStatus = StatusStoped.String()
		result.Results = append(result.Results, task)
		return "Task has stoped"
	} else {
		task.TaskID = t.TaskID
		task.TaskStatus = t.TaskStatus
		result.Results = append(result.Results, task)
		return "Task is processing"
	}
}

func actionQueryAll(result *ActionResult) (msg string) {
	taskIDs := GetTaskMgr().GetAllTaskIDs()
	for _, taskID := range taskIDs {
		actionQuery(taskID, result)
	}

	return "Tasks are processing"
}

func actionStop(taskID string, result *ActionResult) (msg string) {
	GetTaskMgr().StopTask(taskID)
	return actionQuery(taskID, result)
}

func actionStopAll(result *ActionResult) (msg string) {
	taskIDs := GetTaskMgr().GetAllTaskIDs()
	for _, taskID := range taskIDs {
		actionStop(taskID, result)
	}

	return "Tasks stop successfully"
}
