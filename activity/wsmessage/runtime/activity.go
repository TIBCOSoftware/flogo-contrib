package sendWSMessage

import (
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"net/url"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-sendWSMessage")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&MyActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval - Sends a message to a WebSocket enabled server like TIBCO eFTL
func (a *MyActivity) Eval(context activity.Context) (done bool, evalError *activity.Error) {

	// Get the activity data from the context
	wsHost := context.GetInput("Server").(string)
	wsChannel := context.GetInput("Channel").(string)
	wsDestination := context.GetInput("Destination").(string)
	wsMessage := context.GetInput("Message").(string)
	wsUser := context.GetInput("Username").(string)
	wsPassword := context.GetInput("Password").(string)

	wsURL := url.URL{Scheme: "ws", Host: wsHost, Path: wsChannel}
	log.Debugf("connecting to %s", wsURL.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Debugf("Error while dialing to wsHost: ", err)
	}

	loginMessage := `{"op": 1, "client_type": "js", "client_version": "3.0.0   V9", "user":"` + wsUser + `", "password":"` + wsPassword + `", "login_options": {"_qos": "true"}}`

	log.Debugf("Preparing to send login message: [%s]", loginMessage)

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(loginMessage))
	if err != nil {
		log.Debugf("Error while sending login message to wsHost: [%s]", err)
		return
	}

	textMessage := `{"op": 8, "body": {"_dest":"` + wsDestination + `", "text":"` + wsMessage + `"}, "seq": 1}`

	log.Debugf("Preparing to send message: [%s]", textMessage)

	err = wsConn.WriteMessage(websocket.TextMessage, []byte(textMessage))
	if err != nil {
		log.Debugf("Error while sending message to wsHost: [%s]", err)
		return
	}

	context.SetOutput("result", "OK")

	return true, nil
}
