package instance

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

const (
	OpStart   = iota // 0
	OpResume         // 1
	OpRestart        // 2
)

// RunOptions the options when running a FlowAction
type RunOptions struct {
	Op           int
	ReturnID     bool
	FlowURI      string
	InitialState *Instance
	ExecOptions  *ExecOptions
}

// ExecOptions are optional Patch & Interceptor to be used during instance execution
type ExecOptions struct {
	Patch       *support.Patch
	Interceptor *support.Interceptor
}

// IDGenerator generates IDs for flow instances
type IDGenerator interface {
	//NewFlowInstanceID generate a new instance ID
	NewFlowInstanceID() string
}

// ApplyExecOptions applies any execution options to the flow instance
func ApplyExecOptions(instance *Instance, execOptions *ExecOptions) {

	if execOptions != nil {

		if execOptions.Patch != nil {
			logger.Infof("Instance [%s] has patch", instance.ID())
			instance.Patch = execOptions.Patch
			instance.Patch.Init()
		}

		if execOptions.Interceptor != nil {
			logger.Infof("Instance [%s] has interceptor", instance.ID())
			instance.Interceptor = execOptions.Interceptor
			instance.Interceptor.Init()
		}
	}
}

// IDResponse is a response object consists of an ID
type IDResponse struct {
	ID string `json:"id"`
}


////////////////////////////////////////////////////////////////////////////////////////////////////////
// Task Environment

// ExecEnv is a structure that describes the execution environment for a flow
type ExecEnv struct {
	ID        int
	//Task      *definition.Task
	Instance  *Instance
	//ParentEnv *ExecEnv //need?

	flowDef *definition.Definition
	Attrs         map[string]*data.Attribute

	TaskDatas map[string]*TaskData
	LinkDatas map[int]*LinkData


	returnData      map[string]*data.Attribute
	returnError     error

	//taskID string // for deserialization
}

// init initializes the Task Environment, typically called on deserialization
func (te *ExecEnv) init(flowInst *Instance) {

	if te.Instance == nil {

		te.Instance = flowInst
		//te.Task = flowInst.Flow.GetTask(te.taskID)

		for _, td := range te.TaskDatas {
			td.execEnv = te
			td.task = flowInst.Flow.GetTask(td.taskID)
		}

		for _, ld := range te.LinkDatas {
			ld.execEnv = te
			ld.link = flowInst.Flow.GetLink(ld.linkID)
		}
	}
}

// FindOrCreateTaskData finds an existing TaskData or creates ones if not found for the
// specified task the task environment
func (te *ExecEnv) FindOrCreateTaskData(task *definition.Task) (taskData *TaskData, created bool) {

	taskData, ok := te.TaskDatas[task.ID()]

	created = false

	if !ok {
		taskData = NewTaskData(te, task)
		te.TaskDatas[task.ID()] = taskData
		te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

		created = true
	}

	return taskData, created
}

// NewTaskData creates a new TaskData object
func (te *ExecEnv) NewTaskData(task *definition.Task) *TaskData {

	taskData := NewTaskData(te, task)
	te.TaskDatas[task.ID()] = taskData
	te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

	return taskData
}

// FindOrCreateLinkData finds an existing LinkData or creates ones if not found for the
// specified link the task environment
func (te *ExecEnv) FindOrCreateLinkData(link *definition.Link) (linkData *LinkData, created bool) {

	linkData, ok := te.LinkDatas[link.ID()]
	created = false

	if !ok {
		linkData = NewLinkData(te, link)
		te.LinkDatas[link.ID()] = linkData
		te.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtAdd, ID: link.ID(), LinkData: linkData})
		created = true
	}

	return linkData, created
}

// releaseTask cleans up TaskData in the task environment any of its dependencies.
// This is called when a task is completed and can be discarded
func (te *ExecEnv) releaseTask(task *definition.Task) {
	delete(te.TaskDatas, task.ID())
	te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtDel, ID: task.ID()})
	links := task.FromLinks()

	for _, link := range links {
		delete(te.LinkDatas, link.ID())
		te.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtDel, ID: link.ID()})
	}
}

func (te *ExecEnv) FlowDefinition() *definition.Definition {
	return te.flowDef
}

// TaskInsts get the task instances
func (te *ExecEnv)  TaskInsts() []model.TaskInst {
	insts := make([]model.TaskInst, len(te.TaskDatas))

	for _, taskData := range te.TaskDatas {
		insts = append(insts, taskData)
	}

	return insts
}

//Status gets the state of the Flow instance
func (te *ExecEnv) State() int {
	return te.Instance.State()
}

//SetStatus sets the state of the Flow instance
func (te *ExecEnv) SetState(state int) {
	te.Instance.SetState(state)
}