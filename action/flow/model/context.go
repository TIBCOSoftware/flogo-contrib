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

	// TaskInsts get the task instances
	TaskInsts() []TaskInst

	////Status gets the state of the Flow instance
	//Status() int
	//
	////SetStatus sets the state of the Flow instance
	//SetStatus(status int)
}

// TaskContext is the execution context of the Task when executing
// a Task Behavior function
type TaskContext interface {

	// Status gets the state of the Task instance
	Status() int

	// SetStatus sets the state of the Task instance
	SetStatus(status int)

	// Task returns the Task associated with this context
	Task() *definition.Task

	// FromInstLinks returns the instances of predecessor Links of the current
	// task.
	FromInstLinks() []LinkInst

	// ToInstLinks returns the instances of successor Links of the current
	// task.
	ToInstLinks() []LinkInst

	// EvalLink evaluates the specified link
	EvalLink(link *definition.Link) (bool, error)

	// HasActivity flag indicating if the task has an Activity
	HasActivity() bool

	// EvalActivity evaluates the Activity associated with the Task
	EvalActivity() (done bool, err error)

	// Failed marks the Activity as failed
	Failed(err error)

	GetSetting(setting string) (value interface{}, exists bool)

	AddWorkingData(attr *data.Attribute)

	UpdateWorkingData(key string, value interface{}) error

	GetWorkingData(key string) (*data.Attribute, bool)
}

// LinkInst is the instance of a link
type LinkInst interface {

	// Link returns the Link associated with this Link Instance
	Link() *definition.Link

	// Status gets the state of the Link instance
	Status() int

	// SetStatus sets the state of the Link instance
	SetStatus(status int)
}

type TaskInst interface {

	// Task returns the Task associated with this Task Instance
	Task() *definition.Task

	// Status gets the state of the Task instance
	Status() int
}

//type Status int
//
//const (
//	// StatusNotStarted indicates that the FlowInstance has not started
//	StatusNotStarted Status = 0
//
//	// StatusActive indicates that the FlowInstance is active
//	StatusActive Status = 100
//
//	// StatusCompleted indicates that the FlowInstance has been completed
//	StatusCompleted Status = 500
//
//	// StatusCancelled indicates that the FlowInstance has been cancelled
//	StatusCancelled Status = 600
//
//	// StatusFailed indicates that the FlowInstance has failed
//	StatusFailed Status = 700
//)
