package reply

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-reply")

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

// NewActivity creates a new ReplyActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ReplyActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ReplyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ReplyActivity) Eval(context activity.Context) (done bool, err error) {

	code := context.GetInput(ivCode).(int)
	data := context.GetInput(ivData)

	log.Debugf("Code :'%d', Data: '%+v'", code, data)

	replyHandler := context.FlowDetails().ReplyHandler()

	//todo support replying with error

	if replyHandler != nil {

		//todo fix to support new ReplyWithData (had to keep old Reply for backwards compatibility)
		replyHandler.Reply(code, data, nil)
	}

	return true, nil
}
