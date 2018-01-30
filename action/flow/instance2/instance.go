package instance2

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

type Instance struct {
	id          string
	stepID      int
	status      model.FlowStatus
	flowModel   *model.FlowModel
	RootExecEnv *ExecEnv
	execEnvId   int

	WorkItemQueue *util.SyncQueue //todo: change to faster non-threadsafe queue
	wiCounter     int
	ChangeTracker *InstanceChangeTracker

	actionCtx *ActionCtx //todo after transition to actionCtx, make sure actionCtx isn't null before executing

	instId int
	parentId int

	//ret to owner/parent //to tell owner you are done or failed
	//ref to master

}

// New creates a new Flow Instance from the specified Flow
func New(instanceID string, flow *definition.Definition, flowModel *model.FlowModel) *Instance {
	var inst Instance
	inst.id = instanceID
	inst.stepID = 0
	inst.flowModel = flowModel
	inst.status = model.FlowStatusNotStarted

	inst.WorkItemQueue = util.NewSyncQueue()
	//inst.ChangeTracker = NewInstanceChangeTracker()

	inst.RootExecEnv = inst.NewExecEnv(nil, flow)

	return &inst
}

func (inst *Instance) NewEmbeddedInstance(flow *definition.Definition) *Instance {

	return inst.id
}

// ID returns the ID of the Flow Instance
func (inst *Instance) ID() string {
	return inst.id
}

// StepID returns the current step ID of the Flow Instance
func (inst *Instance) StepID() int {
	return inst.stepID
}

// Status returns the current status of the Flow Instance
func (inst *Instance) Status() model.FlowStatus {
	return inst.status
}

func (inst *Instance) SetStatus(status model.FlowStatus) {

	inst.status = status
	//inst.ChangeTracker.SetStatus(status)
}

func (inst *Instance) DoStep() bool {

	hasNext := false

	//inst.ResetChanges()

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
			//inst.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtDel, ID: workItem.ID, WorkItem: workItem})
			ExecTask(behavior, workItem.TaskData)
			//inst.execTask(workItem)
			hasNext = true
		} else {
			logger.Debug("flow instance work queue empty")
		}
	}

	return hasNext
}

func (inst *Instance) scheduleEval(taskData *TaskData, evalCode int) {

	inst.wiCounter++

	workItem := NewWorkItem(inst.wiCounter, taskData, evalCode)
	logger.Debugf("Scheduling task: %s\n", taskData.task.Name())

	inst.WorkItemQueue.Push(workItem)
	//inst.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtAdd, ID: workItem.ID, WorkItem: workItem})
}

//////////////////////////////////////////////////////////////
// ActionCtx

func (inst *Instance) ActionContext() action.Context {
	return inst.actionCtx
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
	//	ac.inst.forceCompletion = true
	//	ac.inst.returnData = returnData
	//	ac.inst.returnError = err
}

func (ac *ActionCtx) WorkingData() data.Scope {
	return nil //ac.inst
}

func (ac *ActionCtx) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

// WorkItem describes an item of work (event for a Task) that should be executed on Step
type WorkItem struct {
	ID       int       `json:"id"`
	TaskData *TaskData `json:"-"`
	//ExecType ExecType  `json:"execType"`
	EvalCode int `json:"code"`

	//TaskID string `json:"taskID"` //for now need for ser
	//taskCtxID int `json:"taskCtxID"` //not needed for now
}

// NewWorkItem constructs a new WorkItem for the specified TaskData
func NewWorkItem(id int, taskData *TaskData, evalCode int) *WorkItem {

	var workItem WorkItem

	workItem.ID = id
	workItem.TaskData = taskData
	//workItem.ExecType = execType
	workItem.EvalCode = evalCode

	//workItem.TaskID = taskData.task.ID()

	return &workItem
}
