package log

import (
	"fmt"
	"os"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
)

// activityLog is the default logger for the Log Activity
var activityLog = logging.MustGetLogger("activity-tibco-log")

// format is the log format for the Activity log
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)

var backend = logging.NewLogBackend(os.Stdout, "", 0)
var backendFormatter = logging.NewBackendFormatter(backend, format)
var backendLeveled = logging.AddModuleLevel(backendFormatter)

const (
	ivMessage   = "message"
	ivFlowInfo  = "flowInfo"
	ivAddToFlow = "addToFlow"

	ovMessage = "message"
)

func init() {
	backendLeveled.SetLevel(logging.DEBUG, "")
	activityLog.SetBackend(backendLeveled)
}

// LogActivity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type LogActivity struct {
	metadata *activity.Metadata
}

func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&LogActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *LogActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *LogActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	message, _ := context.GetInput(ivMessage).(string)

	flowInfo, _ := toBool(context.GetInput(ivFlowInfo))
	addToFlow, _ := toBool(context.GetInput(ivAddToFlow))

	msg := message

	if flowInfo {

		msg = fmt.Sprintf("'%s' - FlowInstanceID [%s], Flow [%s], Task [%s]", msg,
			context.FlowDetails().ID(), context.FlowDetails().Name(), context.TaskName())
	}

	activityLog.Info(msg)

	if addToFlow {
		context.SetOutput(ovMessage, msg)
	}

	return true, nil
}

func toBool(val interface{}) (bool, error) {

	b, ok := val.(bool)
	if !ok {
		s, ok := val.(string)

		if !ok {
			return false, fmt.Errorf("unable to convert to boolean")
		}

		var err error
		b, err = strconv.ParseBool(s)

		if err != nil {
			return false, err
		}
	}

	return b, nil
}
