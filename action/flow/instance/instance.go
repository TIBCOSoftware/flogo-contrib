package instance

import (
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/util"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

type Instance struct {
	subFlowId int

	master *IndependentInstance //needed for change tracker
	host   util.FlowHost
	//parent, should it be the taskinst?  Is this also the host context?
	//hostContext HostContext

	isErrorHandler bool
	ErrorInstance  *Instance

	status  model.FlowStatus
	flowDef *definition.Definition
	flowURI string //needed for serialization

	attrs map[string]*data.Attribute

	taskInsts map[string]*TaskInst
	linkInsts map[int]*LinkInst

	forceCompletion bool
	returnData      map[string]*data.Attribute
	returnError     error
}

func (inst *Instance) FlowURI() string {
	return inst.flowURI
}

func (inst *Instance) Name() string {
	return inst.flowDef.Name()
}

func (inst *Instance) ID() string  {
	if inst.subFlowId > 0 {
		return inst.master.id + "-" + strconv.Itoa(inst.subFlowId)
	}

	return inst.master.id
}

// FindOrCreateTaskData finds an existing TaskInst or creates ones if not found for the
// specified task the task environment
func (inst *Instance) FindOrCreateTaskData(task *definition.Task) (taskInst *TaskInst, created bool) {

	taskInst, ok := inst.taskInsts[task.ID()]

	created = false

	if !ok {
		taskInst = NewTaskInst(inst, task)
		inst.taskInsts[task.ID()] = taskInst
		inst.master.ChangeTracker.trackTaskData(inst.subFlowId, &TaskInstChange{ChgType: CtAdd, ID: task.ID(), TaskInst: taskInst})

		created = true
	}

	return taskInst, created
}

// FindOrCreateLinkData finds an existing LinkInst or creates ones if not found for the
// specified link the task environment
func (inst *Instance) FindOrCreateLinkData(link *definition.Link) (linkInst *LinkInst, created bool) {

	linkInst, ok := inst.linkInsts[link.ID()]
	created = false

	if !ok {
		linkInst = NewLinkInst(inst, link)
		inst.linkInsts[link.ID()] = linkInst
		inst.master.ChangeTracker.trackLinkData(inst.subFlowId, &LinkInstChange{ChgType: CtAdd, ID: link.ID(), LinkInst: linkInst})
		created = true
	}

	return linkInst, created
}

func (inst *Instance) releaseTask(task *definition.Task) {
	delete(inst.taskInsts, task.ID())
	inst.master.ChangeTracker.trackTaskData(inst.subFlowId, &TaskInstChange{ChgType: CtDel, ID: task.ID()})
	links := task.FromLinks()

	for _, link := range links {
		delete(inst.linkInsts, link.ID())
		inst.master.ChangeTracker.trackLinkData(inst.subFlowId, &LinkInstChange{ChgType: CtDel, ID: link.ID()})
	}
}

func (inst *Instance) appendErrorData(err error) {

	switch e := err.(type) {
	case *definition.LinkExprError:
		inst.AddAttr("{Error.type}", data.STRING, "link_expr")
		inst.AddAttr("{Error.message}", data.STRING, err.Error())
	case *activity.Error:
		inst.AddAttr("{Error.message}", data.STRING, err.Error())
		inst.AddAttr("{Error.data}", data.OBJECT, e.Data())
		inst.AddAttr("{Error.code}", data.STRING, e.Code())

		if e.ActivityName() != "" {
			inst.AddAttr("{Error.activity}", data.STRING, e.ActivityName())
		}
	case *ActivityEvalError:
		inst.AddAttr("{Error.activity}", data.STRING, e.TaskName())
		inst.AddAttr("{Error.message}", data.STRING, err.Error())
		inst.AddAttr("{Error.type}", data.STRING, e.Type())
	default:
		inst.AddAttr("{Error.message}", data.STRING, err.Error())
	}

	//todo add case for *dataMapperError & *activity.Error
}

/////////////////////////////////////////
// Instance - activity.Host Implementation

func (inst *Instance) Reply(replyData map[string]*data.Attribute, err error) {
	//ac.rh.HandleResult(replyData, err)
}

func (inst *Instance) Return(returnData map[string]*data.Attribute, err error) {
	inst.forceCompletion = true
	inst.returnData = returnData
	inst.returnError = err
}

func (inst *Instance) WorkingData() data.Scope {
	return inst
}

func (inst *Instance) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

func (inst *Instance) GetReturnData() (map[string]*data.Attribute, error) {

	if inst.returnData == nil {

		//construct returnData from instance attributes
		md := inst.flowDef.Metadata()

		if md != nil && md.Output != nil {

			inst.returnData = make(map[string]*data.Attribute)
			for _, mdAttr := range md.Output {
				piAttr, exists := inst.attrs[mdAttr.Name()]
				if exists {
					inst.returnData[piAttr.Name()] = piAttr
				}
			}
		}
	}

	return inst.returnData, inst.returnError
}

/////////////////////////////////////////
// Instance - FlowContext Implementation

// Status returns the current status of the Flow Instance
func (inst *Instance) Status() model.FlowStatus {
	return inst.status
}

func (inst *Instance) SetStatus(status model.FlowStatus) {

	inst.status = status
	inst.master.ChangeTracker.SetStatus(inst.subFlowId, status)
}

// FlowDefinition returns the Flow definition associated with this context
func (inst *Instance) FlowDefinition() *definition.Definition {
	return inst.flowDef
}

// TaskInstances get the task instances
func (inst *Instance) TaskInstances() []model.TaskInstance {

	taskInsts := make([]model.TaskInstance, 0, len(inst.taskInsts))
	for _, value := range inst.taskInsts {
		taskInsts = append(taskInsts, value)
	}
	return taskInsts
}

/////////////////////////////////////////
// Instance - data.Scope Implementation

// GetAttr implements data.Scope.GetAttr
func (inst *Instance) GetAttr(attrName string) (value *data.Attribute, exists bool) {

	if inst.attrs != nil {
		attr, found := inst.attrs[attrName]

		if found {
			return attr, true
		}
	}

	return inst.flowDef.GetAttr(attrName)
}

// SetAttrValue implements api.Scope.SetAttrValue
func (inst *Instance) SetAttrValue(attrName string, value interface{}) error {
	if inst.attrs == nil {
		inst.attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("SetAttr - name: %s, value:%v\n", attrName, value)

	existingAttr, exists := inst.GetAttr(attrName)

	//todo: optimize, use existing attr
	if exists {
		//todo handle error
		attr, _ := data.NewAttribute(attrName, existingAttr.Type(), value)
		inst.attrs[attrName] = attr
		inst.master.ChangeTracker.AttrChange(inst.subFlowId, CtUpd, attr)
		return nil
	}

	return fmt.Errorf("Attr [%s] does not exists", attrName)
}

// AddAttr add a new attribute to the instance
func (inst *Instance) AddAttr(attrName string, attrType data.Type, value interface{}) *data.Attribute {
	if inst.attrs == nil {
		inst.attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("AddAttr - name: %s, type: %s, value:%v\n", attrName, attrType, value)

	var attr *data.Attribute

	existingAttr, exists := inst.GetAttr(attrName)

	if exists {
		attr = existingAttr
	} else {
		//todo handle error
		attr, _ = data.NewAttribute(attrName, attrType, value)
		inst.attrs[attrName] = attr
		inst.master.ChangeTracker.AttrChange(inst.subFlowId, CtAdd, attr)
	}

	return attr
}
