package instance2

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
)

// ChgType denotes the type of change for an object in an instance
type ChgType int

const (
	// CtAdd denotes an addition
	CtAdd ChgType = 1
	// CtUpd denotes an update
	CtUpd ChgType = 2
	// CtDel denotes an deletion
	CtDel ChgType = 3
)

// WorkItemQueueChange represents a change in the WorkItem Queue
type WorkItemQueueChange struct {
	ChgType  ChgType
	ID       int
	WorkItem *WorkItem
}

// TaskDataChange represents a change to a TaskData
type TaskDataChange struct {
	SubFlowID int
	ChgType   ChgType
	ID        string
	TaskData  *TaskData
}

// LinkDataChange represents a change to a LinkData
type LinkDataChange struct {
	SubFlowID int
	ChgType   ChgType
	ID        int
	LinkData  *LinkData
}

// InstanceChange represents a change to the instance
type InstanceChange struct {
	SubFlowID   int
	State       int
	Status      model.FlowStatus
	AttrChanges []*AttributeChange
	SubFlowChg  *SubFlowChange
}

// InstanceChange represents a change to the instance
type SubFlowChange struct {
	SubFlowID int
	TaskID    string
	ChgType   ChgType
}

// AttributeChange represents a change to an Attribute
type AttributeChange struct {
	SubFlowID int
	ChgType   ChgType
	Attribute *data.Attribute
}

// InstanceChangeTracker is used to track all changes to an instance
type InstanceChangeTracker struct {
	wiqChanges map[int]*WorkItemQueueChange

	tdChanges  map[string]*TaskDataChange
	ldChanges  map[int]*LinkDataChange

	instChanges map[int]*InstanceChange //at most 2
}

// NewInstanceChangeTracker creates an InstanceChangeTracker
func NewInstanceChangeTracker() *InstanceChangeTracker {

	var tracker InstanceChangeTracker
	tracker.instChanges =make(map[int]*InstanceChange)
	return &tracker
}

// SetStatus is called to track a state change on an instance
func (ict *InstanceChangeTracker) SetState(subFlowId int, state int) {

	ict.instChanges[subFlowId].State = state
}

// SetStatus is called to track a status change on an instance
func (ict *InstanceChangeTracker) SetStatus(subFlowId int, status model.FlowStatus) {

	ict.instChanges[subFlowId].Status = status
}

// AttrChange is called to track a status change of an Attribute
func (ict *InstanceChangeTracker) AttrChange(subFlowId int, chgType ChgType, attribute *data.Attribute) {

	var attrChange AttributeChange
	attrChange.ChgType = chgType

	attrChange.Attribute = attribute
	ict.instChanges[subFlowId].AttrChanges = append(ict.instChanges[subFlowId].AttrChanges, &attrChange)
}

// AttrChange is called to track a status change of an Attribute
func (ict *InstanceChangeTracker) SubFlowChange(parentFlowId int, chgType ChgType, subFlowId int, taskID string) {

	var change SubFlowChange
	change.ChgType = chgType
	change.SubFlowID = subFlowId
	change.TaskID = taskID

	ict.instChanges[parentFlowId].SubFlowChg = &change
}

// trackWorkItem records a WorkItem Queue change
func (ict *InstanceChangeTracker) trackWorkItem(wiChange *WorkItemQueueChange) {

	if ict.wiqChanges == nil {
		ict.wiqChanges = make(map[int]*WorkItemQueueChange)
	}
	ict.wiqChanges[wiChange.ID] = wiChange
}

// trackTaskData records a TaskData change
func (ict *InstanceChangeTracker) trackTaskData(tdChange *TaskDataChange) {

	if ict.tdChanges == nil {
		ict.tdChanges = make(map[string]*TaskDataChange)
	}

	ict.tdChanges[tdChange.ID] = tdChange
}

// trackLinkData records a LinkData change
func (ict *InstanceChangeTracker) trackLinkData(ldChange *LinkDataChange) {

	if ict.ldChanges == nil {
		ict.ldChanges = make(map[int]*LinkDataChange)
	}
	ict.ldChanges[ldChange.ID] = ldChange
}

// ResetChanges is used to reset any tracking data stored on instance objects
func (ict *InstanceChangeTracker) ResetChanges() {

	// reset TaskData objects
	if ict.tdChanges != nil {
		for _, v := range ict.tdChanges {
			if v.TaskData != nil {
				//v.TaskData.ResetChanges()
			}
		}
	}

	// reset LinkData objects
	if ict.ldChanges != nil {
		for _, v := range ict.ldChanges {
			if v.LinkData != nil {
				//v.LinkData.ResetChanges()
			}
		}
	}
}
