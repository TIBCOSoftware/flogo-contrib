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
	//output = "Output"

	high = "High"
	//low = "Low"

	up = "Up"
	down = "Down"
	//off = "off"

	//ouput

	result = "result"
)

type GPIOActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
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
	log.Debug("Running gpio activity.")
	methodInput := context.GetInput(method)

	ivmethod, ok := methodInput.(string)
	if !ok {
		return true, errors.New("Method field not set.")
	}

	//get pinNumber
	ivPinNumber, ok := context.GetInput(pinNumber).(int)

	if !ok {
		return true, errors.New("Pin number must exist")
	}

	log.Debugf("Method '%s' and pin number '%d'", methodInput, ivPinNumber)
	//Open pin
	openErr := rpio.Open()
	if openErr != nil {
		log.Errorf("Open RPIO error: %+v", openErr.Error())
		return true, errors.New("Open RPIO error: " + openErr.Error())
	}

	pin := rpio.Pin(ivPinNumber)

	switch ivmethod {
	case direction:
		ivDirectionField, ok := context.GetInput(directionState).(string)
		if !ok {
			return true, errors.New("Direction field not set.")
		}
		if strings.EqualFold(input, ivDirectionField) {
			log.Debugf("Set pin %d direction to input", pin)
			pin.Input()
		} else {
			log.Debugf("Set pin %d direction to output", pin)
			pin.Output()
		}
	case setState:
		ivState, ok := context.GetInput(state).(string)
		if !ok {
			return true, errors.New("State field not set.")
		}

		if strings.EqualFold(high, ivState) {
			log.Debugf("Set pin %d state to High", pin)
			pin.High()
		} else {
			log.Debugf("Set pin %d state to low", pin)
			pin.Low()
		}
	case readState:
		log.Debugf("Read pin %d state..", pin)
		state := pin.Read()
		log.Debugf("Read state and state: %s", state)
		context.SetOutput(result, int(state))
		return true, nil
	case pull:
		ivPull, ok := context.GetInput(pull).(string)
		if !ok {
			return true, errors.New("Pull field not set.")
		}

		if strings.EqualFold(up, ivPull) {
			log.Debugf("Pull pin %d  to Up", pin)
			pin.PullUp()
		} else if strings.EqualFold(down, ivPull) {
			log.Debugf("Pull pin %d to Down", pin)
			pin.PullDown()
		} else {
			log.Debugf("Pull pin %d to Up", pin)
			pin.PullOff()
		}
	default:
		log.Errorf("Cannot found method %s ", ivmethod)
		return true, errors.New("Cannot found method %s " + ivmethod)
	}

	context.SetOutput(result, 0)
	return true, nil
}