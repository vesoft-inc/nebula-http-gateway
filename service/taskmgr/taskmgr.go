package taskmgr

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
	Runner *cmd.Runner
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
	tid = fmt.Sprintf("%d", taskID)
	taskID++
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
		for _, r := range task.Runner.Readers {
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
	Query TaskAction = iota
	QueryAll
	Stop
	StopAll
	UnknownAction
)

func NewTaskAction(action string) TaskAction {
	switch action {
	case "stop":
		return Stop
	case "stopAll":
		return StopAll
	case "query":
		return Query
	case "queryAll":
		return QueryAll
	default:
		return UnknownAction
	}
}

func (action TaskAction) String() string {
	switch action {
	case Stop:
		return "stop"
	case StopAll:
		return "stopAll"
	case Query:
		return "query"
	case QueryAll:
		return "queryAll"
	default:
		return "unknownAction"
	}
}
