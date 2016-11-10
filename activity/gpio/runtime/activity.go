package salesforce

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/op/go-logging"
	"errors"
	"github.com/stianeikeland/go-rpio"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-rest")

const (
	method = "method"
	pinNumber = "pinNumber"
	directionState = "direction"
	state = "state"
	direction = "Direction"
	setState = "Set State"
	readState = "Read State"
	pull = "Pull"

	input = "Input"
	output = "Output"

	high = "High"
	low = "Low"

	up = "Up"
	down = "Down"
	off = "off"

	//ouput

	reust = "result"
)

type GPIOActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	log.Info("Init and start init activities")
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&GPIOActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *GPIOActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *GPIOActivity) Eval(context activity.Context) (done bool, err error) {
	//getmethod
	methodInput := context.GetInput(method)

	if methodInput == nil {
		return true, errors.New("Method field not set.")
	}

	//get pinNumber
	pinNumberInput := context.GetInput(pinNumber)

	if pinNumberInput == nil {
		return true, errors.New("Pin number must exist")
	}

	ivPinNumber := pinNumberInput.(int)

	log.Debugf("Method '%s' and pin number '%d'", methodInput, ivPinNumber)
	//Open pin
	openerr := rpio.Open()
	if openerr != nil {
		return true, errors.New("Open RPIO error: "+ openerr.Error())
	}

	pin := rpio.Pin(ivPinNumber)

	ivmethod := methodInput.(string)

	switch ivmethod {
	case direction:
		directionStateInput := context.GetInput(directionState)
		if directionStateInput == nil {
			return true, errors.New("Direction field not set.")
		}

		ivDirectionField := directionStateInput.(string)

		if strings.EqualFold(input, ivDirectionField) {
			pin.Input()
		}else {
			pin.Output()
		}
	case setState:
		stateInput := context.GetInput(state)
		if stateInput == nil {
			return true, errors.New("State field not set.")
		}

		ivstate := stateInput.(string)

		if strings.EqualFold(high, ivstate) {
			pin.High()
		}else {
			pin.Low()
		}
	case readState:
		log.Debugf("Read state and state: %s", readState)

		readState := pin.Read()
		log.Debugf("Read state and state: %s", readState)
		context.SetOutput("result", readState)
	case pull:
		pullInput := context.GetInput(pull)
		if pullInput == nil {
			return true, errors.New("Pull field not set.")
		}

		ivpull := pullInput.(string)

		if strings.EqualFold(up, ivpull) {
			pin.PullUp()
		}else if strings.EqualFold(down, ivpull) {
			pin.PullDown()
		}else {
			pin.PullOff()
		}
	default:
		log.Errorf("Cannot found method %s ", ivmethod)
		return true, errors.New("Cannot found method %s " + ivmethod)
	}

	context.SetOutput("result", "done")
	return true, nil
}