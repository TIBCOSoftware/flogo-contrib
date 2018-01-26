package model

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/util"
)

// FlowModel defines the execution Model for a Flow.  It contains the
// execution behaviors for Flows and Tasks.
type FlowModel struct {
	name                string
	flowBehavior        FlowBehavior
	defaultTaskBehavior TaskBehavior

	taskBehaviors       map[int]TaskBehavior
	taskBehaviorAliases map[string]int
}

// New creates a new FlowModel from the specified Behaviors
func New(name string) *FlowModel {

	var flowModel FlowModel
	flowModel.name = name
	flowModel.taskBehaviors = make(map[int]TaskBehavior)


	return &flowModel
}

// Name returns the name of the FlowModel
func (fm *FlowModel) Name() string {
	return fm.name
}

// RegisterFlowBehavior registers the specified FlowBehavior with the Model
func (fm *FlowModel) RegisterFlowBehavior(flowBehavior FlowBehavior) {

	fm.flowBehavior = flowBehavior
}

// GetFlowBehavior returns FlowBehavior of the FlowModel
func (fm *FlowModel) GetFlowBehavior() FlowBehavior {
	return fm.flowBehavior
}

// RegisterDefaultTaskBehavior registers the default TaskBehavior for the Model
func (fm *FlowModel) RegisterDefaultTaskBehavior(taskBehavior TaskBehavior) {
	fm.defaultTaskBehavior = taskBehavior

	fm.taskBehaviors[0] = taskBehavior
	fm.taskBehaviors[1] = taskBehavior //for backwards compatibility
	util.RegisterIntAlias(fm.name + "-" + "default", 0)
	util.RegisterIntAlias(fm.name + "-" + "default", 1) //for backwards compatibility
}

// RegisterDefaultTaskBehavior registers the default TaskBehavior for the Model
func (fm *FlowModel) GetDefaultTaskBehavior() TaskBehavior {
	return fm.defaultTaskBehavior
}


// RegisterTaskBehavior registers the specified TaskBehavior with the Model
func (fm *FlowModel) RegisterTaskBehavior(id int, alias string, taskBehavior TaskBehavior) {
	if id > 1 {
		fm.taskBehaviors[id] = taskBehavior

		if alias != "" {
			util.RegisterIntAlias(fm.name + "-" + alias, id)
		}
	}

	fm.taskBehaviors[id] = taskBehavior
	if alias != "" {
		util.RegisterIntAlias(fm.name + "-" + alias, id)
	}
}

// GetTaskBehavior returns TaskBehavior with the specified ID in he FlowModel
func (fm *FlowModel) GetTaskBehavior(id int) TaskBehavior {
	if id < 2 {
		return fm.defaultTaskBehavior
	}

	return fm.taskBehaviors[id]
}

// GetTaskBehaviorByAlias returns TaskBehavior with the specified alias in he FlowModel
func (fm *FlowModel) GetTaskBehaviorByAlias(alias string) TaskBehavior {
	id, found := util.GetIntFromAlias(fm.name + "-" + alias)
	if found {
		return fm.taskBehaviors[id]
	}
	return nil
}
