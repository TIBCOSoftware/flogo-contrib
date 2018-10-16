package instance

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/app"
	"time"
)

type Status string

const (
	CREATED   = "Created"
	COMPLETED = "Completed"
	CANCELLED = "Cancelled"
	FAILED    = "Failed"
	SCHEDULED = "Scheduled"
	SKIPPED   = "Skipped"
	STARTED   = "Started"
	WAITING   = "Waiting"
	UNKNOWN   = "Created"
)

const FLOW_EVENT_TYPE = "flowevent"



// FlowEvent provides access to flow instance execution details
type FlowEvent struct {
	time         time.Time
	status       Status
	flowInstance *Instance
}

// Returns flow name
func (fe *FlowEvent) Name() string {
	return fe.flowInstance.Name()
}

// Returns flow ID
func (fe *FlowEvent) ID() string {
	return fe.flowInstance.ID()
}

// In case of subflow, returns parent flow name
func (fe *FlowEvent) ParentName() string {
	if fe.flowInstance.master != nil {
		return fe.flowInstance.master.Name()
	}
	return ""
}

// In case of subflow, returns parent flow ID
func (fe *FlowEvent) ParentID() string {
	if fe.flowInstance.master != nil {
		return fe.flowInstance.master.ID()
	}
	return ""
}

// Returns event time
func (fe *FlowEvent) Time() time.Time {
	return fe.time
}

// Returns application name
func (fe *FlowEvent) AppName() string {
	return app.GetName()
}

// Returns application version
func (fe *FlowEvent) AppVersion() string {
	return app.GetVersion()
}

// Returns current flow status
func (fe *FlowEvent) Status() Status {
	return fe.status
}

// Returns output data for completed flow instance
func (fe *FlowEvent) Output() map[string]interface{} {
	attrs := make(map[string]interface{})
	if fe.Status() == COMPLETED && fe.flowInstance.returnData != nil && len(fe.flowInstance.returnData) > 0 {
		for k, v := range fe.flowInstance.returnData {
			attrs[k] = v.Value()
		}
	}
	return attrs
}

// Returns error for failed flow instance
func (fe *FlowEvent) Error() error {
	return fe.flowInstance.returnError
}

// TaskEvent provides access to task instance execution details
type TaskEvent struct {
	time         time.Time
	status       Status
	taskInstance *TaskInst
}

// Returns flow name
func (te *TaskEvent) FlowName() string {
	return te.taskInstance.flowInst.Name()
}

// Returns flow ID
func (te *TaskEvent) FlowID() string {
	return te.taskInstance.flowInst.ID()
}

// Returns task name
func (te *TaskEvent) Name() string {
	return te.taskInstance.task.Name()
}

// Returns task type
func (te *TaskEvent) Type() string {
	return te.taskInstance.task.TypeID()
}

// Returns task status
func (te *TaskEvent) Status() Status {
	return te.status
}

// Returns application name
func (te *TaskEvent) AppName() string {
	return app.GetName()
}

// Returns application version
func (te *TaskEvent) AppVersion() string {
	return app.GetVersion()
}

// Returns event time
func (te *TaskEvent) Time() time.Time {
	return te.time
}

// Returns working data of current instance. e.g. key and value of current iteration for iterator task.
func (te *TaskEvent) GetWorkingData() map[string]interface{} {
	attrs := make(map[string]interface{})
	if te.taskInstance.HasWorkingData() {
		for name, value := range te.taskInstance.workingData {
			attrs[name] = value.Value()
		}
	}
	return attrs
}

// Returns activity input data
func (te *TaskEvent) ActivityInput() map[string]interface{} {
	attrs := make(map[string]interface{})
	if te.taskInstance.task.ActivityConfig().GetInputAttrs() != nil && te.taskInstance.inScope != nil {
		for name := range te.taskInstance.task.ActivityConfig().GetInputAttrs() {
			inVal, _ := te.taskInstance.inScope.GetAttr(name)
			attrs[name] = inVal
		}
	}
	return attrs
}

// Returns output data for completed activity
func (te *TaskEvent) ActivityOutput() map[string]interface{} {
	attrs := make(map[string]interface{})
	if te.Status() == COMPLETED && te.taskInstance.task.ActivityConfig().GetOutputAttrs() != nil && te.taskInstance.outScope != nil {
		for name := range te.taskInstance.task.ActivityConfig().GetOutputAttrs() {
			outVal, _ := te.taskInstance.outScope.GetAttr(name)
			attrs[name] = outVal
		}
	}
	return attrs
}

// Returns error for failed task
func (te *TaskEvent) Error() error {
	return te.taskInstance.returnError
}

func convertFlowStatus(code model.FlowStatus) Status {
	switch code {
	case model.FlowStatusNotStarted:
		return CREATED
	case model.FlowStatusActive:
		return STARTED
	case model.FlowStatusCancelled:
		return CANCELLED
	case model.FlowStatusCompleted:
		return COMPLETED
	case model.FlowStatusFailed:
		return FAILED
	}
	return UNKNOWN
}

func convertTaskStatus(code model.TaskStatus) Status {
	switch code {
	case model.TaskStatusNotStarted:
		return CREATED
	case model.TaskStatusEntered:
		return SCHEDULED
	case model.TaskStatusSkipped:
		return SKIPPED
	case model.TaskStatusReady:
		return STARTED
	case model.TaskStatusFailed:
		return FAILED
	case model.TaskStatusDone:
		return COMPLETED
	case model.TaskStatusWaiting:
		return WAITING
	}
	return UNKNOWN
}

