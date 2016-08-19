package rest

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-reply")

const (
	ivCode = "code"
	ivData = "data"
)

// ReplyActivity is an Activity that is used to reply via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type ReplyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&ReplyActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *ReplyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ReplyActivity) Eval(context activity.Context) (done bool, evalError *activity.Error) {

	code := context.GetInput(ivCode).(int)
	data := context.GetInput(ivData)

	replyHandler := context.FlowDetails().ReplyHandler()

	//todo support replying with error

	if replyHandler != nil {
		replyHandler.Reply(code, data, nil)
	}

	return true, nil
}
