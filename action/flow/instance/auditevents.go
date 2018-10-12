package instance

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"sync"
)

type Status int

const (
	Created   Status = iota
	Completed
	Cancelled
	Failed
	Scheduled
	Skipped
	Started
	Waiting
	Unknown
)

type FlowEventListenerFunc func(*FlowEventContext)
type TaskEventListenerFunc func(*TaskEventContext)

var flowEventListeners []FlowEventListenerFunc
var taskEventListeners []TaskEventListenerFunc

var feLock = &sync.Mutex{}
var teLock = &sync.Mutex{}

// FlowEventContext provides access to flow instance execution details
type FlowEventContext struct {
	flowInstance *Instance
}

// Returns flow name
func (fe *FlowEventContext) Name() string {
	return fe.flowInstance.Name()
}

// Returns flow ID
func (fe *FlowEventContext) ID() string {
	return fe.flowInstance.ID()
}

// In case of subflow, returns parent flow name
func (fe *FlowEventContext) ParentName() string {
	if fe.flowInstance.master != nil {
		return fe.flowInstance.master.Name()
	}
	return ""
}

// In case of subflow, returns parent flow ID
func (fe *FlowEventContext) ParentID() string {
	if fe.flowInstance.master != nil {
		return fe.flowInstance.master.ID()
	}
	return ""
}

// Returns current flow status
func (fe *FlowEventContext) Status() Status {
	return convertFlowStatus(fe.flowInstance.Status())
}

// Returns flow input data
func (fe *FlowEventContext) Input() map[string]interface{} {
	attrs := make(map[string]interface{})
	if fe.flowInstance.attrs != nil && len(fe.flowInstance.attrs) > 0 {
		for k, v := range fe.flowInstance.attrs {
			attrs[k] = v.Value()
		}
	}
	return attrs
}

// Returns flow output data
func (fe *FlowEventContext) Output() map[string]interface{} {
	attrs := make(map[string]interface{})
	if fe.flowInstance.returnData != nil && len(fe.flowInstance.returnData) > 0 {
		for k, v := range fe.flowInstance.returnData {
			attrs[k] = v.Value()
		}
	}
	return attrs
}

// TaskEventContext provides access to task instance execution details
type TaskEventContext struct {
	ti *TaskInst
}

// Returns flow name
func (te *TaskEventContext) FlowName() string {
	return te.ti.flowInst.Name()
}

// Returns flow ID
func (te *TaskEventContext) FlowID() string {
	return te.ti.flowInst.ID()
}

// Returns task name
func (te *TaskEventContext) Name() string {
	return te.ti.task.Name()
}

// Returns task status
func (te *TaskEventContext) Status() Status {
	return convertTaskStatus(te.ti.status)
}

//TODO
// Returns task input data
func (te *TaskEventContext) Input() map[string]interface{} {
	attrs := make(map[string]interface{})
	return attrs
}

//TODO
// Returns task output data
func (te *TaskEventContext) Output() map[string]interface{} {
	attrs := make(map[string]interface{})
	return attrs
}

func convertFlowStatus(code model.FlowStatus) Status {
	switch code {
	case model.FlowStatusNotStarted:
		return Created
	case model.FlowStatusActive:
		return Started
	case model.FlowStatusCancelled:
		return Cancelled
	case model.FlowStatusCompleted:
		return Completed
	case model.FlowStatusFailed:
		return Failed
	}
	return Unknown
}

func convertTaskStatus(code model.TaskStatus) Status {
	switch code {
	case model.TaskStatusNotStarted:
		return Created
	case model.TaskStatusEntered:
		return Scheduled
	case model.TaskStatusSkipped:
		return Skipped
	case model.TaskStatusReady:
		return Started
	case model.TaskStatusFailed:
		return Failed
	case model.TaskStatusDone:
		return Completed
	case model.TaskStatusWaiting:
		return Waiting
	}
	return Unknown
}

// Registers listener for flow events
func RegisterFlowEventListener(fel FlowEventListenerFunc) {
	feLock.Lock()
	defer feLock.Unlock()
	flowEventListeners = append(flowEventListeners, fel)
}

// Registers listener for task events
func RegisterTaskEventListener(tel TaskEventListenerFunc) {
	teLock.Lock()
	defer teLock.Unlock()
	taskEventListeners = append(taskEventListeners, tel)
}

func publishFlowEvent(fe *FlowEventContext) {
	for _, fel := range flowEventListeners {
		go fel(fe)
	}
}

func publishTaskEvent(te *TaskEventContext) {
	for _, tel := range taskEventListeners {
		go tel(te)
	}
}
