package instance2

import (
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Execution Environment

func (inst *Instance) NewExecEnv(parentEnv *ExecEnv, flow *definition.Definition) *ExecEnv {
	var execEnv *ExecEnv
	execEnv.ID = inst.execEnvId
	execEnv.Instance = inst
	execEnv.flowDef = flow

	execEnv.TaskDatas = make(map[string]*TaskData)
	execEnv.LinkDatas = make(map[int]*LinkData)

	inst.execEnvId++

	return execEnv
}

// ExecEnv is a structure that describes the execution environment for a flow
type ExecEnv struct {
	ParentEnv *ExecEnv
	ErrorEnv  *ExecEnv
	ID        int
	flowDef   *definition.Definition
	Instance  *Instance

	//Task      *definition.Task

	Attrs map[string]*data.Attribute

	TaskDatas map[string]*TaskData
	LinkDatas map[int]*LinkData

	status model.FlowStatus

	isErrorHandler  bool
	forceCompletion bool
	returnData      map[string]*data.Attribute
	returnError     error

	//taskID string // for deserialization
}

// FindOrCreateTaskData finds an existing TaskData or creates ones if not found for the
// specified task the task environment
func (env *ExecEnv) FindOrCreateTaskData(task *definition.Task) (taskData *TaskData, created bool) {

	taskData, ok := env.TaskDatas[task.ID()]

	created = false

	if !ok {
		taskData = NewTaskData(env, task)
		env.TaskDatas[task.ID()] = taskData
		env.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

		created = true
	}

	return taskData, created
}

// FindOrCreateLinkData finds an existing LinkData or creates ones if not found for the
// specified link the task environment
func (env *ExecEnv) FindOrCreateLinkData(link *definition.Link) (linkData *LinkData, created bool) {

	linkData, ok := env.LinkDatas[link.ID()]
	created = false

	if !ok {
		linkData = NewLinkData(env, link)
		env.LinkDatas[link.ID()] = linkData
		env.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtAdd, ID: link.ID(), LinkData: linkData})
		created = true
	}

	return linkData, created
}

func (env *ExecEnv) releaseTask(task *definition.Task) {
	delete(env.TaskDatas, task.ID())
	env.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtDel, ID: task.ID()})
	links := task.FromLinks()

	for _, link := range links {
		delete(env.LinkDatas, link.ID())
		env.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtDel, ID: link.ID()})
	}
}

// execTask executes the specified Work Item of the Flow Instance
func ExecTask(behavior model.TaskBehavior, taskData *TaskData) {

	defer func() {
		if r := recover(); r != nil {

			err := fmt.Errorf("Unhandled Error executing task '%s' : %v\n", taskData.task.Name(), r)
			logger.Error(err)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			//env.appendErrorData(NewActivityEvalError(workItem.TaskData.task.Name(), "unhandled", err.Error()))
			//if workItem.TaskData.execEnv.ID != idEhExecEnv {
			//	//=======
			//	//			env.appendActivityErrorData(workItem.TaskData, activity.NewError(err.Error(), "", nil))
			//	//			if workItem.TaskData.execEnv.ID != idEhExecEnv {
			//	//>>>>>>> initial sub-flow support
			//	//not already in global handler, so handle it
			//	env.HandleGlobalError()
			//}
		}
	}()

	//var done bool
	var doneCode int
	var err error

	//todo: should validate process activities

	var evalResult model.EvalResult

	//if workItem.ExecType == EtEval {
	//

	evalResult, doneCode, err = behavior.Eval(taskData, 0)

	//} else {
	//	done, doneCode, err = taskBehavior.PostEval(taskData, workItem.EvalCode, nil)
	//	if done {
	//		evalResult = model.EVAL_DONE
	//	} else {
	//		evalResult = model.EVAL_WAIT
	//	}
	//}

	env := taskData.execEnv

	if err != nil {
		env.handleTaskError(behavior, taskData, err)
		return
	}

	if evalResult == model.EVAL_DONE {
		//task was done
		env.handleTaskDone(behavior, taskData, doneCode)
	} else if evalResult == model.EVAL_REPEAT {
		//task needs to iterate or retry
		env.Instance.scheduleEval(taskData, 0)
	}
}

