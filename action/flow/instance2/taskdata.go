package instance2

import (
	"fmt"
	"errors"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func NewTaskData(inst *Instance, task *definition.Task) *TaskData {
	var taskData TaskData

	taskData.inst = inst
	taskData.task = task

	//taskData.TaskID = task.ID

	return &taskData
}

type TaskData struct {
	inst *Instance
	task    *definition.Task
	status  model.TaskStatus

	workingData map[string]*data.Attribute

	inScope  data.Scope
	outScope data.Scope

	taskID string //needed for serialization
}

// InputScope get the InputScope of the task instance
func (td *TaskData) InputScope() data.Scope {

	if td.inScope != nil {
		return td.inScope
	}

	if len(td.task.ActivityConfig().Ref()) > 0 {

		act := activity.Get(td.task.ActivityConfig().Ref())
		td.inScope = NewFixedTaskScope(act.Metadata().Input, td.task, true)

	} else if td.task.IsScope() {

		//add flow scope
	}

	return td.inScope
}

// OutputScope get the InputScope of the task instance
func (td *TaskData) OutputScope() data.Scope {

	if td.outScope != nil {
		return td.outScope
	}

	if len(td.task.ActivityConfig().Ref()) > 0 {

		act := activity.Get(td.task.ActivityConfig().Ref())
		td.outScope = NewFixedTaskScope(act.Metadata().Output, td.task, false)

		logger.Debugf("OutputScope: %v\n", td.outScope)
	} else if td.task.IsScope() {

		//add flow scope
	}

	return td.outScope
}

/////////////////////////////////////////
// TaskData - activity.Context Implementation

func (td *TaskData) Host() activity.Host {
	return td.inst
}

// Name implements activity.Context.Name method
func (td *TaskData) Name() string {
	return td.task.Name()
}

// GetInput implements activity.Context.GetInput
func (td *TaskData) GetInput(name string) interface{} {

	val, found := td.InputScope().GetAttr(name)
	if found {
		return val.Value()
	}

	return nil
}

// GetOutput implements activity.Context.GetOutput
func (td *TaskData) GetOutput(name string) interface{} {

	val, found := td.OutputScope().GetAttr(name)
	if found {
		return val.Value()
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (td *TaskData) SetOutput(name string, value interface{}) {

	logger.Debugf("SET OUTPUT: %s = %v\n", name, value)
	td.OutputScope().SetAttrValue(name, value)
}

// TaskName implements activity.Context.TaskName method
// Deprecated
func (td *TaskData) TaskName() string {
	return td.task.Name()
}

/////////////////////////////////////////
// TaskData - TaskContext Implementation

// Status implements flow.TaskContext.GetState
func (td *TaskData) Status() model.TaskStatus {
	return td.status
}

// SetStatus implements flow.TaskContext.SetStatus
func (td *TaskData) SetStatus(status model.TaskStatus) {
	td.status = status
	td.inst.master.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtUpd, ID: td.task.ID(), TaskData: td})
}

func (td *TaskData) HasWorkingData() bool {
	return td.workingData != nil
}

func (td *TaskData) GetSetting(setting string) (value interface{}, exists bool) {

	value, exists = td.task.GetSetting(setting)

	if !exists {
		return nil, false
	}

	strValue, ok := value.(string)

	if ok && strValue[0] == '$' {

		v, err := definition.GetDataResolver().Resolve(strValue, td.inst)
		if err != nil {
			return nil, false
		}

		return v, true

	} else {
		return value, true
	}
}

func (td *TaskData) AddWorkingData(attr *data.Attribute) {

	if td.workingData == nil {
		td.workingData = make(map[string]*data.Attribute)
	}
	td.workingData[attr.Name()] = attr
}

func (td *TaskData) UpdateWorkingData(key string, value interface{}) error {

	if td.workingData == nil {
		return errors.New("working data '" + key + "' not defined")
	}

	attr, ok := td.workingData[key]

	if ok {
		attr.SetValue(value)
	} else {
		return errors.New("working data '" + key + "' not defined")
	}

	return nil
}

func (td *TaskData) GetWorkingData(key string) (*data.Attribute, bool) {
	if td.workingData == nil {
		return nil, false
	}

	v, ok := td.workingData[key]
	return v, ok
}

// Task implements model.TaskContext.Task, by returning the Task associated with this
// TaskData object
func (td *TaskData) Task() *definition.Task {
	return td.task
}

// FromInstLinks implements model.TaskContext.FromInstLinks
func (td *TaskData) FromInstLinks() []model.LinkInst {

	logger.Debugf("FromInstLinks: task=%v\n", td.Task)

	links := td.task.FromLinks()

	numLinks := len(links)

	if numLinks > 0 {
		linkCtxs := make([]model.LinkInst, numLinks)

		for i, link := range links {
			linkCtxs[i], _ = td.inst.FindOrCreateLinkData(link)
		}
		return linkCtxs
	}

	return nil
}

// ToInstLinks implements model.TaskContext.ToInstLinks,
func (td *TaskData) ToInstLinks() []model.LinkInst {

	logger.Debugf("ToInstLinks: task=%v\n", td.Task)

	links := td.task.ToLinks()

	numLinks := len(links)

	if numLinks > 0 {
		linkCtxs := make([]model.LinkInst, numLinks)

		for i, link := range links {
			linkCtxs[i], _ = td.inst.FindOrCreateLinkData(link)
		}
		return linkCtxs
	}

	return nil
}

// EvalLink implements activity.ActivityContext.EvalLink method
func (td *TaskData) EvalLink(link *definition.Link) (result bool, err error) {

	logger.Debugf("TaskContext.EvalLink: %d\n", link.ID())

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error evaluating link '%s' : %v\n", link.ID(), r)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			if err != nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	mgr := td.inst.flowDef.GetLinkExprManager()

	if mgr != nil {
		result, err = mgr.EvalLinkExpr(link, td.inst)
		return result, err
	}

	return true, nil
}

// HasActivity implements activity.ActivityContext.HasActivity method
func (td *TaskData) HasActivity() bool {
	return activity.Get(td.task.ActivityConfig().Ref()) != nil
}

// EvalActivity implements activity.ActivityContext.EvalActivity method
func (td *TaskData) EvalActivity() (done bool, evalErr error) {

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error executing activity '%s'[%s] : %v\n", td.task.Name(), td.task.ActivityConfig().Ref(), r)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			if evalErr == nil {
				evalErr = NewActivityEvalError(td.task.Name(), "unhandled", fmt.Sprintf("%v", r))
				done = false
			}
		}
		if evalErr != nil {
			logger.Errorf("Execution failed for Activity[%s] in Flow[%s] - %s", td.task.Name(), td.inst.flowDef.Name(), evalErr.Error())
		}
	}()

	eval := true

	if td.task.ActivityConfig().InputMapper() != nil {

		err := applyInputMapper(td)

		if err != nil {

			evalErr = NewActivityEvalError(td.task.Name(), "mapper", err.Error())
			return false, evalErr
		}

		eval = applyInputInterceptor(td)
	}

	if eval {

		act := activity.Get(td.task.ActivityConfig().Ref())
		done, evalErr = act.Eval(td)

		if evalErr != nil {
			e, ok := evalErr.(*activity.Error)
			if ok {
				e.SetActivityName(td.task.Name())
			}

			return false, evalErr
		}
	} else {
		done = true
	}

	if done {

		if td.task.ActivityConfig().OutputMapper() != nil {
			applyOutputInterceptor(td)

			appliedMapper, err := applyOutputMapper(td)

			if err != nil {
				evalErr = NewActivityEvalError(td.task.Name(), "mapper", err.Error())
				return done, evalErr
			}

			if !appliedMapper && !td.task.IsScope() {

				logger.Debug("Mapper not applied")
			}
		}
	}

	return done, nil
}

//// Failed marks the Activity as failed
//func (td *TaskData) Failed(err error) {
//
//	errorMsgAttr := "[A" + td.task.ID() + "._errorMsg]"
//	td.inst.AddAttr(errorMsgAttr, data.STRING, err.Error())
//	errorMsgAttr2 := "[activity." + td.task.ID() + "._errorMsg]"
//	td.inst.AddAttr(errorMsgAttr2, data.STRING, err.Error())
//}
