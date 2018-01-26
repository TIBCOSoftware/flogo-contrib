package instance

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/provider"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
)

const (
	idEhExecEnv   = 0
	idFlowExecEnv = 1
)

// Instance is a structure for representing an instance of a Flow
type Instance struct {
	id            string
	stepID        int
	lock          sync.Mutex
	status        Status
	state         int
	FlowURI       string
	Flow          *definition.Definition
	FlowExecEnv   *ExecEnv
	EhFlowExecEnv *ExecEnv
	FlowModel     *model.FlowModel
	Attrs         map[string]*data.Attribute
	Patch         *support.Patch
	Interceptor   *support.Interceptor

	WorkItemQueue *util.SyncQueue //todo: change to faster non-threadsafe queue

	wiCounter     int
	ChangeTracker *InstanceChangeTracker `json:"-"`

	flowProvider provider.Provider
	replyHandler activity.ReplyHandler
	actionCtx    *ActionCtx //todo after transition to actionCtx, make sure actionCtx isn't null before executing

	forceCompletion bool
	returnData      map[string]*data.Attribute
	returnError     error
}

// New creates a new Flow Instance from the specified Flow
func New(instanceID string, flowURI string, flow *definition.Definition, flowModel *model.FlowModel) *Instance {
	var inst Instance
	inst.id = instanceID
	inst.stepID = 0
	inst.FlowURI = flowURI
	inst.Flow = flow
	inst.FlowModel = flowModel
	inst.status = StatusNotStarted
	inst.WorkItemQueue = util.NewSyncQueue()
	inst.ChangeTracker = NewInstanceChangeTracker()

	var execEnv ExecEnv
	execEnv.ID = idFlowExecEnv
	//execEnv.Task = flow.RootTask()
	//execEnv.taskID = flow.RootTask().ID()
	execEnv.Instance = &inst
	execEnv.TaskDatas = make(map[string]*TaskData)
	execEnv.LinkDatas = make(map[int]*LinkData)

	inst.FlowExecEnv = &execEnv

	return &inst
}

// SetFlowProvider sets the process.Provider that the instance should use
func (pi *Instance) SetFlowProvider(provider provider.Provider) {
	pi.flowProvider = provider
}

// Restart indicates that this FlowInstance was restarted
func (pi *Instance) Restart(id string, provider provider.Provider) {
	pi.id = id
	pi.flowProvider = provider
	pi.Flow, _ = pi.flowProvider.GetFlow(pi.FlowURI)
	pi.FlowModel = model.Get(pi.Flow.ModelID())
	pi.FlowExecEnv.init(pi)
}

// ID returns the ID of the Flow Instance
func (pi *Instance) ID() string {
	return pi.id
}

// Name implements activity.FlowDetails.Name method
func (pi *Instance) Name() string {
	return pi.Flow.Name()
}

// ReplyHandler returns the reply handler for the instance
func (pi *Instance) ReplyHandler() activity.ReplyHandler {
	return &SimpleReplyHandler{pi.actionCtx.rh}
}

// SimpleReplyHandler is a simple ReplyHandler that is pass-thru to the action ResultHandler
type SimpleReplyHandler struct {
	resultHandler action.ResultHandler
}

// Reply implements ReplyHandler.Reply
func (rh *SimpleReplyHandler) Reply(code int, replyData interface{}, err error) {

	dataAttr, _ := data.NewAttribute("data", data.OBJECT, replyData)
	codeAttr, _ := data.NewAttribute("code", data.INTEGER, code)
	resultData := map[string]*data.Attribute{
		"data": dataAttr,
		"code": codeAttr,
	}

	rh.resultHandler.HandleResult(resultData, err)
}

// InitActionContext initialize the action context, should be initialized before execution
func (pi *Instance) InitActionContext(config *action.Config, handler action.ResultHandler) {
	pi.actionCtx = &ActionCtx{inst: pi, config: config, rh: handler}
}