// handleTaskDone handles the completion of a task in the Flow Instance
func (env *ExecEnv) handleTaskDone(taskBehavior model.TaskBehavior, taskData *TaskData, doneCode int) {

	notifyFlow, notifyCode, taskEntries, err := taskBehavior.Done(taskData, doneCode)

	if err != nil {
		env.appendErrorData(err)
		env.HandleGlobalError()
		return
	}

	flowDone := false
	task := taskData.Task()

	if notifyFlow {

		flowBehavior := env.Instance.flowModel.GetFlowBehavior()
		flowDone = flowBehavior.TaskDone(env, notifyCode)
	}

	if flowDone || env.forceCompletion {
		//flow completed or return was called explicitly, so lets complete the flow
		flowBehavior := env.Instance.flowModel.GetFlowBehavior()
		flowBehavior.Done(env)
		flowDone = true
		env.setStatus(model.FlowStatusCompleted)

	} else {
		for _, taskEntry := range taskEntries {

			logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
			taskToEnterBehavior := env.Instance.flowModel.GetTaskBehavior(taskEntry.Task.TypeID())

			enterTaskData, _ := taskData.execEnv.FindOrCreateTaskData(taskEntry.Task)

			eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

			if eval {
				env.Instance.scheduleEval(enterTaskData, evalCode)
			}
		}
	}

	env.releaseTask(task)
}

// handleTaskError handles the completion of a task in the Flow Instance
func (env *ExecEnv) handleTaskError(taskBehavior model.TaskBehavior, taskData *TaskData, err error) {

	handled, taskEntry := taskBehavior.Error(taskData)

	if !handled {
		if env.isErrorHandler {
			//fail
		} else {
			env.appendErrorData(err)
			env.HandleGlobalError()
		}
		return
	}

	if taskEntry != nil {

		logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
		taskToEnterBehavior := env.Instance.flowModel.GetTaskBehavior(taskEntry.Task.TypeID())

		enterTaskData, _ := taskData.execEnv.FindOrCreateTaskData(taskEntry.Task)

		eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

		if eval {
			env.Instance.scheduleEval(enterTaskData, evalCode)
		}
	}

	task := taskData.Task()
	taskData.execEnv.releaseTask(task)
}

// HandleGlobalError handles instance errors
func (env *ExecEnv) HandleGlobalError() {

	if env.flowDef.GetErrorHandlerFlow() != nil {

		// if embedded run in the context of this instance
		// special case, when this sub-flow is done, process is done
		// todo: should we clear out the existing workitem queue for items from this exec env?

		ehFlow := env.flowDef.GetErrorHandlerFlow()

		if env.ErrorEnv == nil {
			var execEnv ExecEnv
			execEnv.ParentEnv = env
			execEnv.isErrorHandler = true
			execEnv.Instance = env.Instance
			execEnv.TaskDatas = make(map[string]*TaskData)
			execEnv.LinkDatas = make(map[int]*LinkData)
			execEnv.flowDef = ehFlow

			env.ErrorEnv = &execEnv
		}

		flowBehavior := env.Instance.flowModel.GetFlowBehavior()
		ok, taskEntries := flowBehavior.StartEmbeddedFlow(env.ErrorEnv, ehFlow)

		if ok {

			for _, taskEntry := range taskEntries {

				logger.Debugf("execTask - TaskEntry: %v\n", taskEntry)
				taskToEnterBehavior := env.Instance.flowModel.GetTaskBehavior(taskEntry.Task.TypeID())

				enterTaskData, _ := env.ErrorEnv.FindOrCreateTaskData(taskEntry.Task)

				eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

				if eval {
					env.Instance.scheduleEval(enterTaskData, evalCode)
				}
			}
		}
	} else {

		//todo: log error information
		env.setStatus(model.FlowStatusFailed)
	}
}

