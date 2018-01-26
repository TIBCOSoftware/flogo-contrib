package instance2

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"runtime/debug"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"errors"
)

type TaskData struct {
	execEnv *ExecEnv
	task    *definition.Task
	status  int


	workingData    map[string]*data.Attribute

	//state   int
	//done    bool
	//attrs   map[string]*data.Attribute
	//
	//inScope  data.Scope
	//outScope data.Scope
	//
	//changes int
	//
	//taskID string //needed for serialization
}

func NewTaskData(execEnv *ExecEnv, task *definition.Task) *TaskData {
	var taskData TaskData

	taskData.execEnv = execEnv
	taskData.task = task

	//taskData.TaskID = task.ID

	return &taskData
}


/////////////////////////////////////////
// TaskData - TaskContext Implementation

// Status implements flow.TaskContext.GetState
func (td *TaskData) Status() int {
	return td.status
}

// SetStatus implements flow.TaskContext.SetStatus
func (td *TaskData) SetStatus(status int) {
	td.status = status
	//td.execEnv.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtUpd, ID: td.task.ID(), TaskData: td})
}

func (td *TaskData) HasWorkingData() bool {
	return td.workingData != nil
}

func (td *TaskData) GetSetting(setting string) (value interface{}, exists bool) {

	value, exists = td.task.GetSetting(setting)

	if !exists {
		return nil, false
	}

	strValue, ok := value.(string)

	if ok && strValue[0] == '$' {

		v, err := definition.GetDataResolver().Resolve(strValue, td.execEnv.Instance)
		if err != nil {
			return nil, false
		}

		return v, true

	} else {
		return value, true
	}
}

func (td *TaskData) AddWorkingData(attr *data.Attribute) {

	if td.workingData == nil {
		td.workingData = make(map[string]*data.Attribute)
	}
	td.workingData[attr.Name()] = attr
}


func (td *TaskData) UpdateWorkingData(key string, value interface{}) error {

	if td.workingData == nil {
		return errors.New("working data '" + key + "' not defined")
	}

	attr, ok := td.workingData[key]

	if ok {
		attr.SetValue(value)
	} else {
		return errors.New("working data '" + key + "' not defined")
	}

	return nil
}

func (td *TaskData) GetWorkingData(key string) (*data.Attribute, bool) {
	if td.workingData == nil {
		return nil, false
	}

	v, ok := td.workingData[key]
	return v, ok
}

// Task implements model.TaskContext.Task, by returning the Task associated with this
// TaskData object
func (td *TaskData) Task() *definition.Task {
	return td.task
}

// FromInstLinks implements model.TaskContext.FromInstLinks
func (td *TaskData) FromInstLinks() []model.LinkInst {

	logger.Debugf("FromInstLinks: task=%v\n", td.Task)

	links := td.task.FromLinks()

	numLinks := len(links)

	if numLinks > 0 {
		linkCtxs := make([]model.LinkInst, numLinks)

		for i, link := range links {
			linkCtxs[i], _ = td.execEnv.FindOrCreateLinkData(link)
		}
		return linkCtxs
	}

	return nil
}

// ToInstLinks implements model.TaskContext.ToInstLinks,
func (td *TaskData) ToInstLinks() []model.LinkInst {

	logger.Debugf("ToInstLinks: task=%v\n", td.Task)

	links := td.task.ToLinks()

	numLinks := len(links)

	if numLinks > 0 {
		linkCtxs := make([]model.LinkInst, numLinks)

		for i, link := range links {
			linkCtxs[i], _ = td.execEnv.FindOrCreateLinkData(link)
		}
		return linkCtxs
	}

	return nil
}

// EvalLink implements activity.ActivityContext.EvalLink method
func (td *TaskData) EvalLink(link *definition.Link) (result bool, err error) {

	logger.Debugf("TaskContext.EvalLink: %d\n", link.ID())

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error evaluating link '%s' : %v\n", link.ID(), r)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			if err != nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	mgr := td.execEnv.flowDef.GetLinkExprManager()

	if mgr != nil {
		result, err = mgr.EvalLinkExpr(link, td.execEnv.Instance)
		return result, err
	}

	return true, nil
}

// HasActivity implements activity.ActivityContext.HasActivity method
func (td *TaskData) HasActivity() bool {
	return activity.Get(td.task.ActivityConfig().Ref()) != nil
}



// LinkData represents data associated with an instance of a Link
type LinkData struct {
	execEnv *ExecEnv
	link    *definition.Link
	status  int

	changes int

	linkID int //needed for serialization
}

// NewLinkData creates a LinkData for the specified link in the specified task
// environment
func NewLinkData(execEnv *ExecEnv, link *definition.Link) *LinkData {
	var linkData LinkData

	linkData.execEnv = execEnv
	linkData.link = link

	return &linkData
}

// Status returns the current state indicator for the LinkData
func (ld *LinkData) Status() int {
	return ld.status
}

// SetStatus sets the current state indicator for the LinkData
func (ld *LinkData) SetStatus(status int) {
	ld.status = status
	//ld.execEnv.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtUpd, ID: ld.link.ID(), LinkData: ld})
}

// Link returns the Link associated with ld context
func (ld *LinkData) Link() *definition.Link {
	return ld.link
}