// FlowDefinition returns the Flow that the instance is of
func (pi *Instance) FlowDefinition() *definition.Definition {
	return pi.Flow
}

// StepID returns the current step ID of the Flow Instance
func (pi *Instance) StepID() int {
	return pi.stepID
}

// Status returns the current status of the Flow Instance
func (pi *Instance) Status() Status {
	return pi.status
}

func (pi *Instance) setStatus(status Status) {

	pi.status = status
	pi.ChangeTracker.SetStatus(status)
}

// Status returns the state indicator of the Flow Instance
func (pi *Instance) State() int {
	return pi.state
}

// SetStatus sets the state indicator of the Flow Instance
func (pi *Instance) SetState(state int) {
	pi.state = state
	pi.ChangeTracker.SetState(state)
}

// TaskInsts get the task instances
func (pi *Instance) TaskInsts() []model.TaskInst {
	return nil
}

// UpdateAttrs updates the attributes of the Flow Instance
func (pi *Instance) UpdateAttrs(attrs []*data.Attribute) {

	if attrs != nil {

		logger.Debugf("Updating flow attrs: %v", attrs)

		if pi.Attrs == nil {
			pi.Attrs = make(map[string]*data.Attribute, len(attrs))
		}

		for _, attr := range attrs {
			pi.Attrs[attr.Name()] = attr
		}
	}
}

// Start will start the Flow Instance, returns a boolean indicating
// if it was able to start
func (pi *Instance) Start(startAttrs []*data.Attribute) bool {

	pi.setStatus(StatusActive)

	if pi.Attrs == nil {
		pi.Attrs = make(map[string]*data.Attribute)
	}

	for _, attr := range startAttrs {
		pi.Attrs[attr.Name()] = attr
	}

	logger.Infof("FlowInstance Flow: %v", pi.FlowModel)
	flowBehavior := pi.FlowModel.GetFlowBehavior()

	//todo: error if flowBehavior not found

	ok, taskEntries := flowBehavior.Start(pi)

	if ok {

		for _, taskEntry := range taskEntries {

			logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
			taskToEnterBehavior := pi.FlowModel.GetTaskBehavior(taskEntry.Task.TypeID())

			enterTaskData, _ := pi.FlowExecEnv.FindOrCreateTaskData(taskEntry.Task)

			eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

			if eval {
				pi.scheduleEval(enterTaskData, evalCode)
			}
		}

		//rootTaskData := pi.FlowExecEnv.NewTaskData(pi.Flow.RootTask())
		//pi.scheduleEval(rootTaskData, evalCode)
	}

	return ok
}

////Resume resumes a Flow Instance
//func (pi *Instance) Resume(flowData map[string]interface{}) bool {
//
//	model := pi.FlowModel.GetFlowBehavior(pi.Flow.TypeID())
//
//	pi.setStatus(StatusActive)
//	pi.UpdateAttrs(flowData)
//
//	return model.Resume(pi)
//}

// DoStep performs a single execution 'step' of the Flow Instance
func (pi *Instance) DoStep() bool {

	hasNext := false

	pi.ResetChanges()

	pi.stepID++

	if pi.status == StatusActive {

		item, ok := pi.WorkItemQueue.Pop()

		if ok {
			logger.Debug("popped item off queue")

			workItem := item.(*WorkItem)

			pi.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtDel, ID: workItem.ID, WorkItem: workItem})

			pi.execTask(workItem)
			hasNext = true
		} else {
			logger.Debug("queue emtpy")
		}
	}

	return hasNext
}

// GetChanges returns the Change Tracker object
func (pi *Instance) GetChanges() *InstanceChangeTracker {
	return pi.ChangeTracker
}

// ResetChanges resets an changes that were being tracked
func (pi *Instance) ResetChanges() {

	if pi.ChangeTracker != nil {
		pi.ChangeTracker.ResetChanges()
	}

	//todo: can we reuse this to avoid gc
	pi.ChangeTracker = NewInstanceChangeTracker()
}

