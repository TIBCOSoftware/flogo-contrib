package instance2

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"runtime/debug"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Execution Environment

func (env *Instance) NewExecEnv(parentEnv *ExecEnv, flow *definition.Definition) *ExecEnv {
	var execEnv *ExecEnv
	execEnv.ID = env.execEnvId
	execEnv.Instance = env
	execEnv.flowDef = flow

	env.execEnvId++
	////execEnv.Task = flow.RootTask()
	////execEnv.taskID = flow.RootTask().ID()
	//execEnv.Instance = &env
	//execEnv.TaskDatas = make(map[string]*TaskData)
	//execEnv.LinkDatas = make(map[int]*LinkData)

	return execEnv
}

// ExecEnv is a structure that describes the execution environment for a flow
type ExecEnv struct {
	ParentEnv *ExecEnv
	ID        int
	flowDef   *definition.Definition
	Instance  *Instance

	//Task      *definition.Task

	//Attrs         map[string]*data.Attribute

	TaskDatas map[string]*TaskData
	LinkDatas map[int]*LinkData

	//returnData      map[string]*data.Attribute
	//returnError     error

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
		//env.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

		created = true
	}

	return taskData, created
}

// NewTaskData creates a new TaskData object
func (env *ExecEnv) NewTaskData(task *definition.Task) *TaskData {

	taskData := NewTaskData(env, task)
	env.TaskDatas[task.ID()] = taskData
	//env.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

	return taskData
}

// FindOrCreateLinkData finds an existing LinkData or creates ones if not found for the
// specified link the task environment
func (env *ExecEnv) FindOrCreateLinkData(link *definition.Link) (linkData *LinkData, created bool) {

	linkData, ok := env.LinkDatas[link.ID()]
	created = false

	if !ok {
		linkData = NewLinkData(env, link)
		env.LinkDatas[link.ID()] = linkData
		//env.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtAdd, ID: link.ID(), LinkData: linkData})
		created = true
	}

	return linkData, created
}

func (env *ExecEnv) releaseTask(task *definition.Task) {
	delete(env.TaskDatas, task.ID())
	//env.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtDel, ID: task.ID()})
	links := task.FromLinks()

	for _, link := range links {
		delete(env.LinkDatas, link.ID())
		//env.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtDel, ID: link.ID()})
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
		env.handleTaskDone(behavior, taskData, doneCode)
	} else if evalResult == model.EVAL_REPEAT {
		//iterate or retry
		env.Instance.scheduleEval(taskData, 0)
	}
}

// handleTaskDone handles the completion of a task in the Flow Instance
func (env *ExecEnv) handleTaskDone(taskBehavior model.TaskBehavior, taskData *TaskData, doneCode int) {

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

	env.releaseTask(task)
}

// handleTaskError handles the completion of a task in the Flow Instance
func (env *ExecEnv) handleTaskError(taskBehavior model.TaskBehavior, taskData *TaskData, err error) {

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
