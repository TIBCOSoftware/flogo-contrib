package instance

import (
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

type IndependentInstance struct {
	*Instance

	id            string
	stepID        int
	workItemQueue *util.SyncQueue //todo: change to faster non-threadsafe queue
	wiCounter     int

	ChangeTracker *InstanceChangeTracker

	subFlowCtr  int
	flowModel   *model.FlowModel
	patch       *support.Patch
	interceptor *support.Interceptor

	subFlows map[int]*EmbeddedInstance
}

type EmbeddedInstance struct {
	*Instance

	instId   int
	parentId int
	//master *Instance //to access queue
	parent *Instance // could change to "container" and move to instance
}

func NewEmbeddedInstance() {
	//ref to the flow
	//Host
	//set pass along inputs

}

// New creates a new Flow Instance from the specified Flow
func NewIndependentInstance(instanceID string, flow *definition.Definition) *IndependentInstance {
	var inst IndependentInstance
	inst.id = instanceID
	inst.stepID = 0
	inst.workItemQueue = util.NewSyncQueue()
	inst.flowDef = flow

	inst.status = model.FlowStatusNotStarted
	inst.ChangeTracker = NewInstanceChangeTracker()

	inst.taskInsts = make(map[string]*TaskInst)
	inst.linkInsts = make(map[int]*LinkInst)

	return &inst
}

func (inst *IndependentInstance) NewEmbeddedInstance(containerInst *Instance, flow *definition.Definition) *EmbeddedInstance {

	inst.subFlowCtr++

	var embeddedInst EmbeddedInstance
	embeddedInst.flowDef = flow
	embeddedInst.subFlowId = inst.subFlowCtr
	embeddedInst.master = inst

	if flow.IsErrorHandler() {

	} else {

	}

	embeddedInst.parent = containerInst

	if inst.subFlows == nil {
		inst.subFlows = make(map[int]*EmbeddedInstance)
	}
	inst.subFlows[embeddedInst.subFlowId] = &embeddedInst

	inst.ChangeTracker.SubFlowChange(containerInst.subFlowId, CtAdd, embeddedInst.subFlowId, "")

	return &embeddedInst
}

// ID returns the ID of the Flow Instance
func (inst *IndependentInstance) ID() string {
	return inst.id
}

func (inst *IndependentInstance) Start(startAttrs []*data.Attribute) bool {

	if inst.attrs == nil {
		inst.attrs = make(map[string]*data.Attribute)
	}

	for _, attr := range startAttrs {
		inst.attrs[attr.Name()] = attr
	}

	return inst.startInstance(inst.Instance)
}

func (inst *IndependentInstance) ApplyPatch(patch *support.Patch) {
	if inst.patch == nil {
		inst.patch = patch
		inst.patch.Init()
	}
}

func (inst *IndependentInstance) ApplyInterceptor(interceptor *support.Interceptor) {
	if inst.interceptor == nil {
		inst.interceptor = interceptor
		inst.interceptor.Init()
	}
}

// GetChanges returns the Change Tracker object
func (inst *IndependentInstance) GetChanges() *InstanceChangeTracker {
	return inst.ChangeTracker
}

// ResetChanges resets an changes that were being tracked
func (inst *IndependentInstance) ResetChanges() {

	if inst.ChangeTracker != nil {
		inst.ChangeTracker.ResetChanges()
	}

	//todo: can we reuse this to avoid gc
	inst.ChangeTracker = NewInstanceChangeTracker()
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

		item, ok := inst.workItemQueue.Pop()

		if ok {
			logger.Debug("retrieved item from flow instance work queue")

			workItem := item.(*WorkItem)

			behavior := inst.flowModel.GetDefaultTaskBehavior()
			if typeID := workItem.taskInst.task.TypeID(); typeID > 1 {
				behavior = inst.flowModel.GetTaskBehavior(typeID)
			}

			inst.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtDel, ID: workItem.ID, WorkItem: workItem})
			inst.execTask(behavior, workItem.taskInst)
			hasNext = true
		} else {
			logger.Debug("flow instance work queue empty")
		}
	}

	return hasNext
}

func (inst *IndependentInstance) scheduleEval(taskInst *TaskInst) {

	inst.wiCounter++

	workItem := NewWorkItem(inst.wiCounter, taskInst)
	logger.Debugf("Scheduling task: %s\n", taskInst.task.Name())

	inst.workItemQueue.Push(workItem)
	inst.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtAdd, ID: workItem.ID, WorkItem: workItem})
}

