package subflow

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
)

// log is the default package logger
var log = logger.GetLogger("activity-flogo-subFlow")

const (
	settingFlowPath = "flowPath"
)

// SubFlowActivity is an Activity that is used to start a sub-flow, can only be used within the
// context of an flow
// settings: {flowPath}
// input : {sub-flow's input}
// output: {sub-flow's output}
type SubFlowActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new SubFlowActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &SubFlowActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *SubFlowActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *SubFlowActivity) Eval(ctx activity.Context) (done bool, err error) {

	flowPath := ctx.GetSetting(settingFlowPath).(string)
	log.Debugf("Starting SubFlow: %s", flowPath)


	//r := ctx.ActivityHost().GetResolver()
	//r.Resolve()

	instance.StartSubFlow(ctx, flowPath)

	return false, nil
}
