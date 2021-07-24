package model

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// FlowContext is the execution context of the Flow when executing
// a Flow Behavior function
type FlowContext interface {
	// FlowDefinition returns the Flow definition associated with this context
	FlowDefinition() *definition.Definition

	// TaskInstances get the task instances
	TaskInstances() []TaskInstance

	// Status gets the state of the Flow instance
	Status() FlowStatus
}

// TaskContext is the execution context of the Task when executing
// a Task Behavior function
type TaskContext interface {

	// Status gets the state of the Task instance
	Status() TaskStatus

	// SetStatus sets the state of the Task instance
	SetStatus(status TaskStatus)

	// SetStatus sets the state of the Task instance
	SetStatusWithError(status TaskStatus, err error)

	// Task returns the Task associated with this context
	Task() *definition.Task

	// GetFromLinkInstances returns the instances of predecessor Links of the current task.
	GetFromLinkInstances() []LinkInstance

	// GetToLinkInstances returns the instances of successor Links of the current task.
	GetToLinkInstances() []LinkInstance

	// EvalLink evaluates the specified link
	EvalLink(link *definition.Link) (bool, error)

	// EvalActivity evaluates the Activity associated with the Task
	EvalActivity() (done bool, err error)

	// PostActivity does post evaluation of the Activity associated with the Task
	PostEvalActivity() (done bool, err error)

	Resolve(toResolve string) (value interface{}, err error)

	//todo  move to a mutable scope

	AddWorkingData(attr *data.Attribute)

	UpdateWorkingData(key string, value interface{}) error

	GetWorkingData(key string) (*data.Attribute, bool)
}

// LinkInstance is the instance of a link
type LinkInstance interface {

	// Link returns the Link associated with this Link Instance
	Link() *definition.Link

	// Status gets the state of the Link instance
	Status() LinkStatus

	// SetStatus sets the state of the Link instance
	SetStatus(status LinkStatus)
}

type TaskInstance interface {

	// Task returns the Task associated with this Task Instance
	Task() *definition.Task

	// Status gets the state of the Task instance
	Status() TaskStatus
}