// execTask executes the specified Work Item of the Flow Instance
func (inst *IndependentInstance) execTask(behavior model.TaskBehavior, taskInst *TaskInst) {

	defer func() {
		if r := recover(); r != nil {

			err := fmt.Errorf("Unhandled Error executing task '%s' : %v\n", taskInst.task.Name(), r)
			logger.Error(err)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			if !taskInst.flowInst.isErrorHandler {

				taskInst.flowInst.appendErrorData(NewActivityEvalError(taskInst.task.Name(), "unhandled", err.Error()))
				inst.HandleGlobalError(taskInst.flowInst)
			}
			// else what should we do?
		}
	}()

	var err error

	var evalResult model.EvalResult

	if taskInst.status == model.TaskStatusWaiting {

		evalResult, err = behavior.PostEval(taskInst)

	} else {
		evalResult, err = behavior.Eval(taskInst)
	}

	if err != nil {
		inst.handleTaskError(behavior, taskInst, err)
		return
	}

	if evalResult == model.EVAL_DONE {
		//task was done
		inst.handleTaskDone(behavior, taskInst)
	} else if evalResult == model.EVAL_REPEAT {
		//task needs to iterate or retry
		inst.scheduleEval(taskInst)
	}
}

// handleTaskDone handles the completion of a task in the Flow Instance
func (inst *IndependentInstance) handleTaskDone(taskBehavior model.TaskBehavior, taskInst *TaskInst) {

	notifyFlow, taskEntries, err := taskBehavior.Done(taskInst)

	containerInst := taskInst.flowInst

	if err != nil {
		containerInst.appendErrorData(err)
		inst.HandleGlobalError(containerInst)
		return
	}

	flowDone := false
	task := taskInst.Task()

	if notifyFlow {

		flowBehavior := inst.flowModel.GetFlowBehavior()
		flowDone = flowBehavior.TaskDone(containerInst)
	}

	if flowDone || containerInst.forceCompletion {
		//flow completed or return was called explicitly, so lets complete the flow
		flowBehavior := inst.flowModel.GetFlowBehavior()
		flowBehavior.Done(containerInst)
		flowDone = true
		containerInst.SetStatus(model.FlowStatusCompleted)

		// if error flow, we should mark the parent instance as done
		// else if containerInst != inst
		//   if flow was started by an activity, we should schedule the post-eval
		//   -- copy return values to activity output

	} else {
		inst.enterTasks(containerInst, taskEntries)
	}

	//inst.releaseTask(taskInst)
	containerInst.releaseTask(task)
}

//func (inst *IndependentInstance) releaseTask(taskInst *TaskInst) {
//
//	task := taskInst.Task()
//
//	delete(taskInst.inst.TaskDatas, task.ID())
//
//	inst.ChangeTracker.trackTaskData(&TaskInstChange{ChgType: CtDel, SubFlowID: taskInst.inst.subFlowId, ID: task.ID()})
//	links := task.FromLinks()
//
//	for _, link := range links {
//		delete(taskInst.inst.LinkDatas, link.ID())
//		inst.ChangeTracker.trackLinkData(&LinkInstChange{ChgType: CtDel, SubFlowID: taskInst.inst.subFlowId,ID: link.ID()})
//	}
//}

// handleTaskError handles the completion of a task in the Flow Instance
func (inst *IndependentInstance) handleTaskError(taskBehavior model.TaskBehavior, taskInst *TaskInst, err error) {

	handled, taskEntries := taskBehavior.Error(taskInst, err)

	containerInst := taskInst.flowInst

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

	containerInst.releaseTask(taskInst.Task())
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

		errorInst := inst.NewEmbeddedInstance(containerInst, containerInst.flowDef.GetErrorHandlerFlow())
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
	taskInst *TaskInst `json:"-"`

	TaskID    string `json:"taskID"`
	SubFlowID int    `json:"subFlowId"`
}

// NewWorkItem constructs a new WorkItem for the specified TaskInst
func NewWorkItem(id int, taskInst *TaskInst) *WorkItem {

	var workItem WorkItem

	workItem.ID = id
	workItem.taskInst = taskInst
	workItem.TaskID = taskInst.task.ID()
	workItem.SubFlowID = taskInst.flowInst.subFlowId

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
