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

type FlowEventListenerFunc func(*FlowAuditEvent)
type TaskEventListenerFunc func(*TaskAuditEvent)

var flowEventListeners []FlowEventListenerFunc
var taskEventListeners []TaskEventListenerFunc

var feLock = &sync.Mutex{}
var teLock = &sync.Mutex{}



type FlowAuditEvent struct {
	status Status
	flowName, flowId, parentFlowName, parentFlowId string
}

func (fe *FlowAuditEvent) GetFlowName() string {
	return fe.flowName
}

func (fe *FlowAuditEvent) GetFlowId() string {
	return fe.flowId
}

func (fe *FlowAuditEvent) GetParentFlowName() string {
	return fe.parentFlowName
}

func (fe *FlowAuditEvent) GetParentFlowId() string {
	return fe.parentFlowId
}

func (fe *FlowAuditEvent) GetStatus() Status {
	return fe.status
}

type TaskAuditEvent struct {
	status Status
	flowName, flowId, taskName string
}

func (te *TaskAuditEvent) GetFlowName() string {
	return te.flowName
}

func (te *TaskAuditEvent) GetFlowId() string {
	return te.flowId
}

func (te *TaskAuditEvent) GetTaskName() string {
	return te.taskName
}


func (te *TaskAuditEvent) GetStatus() Status {
	return te.status
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

func RegisterFlowEventListener(fel FlowEventListenerFunc) {
	feLock.Lock()
	defer feLock.Unlock()
	flowEventListeners = append(flowEventListeners, fel)
}

func RegisterTaskEventListener(tel TaskEventListenerFunc) {
	teLock.Lock()
	defer teLock.Unlock()
	taskEventListeners = append(taskEventListeners, tel)
}


func publishFlowEvent(fe *FlowAuditEvent) {
	for _, fel := range flowEventListeners {
		go fel(fe)
	}
}

func publishTaskEvent(te *TaskAuditEvent) {
	for _, tel := range taskEventListeners {
		go tel(te)
	}
}