func (pi *Instance) scheduleEval(taskData *TaskData, evalCode int) {

	pi.wiCounter++

	workItem := NewWorkItem(pi.wiCounter, taskData, EtEval, evalCode)
	logger.Debugf("Scheduling EVAL on task: %s\n", taskData.task.Name())

	pi.WorkItemQueue.Push(workItem)
	pi.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtAdd, ID: workItem.ID, WorkItem: workItem})
}

// execTask executes the specified Work Item of the Flow Instance
func (pi *Instance) execTask(workItem *WorkItem) {

	defer func() {
		if r := recover(); r != nil {

			err := fmt.Errorf("Unhandled Error executing task '%s' : %v\n", workItem.TaskData.task.Name(), r)
			logger.Error(err)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

//<<<<<<< HEAD
			pi.appendErrorData(NewActivityEvalError(workItem.TaskData.task.Name(), "unhandled", err.Error()))
			if workItem.TaskData.execEnv.ID != idEhExecEnv {
//=======
//			pi.appendActivityErrorData(workItem.TaskData, activity.NewError(err.Error(), "", nil))
//			if workItem.TaskData.execEnv.ID != idEhExecEnv {
//>>>>>>> initial sub-flow support
				//not already in global handler, so handle it
				pi.HandleGlobalError()
			}
		}
	}()

	taskData := workItem.TaskData
	taskBehavior := pi.FlowModel.GetTaskBehavior(taskData.task.TypeID())

	var done bool
	var doneCode int
	var err error

	//todo: should validate process activities

	var evalResult model.EvalResult

	if workItem.ExecType == EtEval {

		evalResult, doneCode, err = taskBehavior.Eval(taskData, workItem.EvalCode)

	} else {
		done, doneCode, err = taskBehavior.PostEval(taskData, workItem.EvalCode, nil)
		if done {
			evalResult = model.EVAL_DONE
		} else {
			evalResult = model.EVAL_WAIT
		}
	}

	if err != nil {
		pi.handleTaskError(taskBehavior, taskData, err)
		return
	}

	if evalResult == model.EVAL_DONE {
		pi.handleTaskDone(taskBehavior, taskData, doneCode)
	} else if evalResult == model.EVAL_REPEAT {
		//iterate or retry
		pi.scheduleEval(taskData, workItem.EvalCode)
	}
}

// handleTaskDone handles the completion of a task in the Flow Instance
func (pi *Instance) handleTaskDone(taskBehavior model.TaskBehavior, taskData *TaskData, doneCode int) {

	notifyFlow, notifyCode, taskEntries, err := taskBehavior.Done(taskData, doneCode)

	if err != nil {
		pi.appendErrorData(err)
		pi.HandleGlobalError()
		return
	}

	flowDone := false
	task := taskData.Task()

	if notifyFlow {

		flowBehavior := pi.FlowModel.GetFlowBehavior()
		flowDone = flowBehavior.TaskDone(pi, notifyCode)
	}

	if flowDone || pi.forceCompletion {
		//flow completed or return was called explicitly, so lets complete the flow
		flowBehavior := pi.FlowModel.GetFlowBehavior()
		flowBehavior.Done(pi)
		flowDone = true
		pi.setStatus(StatusCompleted)

	} else {
		for _, taskEntry := range taskEntries {

			logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
			taskToEnterBehavior := pi.FlowModel.GetTaskBehavior(taskEntry.Task.TypeID())

			enterTaskData, _ := taskData.execEnv.FindOrCreateTaskData(taskEntry.Task)

			eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

			if eval {
				pi.scheduleEval(enterTaskData, evalCode)
			}
		}
	}

	taskData.execEnv.releaseTask(task)
}