func (env *ExecEnv) appendErrorData(err error) {

	switch e := err.(type) {
	case *definition.LinkExprError:
		env.AddAttr("{Error.type}", data.STRING, "link_expr")
		env.AddAttr("{Error.message}", data.STRING, err.Error())
	case *activity.Error:
		env.AddAttr("{Error.message}", data.STRING, err.Error())
		env.AddAttr("{Error.data}", data.OBJECT, e.Data())
		env.AddAttr("{Error.code}", data.STRING, e.Code())

		if e.ActivityName() != "" {
			env.AddAttr("{Error.activity}", data.STRING, e.ActivityName())
		}
	case *ActivityEvalError:
		env.AddAttr("{Error.activity}", data.STRING, e.TaskName())
		env.AddAttr("{Error.message}", data.STRING, err.Error())
		env.AddAttr("{Error.type}", data.STRING, e.Type())
	default:
		env.AddAttr("{Error.message}", data.STRING, err.Error())
	}

	//todo add case for *dataMapperError & *activity.Error
}

/////////////////////////////////////////
// ExecEnv - FlowContext Implementation

// Status returns the current status of the Flow Instance
func (env *ExecEnv) Status() model.FlowStatus {
	return env.status
}

func (env *ExecEnv) setStatus(status model.FlowStatus) {

	env.status = status
	env.Instance.ChangeTracker.SetStatus(status) //should this be at the exec level?
}

// FlowDefinition returns the Flow definition associated with this context
func (env *ExecEnv) FlowDefinition() *definition.Definition {
	return env.flowDef
}

// TaskInsts get the task instances
func (env *ExecEnv) TaskInsts() []model.TaskInst {
	return nil
}

/////////////////////////////////////////
// ExecEnv - data.Scope Implementation

//func (env *ExecEnv) ID() string {
//
//}
//
//// The action reference
//func (env *ExecEnv) Ref() string {
//
//}

func (env *ExecEnv) Reply(replyData map[string]*data.Attribute, err error) {
	//ac.rh.HandleResult(replyData, err)
	//only allow reply if top level flow (parent == nil)

}

// Return is used to complete the action and return the results of the execution
func (env *ExecEnv) Return(returnData map[string]*data.Attribute, err error) {
	env.forceCompletion = true
	env.returnData = returnData
	env.returnError = err
}

func (env *ExecEnv) WorkingData() data.Scope {
	return env
}

func (env *ExecEnv) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

/////////////////////////////////////////
// ExecEnv - data.Scope Implementation

// GetAttr implements data.Scope.GetAttr
func (env *ExecEnv) GetAttr(attrName string) (value *data.Attribute, exists bool) {

	if env.Attrs != nil {
		attr, found := env.Attrs[attrName]

		if found {
			return attr, true
		}
	}

	return env.flowDef.GetAttr(attrName)
}

// SetAttrValue implements api.Scope.SetAttrValue
func (env *ExecEnv) SetAttrValue(attrName string, value interface{}) error {
	if env.Attrs == nil {
		env.Attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("SetAttr - name: %s, value:%v\n", attrName, value)

	existingAttr, exists := env.GetAttr(attrName)

	//todo: optimize, use existing attr
	if exists {
		//todo handle error
		attr, _ := data.NewAttribute(attrName, existingAttr.Type(), value)
		env.Attrs[attrName] = attr
		env.Instance.ChangeTracker.AttrChange(CtUpd, attr)
		return nil
	}

	return fmt.Errorf("Attr [%s] does not exists", attrName)
}

// AddAttr add a new attribute to the instance
func (env *ExecEnv) AddAttr(attrName string, attrType data.Type, value interface{}) *data.Attribute {
	if env.Attrs == nil {
		env.Attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("AddAttr - name: %s, type: %s, value:%v\n", attrName, attrType, value)

	var attr *data.Attribute

	existingAttr, exists := env.GetAttr(attrName)

	if exists {
		attr = existingAttr
	} else {
		//todo handle error
		attr, _ = data.NewAttribute(attrName, attrType, value)
		env.Attrs[attrName] = attr
		env.Instance.ChangeTracker.AttrChange(CtAdd, attr)
	}

	return attr
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
