package log

import (
	"fmt"
	"os"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
)

// activityLog is the default logger for the Log Activity
var activityLog = logging.MustGetLogger("activity-log")

// format is the log format for the Activity log
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)

var backend = logging.NewLogBackend(os.Stdout, "", 0)
var backendFormatter = logging.NewBackendFormatter(backend, format)
var backendLeveled = logging.AddModuleLevel(backendFormatter)

func init() {
	backendLeveled.SetLevel(logging.INFO, "")
	activityLog.SetBackend(backendLeveled)
}

// LogActivity is an Activity that is used to log a message to the console
// inputs : {message, processInfo}
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
func (a *LogActivity) Eval(context activity.Context) bool {

	message, _ := context.GetAttrValue("message")
	processInfo, _ := context.GetAttrValue("processInfo")

	msg := message

	showInfo, _ := strconv.ParseBool(processInfo)

	if showInfo {

		msg = fmt.Sprintf("%s - ProcessInstanceID [%s], Process [%s], Task [%s]", msg,
			context.ProcessInstanceID(), context.ProcessName(), context.TaskName())
	}

	activityLog.Info(msg)

	//log.Debugf("%s: %s\n", time.Now(), msg)

	return true
}