func (pi *Instance) appendErrorData(err error) {

	switch e := err.(type) {
	case *definition.LinkExprError:
		pi.AddAttr("{Error.type}", data.STRING, "link_expr")
		pi.AddAttr("{Error.message}", data.STRING, err.Error())
	case *activity.Error:
		pi.AddAttr("{Error.message}", data.STRING, err.Error())
		pi.AddAttr("{Error.data}", data.OBJECT, e.Data())
		pi.AddAttr("{Error.code}", data.STRING, e.Code())

		if e.ActivityName() != "" {
			pi.AddAttr("{Error.activity}", data.STRING, e.ActivityName())
		}
	case *ActivityEvalError:
		pi.AddAttr("{Error.activity}", data.STRING, e.TaskName())
		pi.AddAttr("{Error.message}", data.STRING, err.Error())
		pi.AddAttr("{Error.type}", data.STRING, e.Type())
	default:
		pi.AddAttr("{Error.message}", data.STRING, err.Error())
	}

	//todo add case for *dataMapperError & *activity.Error
}

// handleTaskError handles the completion of a task in the Flow Instance
func (pi *Instance) handleTaskError(taskBehavior model.TaskBehavior, taskData *TaskData, err error) {

	handled, taskEntry := taskBehavior.Error(taskData)

	if !handled {
		pi.appendErrorData(err)
		if taskData.execEnv.ID != idEhExecEnv {
			pi.HandleGlobalError()
		}
		return
	}

	//todo add error data for task to flow

	if taskEntry != nil {

		logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
		taskToEnterBehavior := pi.FlowModel.GetTaskBehavior(taskEntry.Task.TypeID())

		enterTaskData, _ := taskData.execEnv.FindOrCreateTaskData(taskEntry.Task)

		eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

		if eval {
			pi.scheduleEval(enterTaskData, evalCode)
		}
	}

	task := taskData.Task()
	taskData.execEnv.releaseTask(task)
}

// HandleGlobalError handles instance errors
func (pi *Instance) HandleGlobalError() {

	if pi.Flow.GetErrorHandlerFlow() != nil {

		// if embedded run in the context of this instance
		// special case, when this sub-flow is done, process is done
		// todo: should we clear out the existing workitem queue?

		ehFlow := pi.Flow.GetErrorHandlerFlow()

		if pi.EhFlowExecEnv == nil {
			var execEnv ExecEnv
			execEnv.ID = idEhExecEnv
			//execEnv.Task = ehTask
			//execEnv.taskID = ehTask.ID()
			execEnv.Instance = pi
			execEnv.TaskDatas = make(map[string]*TaskData)
			execEnv.LinkDatas = make(map[int]*LinkData)

			pi.EhFlowExecEnv = &execEnv
		}

		flowBehavior := pi.FlowModel.GetFlowBehavior()
		ok, taskEntries := flowBehavior.StartEmbeddedFlow(pi, ehFlow)

		if ok {

			for _, taskEntry := range taskEntries {

				logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
				taskToEnterBehavior := pi.FlowModel.GetTaskBehavior(taskEntry.Task.TypeID())

				enterTaskData, _ := pi.EhFlowExecEnv.FindOrCreateTaskData(taskEntry.Task)

				eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

				if eval {
					pi.scheduleEval(enterTaskData, evalCode)
				}
			}
		}
	} else {

		//todo: log error information
		pi.setStatus(StatusFailed)
	}
}

// GetAttr implements data.Scope.GetAttr
func (pi *Instance) GetAttr(attrName string) (value *data.Attribute, exists bool) {

	if pi.Attrs != nil {
		attr, found := pi.Attrs[attrName]

		if found {
			return attr, true
		}
	}

	return pi.Flow.GetAttr(attrName)
}

func (pi *Instance) getInstAttr(attrName string) (value *data.Attribute, exists bool) {

	if pi.Attrs != nil {
		attr, found := pi.Attrs[attrName]
		return attr, found
	}
	return nil, false
}

