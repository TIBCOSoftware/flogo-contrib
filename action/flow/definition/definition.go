package definition

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// Definition is the object that describes the definition of
// a flow.  It contains its data (attributes) and
// structure (tasks & links).
type Definition struct {
	name          string
	modelID       string
	explicitReply bool

	attrs map[string]*data.Attribute

	links       map[int]*Link
	tasks       map[string]*Task

	//ehLinks       map[int]*Link
	//ehTasks       map[string]*Task

	linkExprMgr LinkExprManager

	errorHandlerFlow *Definition
}

// Name returns the name of the definition
func (pd *Definition) Name() string {
	return pd.name
}

// ModelID returns the ID of the model the definition uses
func (pd *Definition) ModelID() string {
	return pd.modelID
}

//// RootTask returns the root task of the definition
//func (pd *Definition) RootTask() *TaskOld {
//	return pd.rootTask
//}

// GetTask returns the task with the specified ID
func (pd *Definition) GetTask(taskID string) *Task {
	task := pd.tasks[taskID]
	return task
}

func (pd *Definition) GetTasks() []*Task {

	tasks := make([]*Task,0,len(pd.tasks))
	for _, task := range pd.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetLink returns the link with the specified ID
func (pd *Definition) GetLink(linkID int) *Link {
	task := pd.links[linkID]
	return task
}


func (pd *Definition) GetLinks() []*Link {
	links := make([]*Link,0,len(pd.links))
	for _, link := range pd.links {
		links = append(links, link)
	}
	return links
}


func (pd *Definition) ExplicitReply() bool {
	return pd.explicitReply
}

func (pd *Definition) GetErrorHandlerFlow() *Definition {
	return pd.errorHandlerFlow
}

//// ErrorHandler returns the error handler task of the definition
//func (pd *Definition) ErrorHandlerTask() *TaskOld {
//	return pd.ehTask
//}
//
//// GetAttr gets the specified attribute
func (pd *Definition) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if pd.attrs != nil {
		attr, found := pd.attrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetTask returns the task with the specified ID
func (pd *Definition) Tasks() []*Task {

	tasks := make([]*Task, len(pd.tasks))
	for _, task := range pd.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}



// SetLinkExprManager sets the LinkOld Expression Manager for the definition
func (pd *Definition) SetLinkExprManager(mgr LinkExprManager) {
	// todo revisit
	pd.linkExprMgr = mgr
}

// GetLinkExprManager gets the LinkOld Expression Manager for the definition
func (pd *Definition) GetLinkExprManager() LinkExprManager {
	return pd.linkExprMgr
}

type ActivityConfig struct {

	Activity    activity.Activity
	inputAttrs  map[string]*data.Attribute
	outputAttrs map[string]*data.Attribute

	inputMapper  data.Mapper
	outputMapper data.Mapper
}


// GetAttr gets the specified input attribute
func (task *ActivityConfig) GetInputAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.inputAttrs != nil {
		attr, found := task.inputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetOutputAttr gets the specified output attribute
func (task *ActivityConfig) GetOutputAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.outputAttrs != nil {
		attr, found := task.outputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// InputMapper returns the InputMapper of the task
func (task *ActivityConfig) InputMapper() data.Mapper {
	return task.inputMapper
}

// OutputMapper returns the OutputMapper of the task
func (task *ActivityConfig) OutputMapper() data.Mapper {
	return task.outputMapper
}

func (task *ActivityConfig) Ref() string {
	return task.Activity.Metadata().ID
}

// Task is the object that describes the definition of
// a task.  It contains its data (attributes) and its
// nested structure (child tasks & child links).
type Task struct {
	id           string
	typeID       int
	name         string
	activityCfg  *ActivityConfig

	isScope      bool

	definition *Definition
	//parent     *Task

	settings    map[string]interface{}
	inputAttrs  map[string]*data.Attribute
	outputAttrs map[string]*data.Attribute

	inputMapper  data.Mapper
	outputMapper data.Mapper

	toLinks   []*Link
	fromLinks []*Link
}

// ID gets the id of the task
func (task *Task) ID() string {
	return task.id
}

// Name gets the name of the task
func (task *Task) Name() string {
	return task.name
}

// TypeID gets the id of the task type
func (task *Task) TypeID() int {
	return task.typeID
}

//// ActivityRef gets the activity ref
//func (task *Task) ActivityRef() string {
//	return task.activityRef
//}

// Parent gets the parent task of the task
//func (task *Task) Parent() *Task {
//	return task.parent
//}

// ChildTasks gets the child tasks of the task
//func (task *Task) ChildTasks() []*Task {
//	return task.tasks
//}

// ChildLinks gets the child tasks of the task
//func (task *Task) ChildLinks() []*Link {
//	return task.links
//}


// GetAttr gets the specified attribute
// DEPRECATED
func (task *Task) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.inputAttrs != nil {
		attr, found := task.inputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetAttr gets the specified input attribute
func (task *Task) GetInputAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.inputAttrs != nil {
		attr, found := task.inputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetOutputAttr gets the specified output attribute
func (task *Task) GetOutputAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.outputAttrs != nil {
		attr, found := task.outputAttrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

func (task *Task) ActivityConfig() *ActivityConfig {
	return task.activityCfg
}

func (task *Task) GetSetting(attrName string) (value interface{}, exists bool) {
	value, exists = task.settings[attrName]
	return value,exists
}

// ToLinks returns the predecessor links of the task
func (task *Task) ToLinks() []*Link {
	return task.toLinks
}

// FromLinks returns the successor links of the task
func (task *Task) FromLinks() []*Link {
	return task.fromLinks
}

func (task *Task) String() string {
	return fmt.Sprintf("Task[%d]:'%s'", task.id, task.name)
}

// IsScope returns flag indicating if the Task is a scope task (a container of attributes)
func (task *Task) IsScope() bool {
	return task.isScope
}

//////////////////////////////////////////////////////////////////////////////
//// TaskOld
//
//// TaskOld is the object that describes the definition of
//// a task.  It contains its data (attributes) and its
//// nested structure (child tasks & child links).
//type TaskOld struct {
//	id           string
//	typeID       int
//	activityType string
//	activityRef  string
//	name         string
//	tasks        []*TaskOld
//	links        []*LinkOld
//	isScope      bool
//
//	definition *Definition
//	parent     *TaskOld
//
//	inputAttrs  map[string]*data.Attribute
//	outputAttrs map[string]*data.Attribute
//
//	inputMapper  data.Mapper
//	outputMapper data.Mapper
//
//	toLinks   []*LinkOld
//	fromLinks []*LinkOld
//}
//
//// ID gets the id of the task
//func (task *TaskOld) ID() string {
//	return task.id
//}
//
//// Name gets the name of the task
//func (task *TaskOld) Name() string {
//	return task.name
//}
//
//// TypeID gets the id of the task type
//func (task *TaskOld) TypeID() int {
//	return task.typeID
//}
//
//// ActivityType gets the activity type
//func (task *TaskOld) ActivityType() string {
//	return task.activityType
//}
//
//// ActivityRef gets the activity ref
//func (task *TaskOld) ActivityRef() string {
//	return task.activityRef
//}
//
//// Parent gets the parent task of the task
//func (task *TaskOld) Parent() *TaskOld {
//	return task.parent
//}
//
//// ChildTasks gets the child tasks of the task
//func (task *TaskOld) ChildTasks() []*TaskOld {
//	return task.tasks
//}
//
//// ChildLinks gets the child tasks of the task
//func (task *TaskOld) ChildLinks() []*LinkOld {
//	return task.links
//}
//
//// GetAttr gets the specified attribute
//// DEPRECATED
//func (task *TaskOld) GetAttr(attrName string) (attr *data.Attribute, exists bool) {
//
//	if task.inputAttrs != nil {
//		attr, found := task.inputAttrs[attrName]
//		if found {
//			return attr, true
//		}
//	}
//
//	return nil, false
//}
//
//// GetAttr gets the specified input attribute
//func (task *TaskOld) GetInputAttr(attrName string) (attr *data.Attribute, exists bool) {
//
//	if task.inputAttrs != nil {
//		attr, found := task.inputAttrs[attrName]
//		if found {
//			return attr, true
//		}
//	}
//
//	return nil, false
//}
//
//// GetOutputAttr gets the specified output attribute
//func (task *TaskOld) GetOutputAttr(attrName string) (attr *data.Attribute, exists bool) {
//
//	if task.outputAttrs != nil {
//		attr, found := task.outputAttrs[attrName]
//		if found {
//			return attr, true
//		}
//	}
//
//	return nil, false
//}
//
//// ToLinks returns the predecessor links of the task
//func (task *TaskOld) ToLinks() []*LinkOld {
//	return task.toLinks
//}
//
//// FromLinks returns the successor links of the task
//func (task *TaskOld) FromLinks() []*LinkOld {
//	return task.fromLinks
//}
//
//// InputMapper returns the InputMapper of the task
//func (task *TaskOld) InputMapper() data.Mapper {
//	return task.inputMapper
//}
//
//// OutputMapper returns the OutputMapper of the task
//func (task *TaskOld) OutputMapper() data.Mapper {
//	return task.outputMapper
//}
//
//func (task *TaskOld) String() string {
//	return fmt.Sprintf("TaskOld[%d]:'%s'", task.id, task.name)
//}
//
//// IsScope returns flag indicating if the TaskOld is a scope task (a container of attributes)
//func (task *TaskOld) IsScope() bool {
//	return task.isScope
//}

////////////////////////////////////////////////////////////////////////////
// Link

// LinkType is an enum for possible Link Types
type LinkType int

const (
	// LtDependency denotes an normal dependency link
	LtDependency LinkType = 0

	// LtExpression denotes a link with an expression
	LtExpression LinkType = 1 //expr language on the model or def?

	// LtLabel denotes 'labelled' link
	LtLabel LinkType = 2

	// LtError denotes an error link
	LtError LinkType = 3
)

// LinkOld is the object that describes the definition of
// a link.
type Link struct {
	id       int
	name     string
	fromTask *Task
	toTask   *Task
	linkType LinkType
	value    string //expression or label

	definition *Definition
}

// ID gets the id of the link
func (link *Link) ID() int {
	return link.id
}

// Type gets the link type
func (link *Link) Type() LinkType {
	return link.linkType
}

// Value gets the "value" of the link
func (link *Link) Value() string {
	return link.value
}

// FromTask returns the task the link is coming from
func (link *Link) FromTask() *Task {
	return link.fromTask
}

// ToTask returns the task the link is going to
func (link *Link) ToTask() *Task {
	return link.toTask
}

func (link *Link) String() string {
	return fmt.Sprintf("Link[%d]:'%s' - [from:%d, to:%d]", link.id, link.name, link.fromTask.id, link.toTask.id)
}

//// LinkOld is the object that describes the definition of
//// a link.
//type LinkOld struct {
//	id       int
//	name     string
//	fromTask *TaskOld
//	toTask   *TaskOld
//	linkType LinkType
//	value    string //expression or label
//
//	definition *Definition
//	parent     *TaskOld
//}
//
//// ID gets the id of the link
//func (link *LinkOld) ID() int {
//	return link.id
//}
//
//// Type gets the link type
//func (link *LinkOld) Type() LinkType {
//	return link.linkType
//}
//
//// Value gets the "value" of the link
//func (link *LinkOld) Value() string {
//	return link.value
//}
//
//// FromTask returns the task the link is coming from
//func (link *LinkOld) FromTask() *TaskOld {
//	return link.fromTask
//}
//
//// ToTask returns the task the link is going to
//func (link *LinkOld) ToTask() *TaskOld {
//	return link.toTask
//}
//
//func (link *LinkOld) String() string {
//	return fmt.Sprintf("LinkOld[%d]:'%s' - [from:%d, to:%d]", link.id, link.name, link.fromTask.id, link.toTask.id)
//}
