package app

import (
	"strings"
	"sync"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)


const (
	ivAttrName = "attribute"
	ivOp       = "operation"
	ivType     = "type"
	ivValue    = "value"

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
func (a *AppActivity) Eval(context activity.Context) (done bool, err error) {

	attrName := context.GetInput(ivAttrName).(string)
	op := strings.ToUpper(context.GetInput(ivOp).(string)) //ADD,UPDATE,GET

	switch op {
	case "ADD":
		logger.Debug("In ADD operation")
		dt, ok := data.ToTypeEnum(strings.ToLower(context.GetInput(ivType).(string)))

		if !ok {
			errorMsg := fmt.Sprintf("Unsupported type '%s'", context.GetInput(ivType).(string))
			logger.Error(errorMsg)
			return false, activity.NewError(errorMsg)
		}

		val := context.GetInput(ivValue)
		//data.CoerceToValue(val, dt)

		data.GetGlobalScope().AddAttr(attrName, dt, val)
		context.SetOutput(ovValue, val)
	case "GET":
		logger.Debug("In GET operation")
		typedVal, ok := data.GetGlobalScope().GetAttr(attrName)

		if !ok {
			errorMsg := fmt.Sprintf("Attribute not defined: " + attrName)
			logger.Error(errorMsg)
			return false, activity.NewError(errorMsg)
		}

		context.SetOutput(ovValue, typedVal.Value)
	case "UPDATE":
		logger.Debug("In UPDATE operation")
		val := context.GetInput(ivValue)
		//data.CoerceToValue(val, dt)

		data.GetGlobalScope().SetAttrValue(attrName, val)
		context.SetOutput(ovValue, val)
	default:
		errorMsg := fmt.Sprintf("Unsupported Op: " + op)
		logger.Error(errorMsg)
		return false, activity.NewError(errorMsg)
	}

	return true, nil
}
