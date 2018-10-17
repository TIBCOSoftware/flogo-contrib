package event

import (
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
const TASK_EVENT_TYPE = "taskevent"


// FlowEvent provides access to flow instance execution details
type FlowEvent interface {
	// Returns flow name
	Name() string
	// Returns flow ID
	ID() string
	// In case of subflow, returns parent flow name
	ParentName() string
	// In case of subflow, returns parent flow ID
	ParentID() string
	// Returns event time
	Time() time.Time
	// Returns application name
	AppName() string
	// Returns application version
	AppVersion() string
	// Returns current flow status
	Status() Status
	// Returns output data for completed flow instance
	Output() map[string]interface{}
	// Returns error for failed flow instance
	Error() error
}


// TaskEvent provides access to task instance execution details
type TaskEvent interface {

	// Returns flow name
	FlowName() string
	// Returns flow ID
	FlowID() string
	// Returns task name
	Name() string
	// Returns task type
	Type() string
	// Returns task status
	Status() Status
	// Returns application name
	AppName() string
	// Returns application version
	AppVersion() string
	// Returns event time
	Time() time.Time
	// Returns task input data
	TaskInput() map[string]interface{}
	// Returns task output data for completed task
	TaskOutput() map[string]interface{}
	// Returns error for failed task
	Error() error
}

