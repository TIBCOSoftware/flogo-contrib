package actreturn

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-flogo-return")

const (
	ivMappings = "mappings"
)

// ReturnActivity is an Activity that is used to return/return via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type ReturnActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new ReturnActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ReturnActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ReturnActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ReturnActivity) Eval(context activity.Context) (done bool, err error) {

	mappings := context.GetInput(ivMappings).([]interface{})

	log.Debugf("Mappings: %+v", mappings)

	//todo move this to a action instance level initialization, need the notion of static inputs or config
	returnMapper, err := mapper.NewBasicMapperFromAnyArray(mappings)

	if err != nil {
		return false, nil
	}

	actionCtx := context.ActionContext()

	outAttrs := actionCtx.InstanceMetadata().Output
	attrs := make([]*data.Attribute, 0, len(outAttrs))

	for _, outAttr := range outAttrs {
		attrs = append(attrs, outAttr)
	}

	//create a fixed scope using the output metadata
	outputScope := data.NewFixedScope(attrs)
	inputScope  :=  actionCtx.WorkingData() //flow data

	err = returnMapper.Apply(inputScope, outputScope)

	if err != nil {
		return false, nil
	}

	actionCtx.Return(outputScope.GetAttrs(), nil)

	return true, nil
}