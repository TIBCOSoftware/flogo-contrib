package app

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
	"strings"
	"github.com/TIBCOSoftware/flogo-lib/core/app"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-app")

const (
	ivAttrName = "attribute"
	ivOp = "operation"
	ivType  = "type"
	ivValue = "value"

	ovValue = "value"
)

// AppActivity is a App Activity implementation
type AppActivity struct {
	sync.Mutex
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&AppActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *AppActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AppActivity) Eval(context activity.Context) (done bool, evalError *activity.Error) {

	attrName := context.GetInput(ivAttrName).(string)
	op := strings.ToUpper(context.GetInput(ivOp).(string)) //ADD,UPDATE,GET
	
	
	switch op {
	case "ADD":
		dt, ok := data.ToTypeEnum(strings.ToLower(context.GetInput(ivType).(string)))

		if !ok {
			return false,  activity.NewError("Unsupported Type: " + context.GetInput(ivType).(string))
		}

		val := context.GetInput(ivValue)
		//data.CoerceToValue(val, dt)
		
		app.GetContext().AddAttr(attrName, dt, val)
		context.SetOutput(ovValue, val)
	case "GET":
		val, ok := app.GetContext().GetAttrValue(attrName)
		
		if !ok {
			return false,  activity.NewError("Attribute not defined: " + attrName)
		}
		
		context.SetOutput(ovValue, val)
	case "UPDATE":

		val := context.GetInput(ivValue)
		//data.CoerceToValue(val, dt)
		
		app.GetContext().SetAttrValue(attrName, val)
		context.SetOutput(ovValue, val)
	default:
		return false, activity.NewError("Unsupported Op: " + op)
	}
	
	return true, nil
}
