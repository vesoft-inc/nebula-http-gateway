package importer

import (
	"fmt"
	"sync"

	"github.com/vesoft-inc/nebula-importer/pkg/cmd"
)

var (
	taskmgr *TaskMgr = &TaskMgr{
		tasks: sync.Map{},
	}

	taskID uint64 = 0
	mux    sync.Mutex
)

type TaskMgr struct {
	tasks sync.Map
}

type Task struct {
	runner *cmd.Runner

	TaskID     string `json:"taskID"`
	TaskStatus string `json:"taskStatus"`
}

func NewTask(taskID string) Task {
	return Task{
		runner:     &cmd.Runner{},
		TaskID:     taskID,
		TaskStatus: StatusProcessing.String(),
	}
}

func (task *Task) GetRunner() *cmd.Runner {
	return task.runner
}

func GetTaskMgr() *TaskMgr {
	return taskmgr
}

func GetTaskID() (tid uint64) {
	mux.Lock()
	defer mux.Unlock()
	tid = taskID
	return tid
}

func NewTaskID() (tid string) {
	mux.Lock()
	defer mux.Unlock()
	taskID++
	tid = fmt.Sprintf("%d", taskID)
	return tid
}

func (mgr *TaskMgr) GetTask(tid string) (*Task, bool) {
	if task, ok := mgr.tasks.Load(tid); ok {
		return task.(*Task), true
	}

	return nil, false
}

func (mgr *TaskMgr) PutTask(tid string, task *Task) {
	mgr.tasks.Store(tid, task)
}

func (mgr *TaskMgr) DelTask(tid string) {
	mgr.tasks.Delete(tid)
}

func (mgr *TaskMgr) StopTask(tid string) bool {
	if task, ok := mgr.GetTask(tid); ok {
		for _, r := range task.GetRunner().Readers {
			r.Stop()
		}

		mgr.DelTask(tid)

		return true
	}

	return false
}

func (mgr *TaskMgr) GetAllTaskIDs() []string {
	ids := make([]string, 0)
	mgr.tasks.Range(func(k, v interface{}) bool {
		ids = append(ids, k.(string))
		return true
	})

	return ids
}

type TaskAction int

const (
	ActionUnknown TaskAction = iota
	ActionQuery
	ActionQueryAll
	ActionStop
	ActionStopAll
)

var taskActionMap = map[TaskAction]string{
	ActionQuery:    "actionQuery",
	ActionQueryAll: "actionQueryAll",
	ActionStop:     "actionStop",
	ActionStopAll:  "actionStopAll",
}

var taskActionRevMap = map[string]TaskAction{
	"actionQuery":    ActionQuery,
	"actionQueryAll": ActionQueryAll,
	"actionStop":     ActionStop,
	"actionStopAll":  ActionStopAll,
}

func NewTaskAction(action string) TaskAction {
	if v, ok := taskActionRevMap[action]; ok {
		return v
	}
	return ActionUnknown
}

func (action TaskAction) String() string {
	if v, ok := taskActionMap[action]; ok {
		return v
	}
	return "actionUnknown"
}

type TaskStatus int

const (
	StatusUnknown TaskStatus = iota
	StatusStoped
	StatusProcessing
	StatusNotExisted
	StatusAborted
)

var taskStatusMap = map[TaskStatus]string{
	StatusStoped:     "statusStoped",
	StatusProcessing: "statusProcessing",
	StatusNotExisted: "statusNotExisted",
	StatusAborted:    "statusAborted",
}

var taskStatusRevMap = map[string]TaskStatus{
	"statusStoped":     StatusStoped,
	"statusProcessing": StatusProcessing,
	"statusNotExisted": StatusNotExisted,
	"statusAborted":    StatusAborted,
}

func NewTaskStatus(status string) TaskStatus {
	if v, ok := taskStatusRevMap[status]; ok {
		return v
	}
	return StatusUnknown
}

func (status TaskStatus) String() string {
	if v, ok := taskStatusMap[status]; ok {
		return v
	}
	return "statusUnknown"
}