// SetAttrValue implements api.Scope.SetAttrValue
func (pi *Instance) SetAttrValue(attrName string, value interface{}) error {
	if pi.Attrs == nil {
		pi.Attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("SetAttr - name: %s, value:%v\n", attrName, value)

	existingAttr, exists := pi.GetAttr(attrName)

	//todo: optimize, use existing attr
	if exists {
		//todo handle error
		attr, _ := data.NewAttribute(attrName, existingAttr.Type(), value)
		pi.Attrs[attrName] = attr
		pi.ChangeTracker.AttrChange(CtUpd, attr)
		return nil
	}

	return fmt.Errorf("Attr [%s] does not exists", attrName)
}

// AddAttr add a new attribute to the instance
func (pi *Instance) AddAttr(attrName string, attrType data.Type, value interface{}) *data.Attribute {
	if pi.Attrs == nil {
		pi.Attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("AddAttr - name: %s, type: %s, value:%v\n", attrName, attrType, value)

	var attr *data.Attribute

	existingAttr, exists := pi.GetAttr(attrName)

	if exists {
		attr = existingAttr
	} else {
		//todo handle error
		attr, _ = data.NewAttribute(attrName, attrType, value)
		pi.Attrs[attrName] = attr
		pi.ChangeTracker.AttrChange(CtAdd, attr)
	}

	return attr
}

func (pi *Instance) ActionContext() action.Context {
	return pi.actionCtx
}

type ActionCtx struct {
	config *action.Config
	inst   *Instance
	rh     action.ResultHandler
}

func (ac *ActionCtx) ID() string {
	return ac.config.Id
}

func (ac *ActionCtx) Ref() string {
	return ac.config.Ref
}

func (ac *ActionCtx) InstanceMetadata() *action.ConfigMetadata {
	return ac.config.Metadata
}

func (ac *ActionCtx) Reply(replyData map[string]*data.Attribute, err error) {
	ac.rh.HandleResult(replyData, err)
}

func (ac *ActionCtx) Return(returnData map[string]*data.Attribute, err error) {
	ac.inst.forceCompletion = true
	ac.inst.returnData = returnData
	ac.inst.returnError = err
}

func (ac *ActionCtx) WorkingData() data.Scope {
	return ac.inst
}

func (ac *ActionCtx) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

func (pi *Instance) GetReturnData() (map[string]*data.Attribute, error) {

	if pi.returnData == nil {

		md := pi.actionCtx.InstanceMetadata()
		//construct returnData from instance attributes

		if md != nil && md.Output != nil {

			pi.returnData = make(map[string]*data.Attribute)
			for _, mdAttr := range md.Output {
				piAttr, exists := pi.Attrs[mdAttr.Name()]
				if exists {
					pi.returnData[piAttr.Name()] = piAttr
				}
			}
		}
	}

	return pi.returnData, pi.returnError
}

// ExecType is the type of execution to perform
type ExecType int

const (
	// EtEval denoted the Eval execution type
	EtEval ExecType = 10

	// EtPostEval denoted the PostEval execution type
	EtPostEval ExecType = 20
)

// WorkItem describes an item of work (event for a Task) that should be executed on Step
type WorkItem struct {
	ID       int       `json:"id"`
	TaskData *TaskData `json:"-"`
	ExecType ExecType  `json:"execType"`
	EvalCode int       `json:"code"`

	TaskID string `json:"taskID"` //for now need for ser
	//taskCtxID int `json:"taskCtxID"` //not needed for now
}

// NewWorkItem constructs a new WorkItem for the specified TaskData
func NewWorkItem(id int, taskData *TaskData, execType ExecType, evalCode int) *WorkItem {

	var workItem WorkItem

	workItem.ID = id
	workItem.TaskData = taskData
	workItem.ExecType = execType
	workItem.EvalCode = evalCode

	workItem.TaskID = taskData.task.ID()

	return &workItem
}
