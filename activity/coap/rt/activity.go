package coap

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
	"github.com/dustin/go-coap"
	"strconv"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-coap")

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodDELETE = "DELETE"

	typeCON = "CONFIRMABLE"
	typeNON = "NONCONFIRMABLE"
	typeACK = "ACKNOWLEDGEMENT"
	typeRST = "RESET"

	ivAddress   = "address"
	ivMethod    = "method"
	ivType      = "type"
	ivPayload   = "payload"
	ivMessageId = "messageId"
	ivOptions   = "options"

	ovResponse = "response"
)

var validMethods = []string{methodGET, methodPOST, methodPUT, methodDELETE}
var validTypes = []string{typeCON, typeNON}

// CoAPActivity is an Activity that is used to send a CoAP message
// inputs : {method,type,payload,messageId}
// outputs: {result}
type CoAPActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&CoAPActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *CoAPActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *CoAPActivity) Eval(context activity.Context) (done bool, evalError *activity.Error) {

	method := strings.ToUpper(context.GetInput(ivMethod).(string))
	msgType := strings.ToUpper(context.GetInput(ivType).(string))
	address := context.GetInput(ivAddress).(string)
	payload := context.GetInput(ivPayload).(string)
	messageId := context.GetInput(ivMessageId).(int)

	val := context.GetInput(ivOptions)

	var options map[string]string

	if val != nil {
		options = val.(map[string]string)
	}

	log.Debugf("CoAP Message: [%s] %s\n", method, payload)

	req := coap.Message{
		Type:      toCoapType(msgType),
		Code:      toCoapCode(method),
		MessageID: uint16(messageId),
		Payload:   []byte(payload),
	}

	if options != nil {
		for k, v := range options {
			op, val := toOption(k, v)
			req.SetOption(op, val)
		}
	}

	c, err := coap.Dial("udp", address)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Debugf("Response payload: %s", rv.Payload)
	}

	context.SetOutput(ovResponse, string(rv.Payload))

	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

func toCoapCode(method string) coap.COAPCode {

	var code coap.COAPCode

	switch method {
	case methodGET:
		code = coap.GET
	case methodPOST:
		code = coap.POST
	case methodPUT:
		code = coap.PUT
	case methodDELETE:
		code = coap.DELETE
	}

	return code
}

func toCoapType(typeStr string) coap.COAPType {

	var ctype coap.COAPType

	switch typeStr {
	case typeCON:
		ctype = coap.Confirmable
	case typeNON:
		ctype = coap.NonConfirmable
	case typeACK:
		ctype = coap.Acknowledgement
	case typeRST:
		ctype = coap.Reset
	}

	return ctype
}

func toOption(name string, value string) (coap.OptionID, interface{}) {

	var opID coap.OptionID
	var val interface{}

	val = value

	switch name {
	case "IFMATCH":
		opID = coap.IfMatch
	case "URIHOST":
		opID = coap.URIHost
	case "ETAG":
		opID = coap.ETag
	//case "IFNONEMATCH":
	//	opID = coap.IfNoneMatch
	case "OBSERVE":
		opID = coap.Observe
		val,_ = strconv.Atoi(value)
	case "URIPORT":
		opID = coap.URIPort
		val,_ = strconv.Atoi(value)
	case "LOCATIONPATH":
		opID = coap.LocationPath
	case "URIPATH":
		opID = coap.URIPath
	case "CONTENTFORMAT":
		opID = coap.ContentFormat
		val,_ = strconv.Atoi(value)
	case "MAXAGE":
		opID = coap.MaxAge
		val,_ = strconv.Atoi(value)
	case "URIQUERY":
		opID = coap.URIQuery
	case "ACCEPT":
		opID = coap.IfMatch
		val,_ = strconv.Atoi(value)
	case "LOCATIONQUERY":
		opID = coap.LocationQuery
	case "PROXYURI":
		opID = coap.ProxyURI
	case "PROXYSCHEME":
		opID = coap.ProxyScheme
	case "SIZE1":
		opID = coap.Size1
		val,_ = strconv.Atoi(value)
	default:
		opID = 0
		val = nil
	}

	return opID, val
}

func methodIsValid(method string) bool {

	if !stringInList(method, validMethods) {
		return false
	}

	return true
}

func stringInList(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}