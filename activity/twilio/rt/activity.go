package twilio

import (
	"github.com/sfreiberg/gotwilio"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-twilio")

const (
	ivAcctSID   = "accountSID"
	ivAuthToken = "authToken"
	ivFrom      = "from"
	ivTo        = "to"
	ivMessage   = "message"
)

// TwilioActivity is a Twilio Activity implementation
type TwilioActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&TwilioActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *TwilioActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *TwilioActivity) Eval(context activity.Context) (done bool, evalError *activity.Error)  {

	accountSID := context.GetInput(ivAcctSID).(string)
	authToken := context.GetInput(ivAuthToken).(string)
	from := context.GetInput(ivFrom).(string)
	to := context.GetInput(ivTo).(string)
	message := context.GetInput(ivMessage).(string)

	twilio := gotwilio.NewTwilioClient(accountSID, authToken)

	resp, _, err :=twilio.SendSMS(from, to, message, "", "")

	if err != nil {
		log.Error("Error sending SMS:", err)
	}

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("Response:", resp)
	}
	
	return true, nil
}
