package subflow

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/pkg/errors"
)

// log is the default package logger
var log = logger.GetLogger("activity-flogo-subFlow")

const (
	settingFlowURI = "flowURI"
)

// SubFlowActivity is an Activity that is used to start a sub-flow, can only be used within the
// context of an flow
// settings: {flowURI}
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

	//todo move to init
	setting, set := ctx.GetSetting(settingFlowURI)

	if !set {
		return false, errors.New("flowURI not set")
	}

	flowPath := setting.(string)
	log.Debugf("Starting SubFlow: %s", flowPath)

	err = instance.StartSubFlow(ctx, flowPath)

	if err != nil {
		return false, err
	}

	return false, nil
}
