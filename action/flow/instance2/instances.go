package instance2

import (
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type IndependentInstance struct {
	*Instance

	id          string
	stepID      int
	WorkItemQueue *util.SyncQueue //todo: change to faster non-threadsafe queue
	wiCounter     int

	subflowCtr  int
}

type EmbeddedInstance struct {
	*Instance

	instId int
	parentId int
	master *Instance //to access queue
	parent *Instance // could change to "container" and move to instance
}

// New creates a new Flow Instance from the specified Flow
func NewIndependentInstance(instanceID string, flow *definition.Definition) *IndependentInstance {
	var inst IndependentInstance
	inst.id = instanceID
	inst.stepID = 0
	inst.WorkItemQueue = util.NewSyncQueue()

	inst.status = model.FlowStatusNotStarted
	inst.ChangeTracker = NewInstanceChangeTracker()

	return &inst
}

func (inst *IndependentInstance) Start(startAttrs []*data.Attribute) bool {

	if inst.Attrs == nil {
		inst.Attrs = make(map[string]*data.Attribute)
	}

	for _, attr := range startAttrs {
		inst.Attrs[attr.Name()] = attr
	}

	return inst.startInstance(inst.Instance)
}


// StepID returns the current step ID of the Flow Instance
func (inst *IndependentInstance) StepID() int {
	return inst.stepID
}

func (inst *IndependentInstance) DoStep() bool {

	hasNext := false

	inst.ResetChanges()

	inst.stepID++

	if inst.status == model.FlowStatusActive {

		item, ok := inst.WorkItemQueue.Pop()

		if ok {
			logger.Debug("retrieved item from flow instance work queue")

			workItem := item.(*WorkItem)

			behavior := inst.flowModel.GetDefaultTaskBehavior()
			if typeID := workItem.TaskData.task.TypeID(); typeID > 1 {
				behavior = inst.flowModel.GetTaskBehavior(typeID)
			}

			workItem.TaskData.inst.ResetChanges()
			inst.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtDel, ID: workItem.ID, WorkItem: workItem})
			inst.execTask(behavior, workItem.TaskData)
			hasNext = true
		} else {
			logger.Debug("flow instance work queue empty")
		}
	}

	return hasNext
}

func (inst *IndependentInstance) scheduleEval(taskData *TaskData) {

	inst.wiCounter++

	workItem := NewWorkItem(inst.wiCounter, taskData)
	logger.Debugf("Scheduling task: %s\n", taskData.task.Name())

	inst.WorkItemQueue.Push(workItem)
	inst.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtAdd, ID: workItem.ID, WorkItem: workItem})
}

// execTask executes the specified Work Item of the Flow Instance
func (inst *IndependentInstance) execTask(behavior model.TaskBehavior, taskData *TaskData) {

	defer func() {
		if r := recover(); r != nil {

			err := fmt.Errorf("Unhandled Error executing task '%s' : %v\n", taskData.task.Name(), r)
			logger.Error(err)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			if !taskData.inst.isErrorHandler {

				taskData.inst.appendErrorData(NewActivityEvalError(taskData.task.Name(), "unhandled", err.Error()))
				inst.HandleGlobalError(taskData.inst)
			}
			// else what should we do?
		}
	}()

	var err error

	var evalResult model.EvalResult

	if taskData.status == model.TaskStatusWaiting {

		evalResult, err = behavior.PostEval(taskData)

	} else {
		evalResult, err = behavior.Eval(taskData)
	}

	if err != nil {
		inst.handleTaskError(behavior, taskData, err)
		return
	}

	if evalResult == model.EVAL_DONE {
		//task was done
		inst.handleTaskDone(behavior, taskData)
	} else if evalResult == model.EVAL_REPEAT {
		//task needs to iterate or retry
		inst.scheduleEval(taskData)
	}
}

