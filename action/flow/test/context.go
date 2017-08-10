package test

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// TestActivityContext is a dummy ActivityContext to assist in testing
type TestActivityContext struct {
	details     activity.FlowDetails
	TaskNameVal string
	Attrs       map[string]*data.Attribute

	metadata *activity.Metadata
	inputs   map[string]*data.Attribute
	outputs  map[string]*data.Attribute
}

// TestFlowDetails simple FlowDetails for use in testing
type TestFlowDetails struct {
	FlowIDVal   string
	FlowNameVal string
}

// NewTestActivityContext creates a new TestActivityContext
func NewTestActivityContext(metadata *activity.Metadata) *TestActivityContext {

	fd := &TestFlowDetails{
		FlowIDVal:   "1",
		FlowNameVal: "Test Flow",
	}

	tc := &TestActivityContext{
		details:     fd,
		TaskNameVal: "Test Task",
		Attrs:       make(map[string]*data.Attribute),
		inputs:      make(map[string]*data.Attribute, len(metadata.Inputs)),
		outputs:     make(map[string]*data.Attribute, len(metadata.Outputs)),
	}

	for _, element := range metadata.Inputs {
		tc.inputs[element.Name] = data.NewAttribute(element.Name, element.Type, nil)
	}
	for _, element := range metadata.Outputs {
		tc.outputs[element.Name] = data.NewAttribute(element.Name, element.Type, nil)
	}

	return tc
}

// ID implements activity.FlowDetails.ID
func (fd *TestFlowDetails) ID() string {
	return fd.FlowIDVal
}

// Name implements activity.FlowDetails.Name
func (fd *TestFlowDetails) Name() string {
	return fd.FlowNameVal
}

// ReplyHandler implements activity.FlowDetails.ReplyHandler
func (fd *TestFlowDetails) ReplyHandler() activity.ReplyHandler {
	return nil
}

// FlowDetails implements activity.Context.FlowDetails
func (c *TestActivityContext) FlowDetails() activity.FlowDetails {
	return c.details
}

// TaskName implements activity.Context.TaskName
func (c *TestActivityContext) TaskName() string {
	return c.TaskNameVal
}

// GetAttrType implements data.Scope.GetAttrType
func (c *TestActivityContext) GetAttrType(attrName string) (attrType data.Type, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Type, true
	}

	return 0, false
}

// GetAttrValue implements data.Scope.GetAttrValue
func (c *TestActivityContext) GetAttrValue(attrName string) (value string, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Value.(string), true
	}

	return "", false
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) SetAttrValue(attrName string, value string) {

	attr, found := c.Attrs[attrName]

	if found {
		attr.Value = value
	}
}

// SetInput implements activity.Context.SetInput
func (c *TestActivityContext) SetInput(name string, value interface{}) {

	attr, found := c.inputs[name]

	if found {
		attr.Value = value
	} else {
		//error?
	}
}

// GetInput implements activity.Context.GetInput
func (c *TestActivityContext) GetInput(name string) interface{} {

	attr, found := c.inputs[name]

	if found {
		return attr.Value
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (c *TestActivityContext) SetOutput(name string, value interface{}) {

	attr, found := c.outputs[name]

	if found {
		attr.Value = value
	} else {
		//error?
	}
}

// GetOutput implements activity.Context.GetOutput
func (c *TestActivityContext) GetOutput(name string) interface{} {

	attr, found := c.outputs[name]

	if found {
		return attr.Value
	}

	return nil
}
