package instance2

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"fmt"
)


type Instance struct {

	ChangeTracker *InstanceChangeTracker

	RootInstance *Instance
	ErrorInstance *Instance
	//hostContext HostContext
	
	flowModel *model.FlowModel //todo remove reference

	status      model.FlowStatus
	flowDef   *definition.Definition

	Attrs map[string]*data.Attribute

	TaskDatas map[string]*TaskData
	LinkDatas map[int]*LinkData

	isErrorHandler  bool

	forceCompletion bool
	returnData      map[string]*data.Attribute
	returnError     error
	
	subflows map[int]*Instance
}


// ID returns the ID of the Flow Instance
func (inst *Instance) ID() string {
	return ""
}

func (inst *Instance) NewEmbeddedInstance(flow *definition.Definition) *EmbeddedInstance {

	var embeddedInst EmbeddedInstance
	embeddedInst.id = instanceID

	embeddedInst.status = model.FlowStatusNotStarted
	embeddedInst.ChangeTracker = NewInstanceChangeTracker()

	return &embeddedInst
}


// FindOrCreateTaskData finds an existing TaskData or creates ones if not found for the
// specified task the task environment
func (inst *Instance) FindOrCreateTaskData(task *definition.Task) (taskData *TaskData, created bool) {

	taskData, ok := inst.TaskDatas[task.ID()]

	created = false

	if !ok {
		taskData = NewTaskData(inst, task)
		inst.TaskDatas[task.ID()] = taskData
		inst.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

		created = true
	}

	return taskData, created
}

// FindOrCreateLinkData finds an existing LinkData or creates ones if not found for the
// specified link the task environment
func (inst *Instance) FindOrCreateLinkData(link *definition.Link) (linkData *LinkData, created bool) {

	linkData, ok := inst.LinkDatas[link.ID()]
	created = false

	if !ok {
		linkData = NewLinkData(inst, link)
		inst.LinkDatas[link.ID()] = linkData
		inst.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtAdd, ID: link.ID(), LinkData: linkData})
		created = true
	}

	return linkData, created
}

func (inst *Instance) releaseTask(task *definition.Task) {
	delete(inst.TaskDatas, task.ID())
	inst.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtDel, ID: task.ID()})
	links := task.FromLinks()

	for _, link := range links {
		delete(inst.LinkDatas, link.ID())
		inst.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtDel, ID: link.ID()})
	}
}

// GetChanges returns the Change Tracker object
func (inst *Instance) GetChanges() *InstanceChangeTracker {
	return inst.ChangeTracker
}

// ResetChanges resets an changes that were being tracked
func (inst *Instance) ResetChanges() {

	if inst.ChangeTracker != nil {
		inst.ChangeTracker.ResetChanges()
	}

	//todo: can we reuse this to avoid gc
	inst.ChangeTracker = NewInstanceChangeTracker()
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
				piAttr, exists := inst.Attrs[mdAttr.Name()]
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
	inst.ChangeTracker.SetStatus(status)
}

// FlowDefinition returns the Flow definition associated with this context
func (inst *Instance) FlowDefinition() *definition.Definition {
	return inst.flowDef
}

// TaskInsts get the task instances
func (inst *Instance) TaskInsts() []model.TaskInst {

	taskInsts := make([]model.TaskInst, 0, len(inst.TaskDatas))
	for _, value := range inst.TaskDatas {
		taskInsts = append(taskInsts, value)
	}
	return taskInsts
}

/////////////////////////////////////////
// Instance - data.Scope Implementation

// GetAttr implements data.Scope.GetAttr
func (inst *Instance) GetAttr(attrName string) (value *data.Attribute, exists bool) {

	if inst.Attrs != nil {
		attr, found := inst.Attrs[attrName]

		if found {
			return attr, true
		}
	}

	return inst.flowDef.GetAttr(attrName)
}

// SetAttrValue implements api.Scope.SetAttrValue
func (inst *Instance) SetAttrValue(attrName string, value interface{}) error {
	if inst.Attrs == nil {
		inst.Attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("SetAttr - name: %s, value:%v\n", attrName, value)

	existingAttr, exists := inst.GetAttr(attrName)

	//todo: optimize, use existing attr
	if exists {
		//todo handle error
		attr, _ := data.NewAttribute(attrName, existingAttr.Type(), value)
		inst.Attrs[attrName] = attr
		inst.ChangeTracker.AttrChange(CtUpd, attr)
		return nil
	}

	return fmt.Errorf("Attr [%s] does not exists", attrName)
}

// AddAttr add a new attribute to the instance
func (inst *Instance) AddAttr(attrName string, attrType data.Type, value interface{}) *data.Attribute {
	if inst.Attrs == nil {
		inst.Attrs = make(map[string]*data.Attribute)
	}

	logger.Debugf("AddAttr - name: %s, type: %s, value:%v\n", attrName, attrType, value)

	var attr *data.Attribute

	existingAttr, exists := inst.GetAttr(attrName)

	if exists {
		attr = existingAttr
	} else {
		//todo handle error
		attr, _ = data.NewAttribute(attrName, attrType, value)
		inst.Attrs[attrName] = attr
		inst.ChangeTracker.AttrChange(CtAdd, attr)
	}

	return attr
}

