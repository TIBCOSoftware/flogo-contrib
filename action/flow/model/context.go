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

	////Status gets the state of the Flow instance
	Status() FlowStatus
	//
	////SetStatus sets the state of the Flow instance
	//SetStatus(status int)
}

// TaskContext is the execution context of the Task when executing
// a Task Behavior function
type TaskContext interface {

	// Status gets the state of the Task instance
	Status() TaskStatus

	// SetStatus sets the state of the Task instance
	SetStatus(status TaskStatus)

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

	// Failed marks the Activity as failed
	//Failed(err error)

	GetSetting(setting string) (value interface{}, exists bool)

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


type FlowHost interface {
	// ID returns the ID of the Action Instance
	ID() string

	// The action reference
	Ref() string

	// Get metadata of the action instance
	//InstanceMetadata() *ConfigMetadata

	// Reply is used to reply with the results of the instance execution
	Reply(replyData map[string]*data.Attribute, err error)

	// Return is used to complete the action and return the results of the execution
	Return(returnData map[string]*data.Attribute, err error)

	//todo rename, essentially the flow's attrs for now
	WorkingData() data.Scope

	//Map with action specific details/properties, flowId, etc.
	//GetDetails() map[string]string

	GetResolver() data.Resolver
}