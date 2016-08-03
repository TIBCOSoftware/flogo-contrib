package rest

import (
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-restreply")

const (
	ivCode = "code"
	ivData = "data"
)

// RESTReplyActivity is an Activity that is used to invoke a REST Operation
// inputs : {method,uri,params}
// outputs: {result}
type RESTReplyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&RESTReplyActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *RESTReplyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *RESTReplyActivity) Eval(context activity.Context) (done bool, evalError *activity.Error) {

	code := context.GetInput(ivCode).(int)
	data := context.GetInput(ivData)

	replyHandler := context.FlowDetails().ReplyHandler()

	if replyHandler != nil {
		replyHandler.Reply(code, data)
	}

	return true, nil
}
