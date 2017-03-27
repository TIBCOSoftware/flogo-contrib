package error

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-error")

const (
	ivMessage = "message"
	ivData    = "data"
)

// ErrorActivity is an Activity that used to cause an explicit error in the flow
// inputs : {message,data}
// outputs: node
type ErrorActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&ErrorActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *ErrorActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ErrorActivity) Eval(context activity.Context) (done bool, err error) {

	mesg := context.GetInput(ivMessage).(string)
	data := context.GetInput(ivData)

	log.Debugf("Message :'%s', Data: '%+v'", mesg, data)

	return false, activity.NewErrorWithData(mesg, data)
}