// handleTaskDone handles the completion of a task in the Flow Instance
func (inst *IndependentInstance) handleTaskDone(taskBehavior model.TaskBehavior, taskData *TaskData) {

	notifyFlow, taskEntries, err := taskBehavior.Done(taskData)

	containerInst := taskData.inst

	if err != nil {
		containerInst.appendErrorData(err)
		inst.HandleGlobalError(containerInst)
		return
	}

	flowDone := false
	task := taskData.Task()

	if notifyFlow {

		flowBehavior := containerInst.flowModel.GetFlowBehavior()
		flowDone = flowBehavior.TaskDone(containerInst)
	}

	if flowDone || containerInst.forceCompletion {
		//flow completed or return was called explicitly, so lets complete the flow
		flowBehavior := containerInst.flowModel.GetFlowBehavior()
		flowBehavior.Done(containerInst)
		flowDone = true
		containerInst.SetStatus(model.FlowStatusCompleted)

		//if error flow, return
		//else if containerInst != inst
		//  notify activity that flow is done (schedule post eval)
		//  in top level case inform action -- copy return values to activity output


	} else {
		inst.enterTasks(containerInst, taskEntries)
	}

	containerInst.releaseTask(task)
}

// handleTaskError handles the completion of a task in the Flow Instance
func (inst *IndependentInstance) handleTaskError(taskBehavior model.TaskBehavior, taskData *TaskData, err error) {

	handled, taskEntries := taskBehavior.Error(taskData, err)

	containerInst := taskData.inst

	if !handled {
		if containerInst.isErrorHandler {
			//fail
			inst.SetStatus(model.FlowStatusFailed)
		} else {
			containerInst.appendErrorData(err)
			inst.HandleGlobalError(containerInst)
		}
		return
	}

	if len(taskEntries) != 0 {
		inst.enterTasks(containerInst, taskEntries)
	}

	containerInst.releaseTask(taskData.Task())
}

// HandleGlobalError handles instance errors
func (inst *IndependentInstance) HandleGlobalError(containerInst *Instance) {

	if containerInst.isErrorHandler {
		//todo: log error information
		containerInst.SetStatus(model.FlowStatusFailed)
		return
	}

	//not an error handler, so we should create the error handler instance and start it.
	if containerInst.flowDef.GetErrorHandlerFlow() != nil {

		// todo: should we clear out the existing workitem queue for items from containerInst?

		errorInst := containerInst.NewEmbeddedInstance(containerInst.flowDef.GetErrorHandlerFlow())
		inst.startInstance(errorInst.Instance)
	}
}

func (inst *IndependentInstance) startInstance(toStart *Instance) bool {

	toStart.SetStatus(model.FlowStatusActive)

	//if pi.Attrs == nil {
	//	pi.Attrs = make(map[string]*data.Attribute)
	//}
	//
	//for _, attr := range startAttrs {
	//	pi.Attrs[attr.Name()] = attr
	//}

	//logger.Infof("FlowInstance Flow: %v", pi.FlowModel)

	//need input mappings

	flowBehavior := inst.flowModel.GetFlowBehavior()
	ok, taskEntries := flowBehavior.Start(toStart)

	if ok {
		inst.enterTasks(toStart, taskEntries)
	}

	return ok
}

func (inst *IndependentInstance) enterTasks(activeInst *Instance, taskEntries []*model.TaskEntry) {

	for _, taskEntry := range taskEntries {

		logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
		taskToEnterBehavior := inst.flowModel.GetTaskBehavior(taskEntry.Task.TypeID())

		enterTaskData, _ := activeInst.FindOrCreateTaskData(taskEntry.Task)

		enterResult := taskToEnterBehavior.Enter(enterTaskData)

		if enterResult == model.ENTER_EVAL {
			inst.scheduleEval(enterTaskData)
		} else if enterResult == model.ENTER_EVAL {
			//skip task
		}
	}
}

//////////////////////////////////////////////////////////////////

// WorkItem describes an item of work (event for a Task) that should be executed on Step
type WorkItem struct {
	ID       int       `json:"id"`
	TaskData *TaskData `json:"-"`

	//TaskID string `json:"taskID"`
	//InstID int `json:"instID"`
}

// NewWorkItem constructs a new WorkItem for the specified TaskData
func NewWorkItem(id int, taskData *TaskData) *WorkItem {

	var workItem WorkItem

	workItem.ID = id
	workItem.TaskData = taskData

	//workItem.TaskID = taskData.task.ID()

	return &workItem
}


func NewActivityEvalError(taskName string, errorType string, errorText string) *ActivityEvalError {
	return &ActivityEvalError{taskName: taskName, errType: errorType, errText: errorText}
}

type ActivityEvalError struct {
	taskName string
	errType  string
	errText  string
}

func (e *ActivityEvalError) TaskName() string {
	return e.taskName
}

func (e *ActivityEvalError) Type() string {
	return e.errType
}

func (e *ActivityEvalError) Error() string {
	return e.errText
}
