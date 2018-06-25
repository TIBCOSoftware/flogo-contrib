package activity

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ACTION_REF = "github.com/TIBCOSoftware/flogo-contrib/action/activity"
)

type ActivityAction struct {
	act        activity.Activity
	md         *action.Metadata
	ioMetadata *data.IOMetadata
}

type actionInstance struct {
	act              activity.Activity

	inputs      map[string]*data.Attribute
	outputScope *data.FixedScope
}

type ActionData struct {
	Ref string `json:"ref"`
}

type ActionFactory struct {
}

//todo fix this
var metadata = &action.Metadata{ID: ACTION_REF, Async: false}

func init() {
	action.RegisterFactory(ACTION_REF, &ActionFactory{})
}

func (ff *ActionFactory) Init() error {
	return nil
}

func (ff *ActionFactory) New(config *action.Config) (action.Action, error) {

	activityAction := &ActivityAction{}

	var actionData ActionData
	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("failed to read activity action data '%s' error '%s'", config.Id, err.Error())
	}

	act := activity.Get(actionData.Ref)

	if act == nil {
		return nil, fmt.Errorf("failed to load activity, activity '%s' not registered", config.Id)
	}

	activityAction.act = act

	iomd := &data.IOMetadata{}
	iomd.Input = act.Metadata().Input
	iomd.Output = act.Metadata().Output
	activityAction.ioMetadata = iomd

	return activityAction, nil
}

// Metadata get the Action's metadata
func (a *ActivityAction) Metadata() *action.Metadata {
	return metadata
}

// IOMetadata get the Action's IO metadata
func (a *ActivityAction) IOMetadata() *data.IOMetadata {
	return a.ioMetadata
}

func (a *ActivityAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error executing activity [%s] : %v\n", a.act.Metadata().ID, r)

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())

			//if evalErr == nil {
			//	evalErr = NewActivityEvalError(a.task.Name(), "unhandled", fmt.Sprintf("%v", r))
			//	done = false
			//}
		}
		//if evalErr != nil {
		//	logger.Errorf("Execution failed for Activity[%s] in Flow[%s] - %s", a.task.Name(), a.flowInst.flowDef.Name(), evalErr.Error())
		//}
	}()

	ai := &actionInstance{inputs: inputs, outputScope: data.NewFixedScope(a.act.Metadata().Output)}

	_, evalErr := a.act.Eval(ai)

	if evalErr != nil {
		return nil, evalErr
	}

	return ai.outputScope.GetAttrs(), nil
}

/////////////////////////////////////////
// activity.Context Implementation

func (ai *actionInstance) ActivityHost() activity.Host {
	return ai
}

func (ai *actionInstance) Name() string {
	return ""
}

func (ai *actionInstance) GetSetting(setting string) (value interface{}, exists bool) {
	return nil, false
}

func (ai *actionInstance) GetInitValue(key string) (value interface{}, exists bool) {
	return nil, false
}

// GetInput implements activity.Context.GetInput
func (ai *actionInstance) GetInput(name string) interface{} {

	val, found := ai.inputs[name]
	if found {
		return val.Value()
	}

	return nil
}

// GetOutput implements activity.Context.GetOutput
func (ai *actionInstance) GetOutput(name string) interface{} {

	val, found := ai.outputScope.GetAttr(name)
	if found {
		return val.Value()
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (ai *actionInstance) SetOutput(name string, value interface{}) {
	ai.outputScope.SetAttrValue(name, value)
}

//Deprecated
func (ai *actionInstance) TaskName() string {
	//ignore
	return ""
}

//Deprecated
func (ai *actionInstance) FlowDetails() activity.FlowDetails {
	//ignore
	return nil
}

/////////////////////////////////////////
// activity.Host Implementation

func (ai *actionInstance) ID() string {
	//ignore
	return ""
}

func (ai *actionInstance) IOMetadata() *data.IOMetadata {
	return nil
}

func (ai *actionInstance) Reply(replyData map[string]*data.Attribute, err error) {
	// ignore
}

func (ai *actionInstance) Return(returnData map[string]*data.Attribute, err error) {
	//ignore
}

func (ai *actionInstance) WorkingData() data.Scope {
	return nil
}

func (ai *actionInstance) GetResolver() data.Resolver {
	return data.GetBasicResolver()
}
