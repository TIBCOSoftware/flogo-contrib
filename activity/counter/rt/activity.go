package counter

import (
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
	"sync"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-counter")

const (
	ivCounterName   = "counterName"
	ivIncrement = "increment"
	ivReset = "reset"

	ovValue = "value"
)


// CounterActivity is a Counter Activity implementation
type CounterActivity struct {
	sync.Mutex
	metadata *activity.Metadata
	counters map[string]int
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&CounterActivity{metadata: md, counters:make(map[string]int)})
}

// Metadata implements activity.Activity.Metadata
func (a *CounterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *CounterActivity) Eval(context activity.Context) (done bool, evalError *activity.Error)  {

	counterName := context.GetInput(ivCounterName).(string)

	var increment,reset bool

	if context.GetInput(ivIncrement) != nil {
		increment = context.GetInput(ivIncrement).(bool)
	}
	if context.GetInput(ivReset) != nil {
		reset = context.GetInput(ivReset).(bool)
	}

	var count int

	if increment {
		count = a.incrementCounter(counterName)

		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("Counter [%s] incremented: %d", counterName, count)
		}
	} else if reset {
		count = a.resetCounter(counterName)

		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("Counter [%s] reset", counterName)
		}
	} else {
		count = a.getCounter(counterName)

		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("Counter [%s] = %d", counterName, count)
		}
	}

	context.SetOutput(ovValue, count)

	return true, nil
}

func (a *CounterActivity) incrementCounter(counterName string) int {
	a.Lock()
	defer a.Unlock()

	count := 1

	if counter, exists := a.counters[counterName]; exists {
		count = counter + 1
	}

	a.counters[counterName] = count

	return count
}

func (a *CounterActivity) resetCounter(counterName string) int {
	a.Lock()
	defer a.Unlock()

	if _, exists := a.counters[counterName]; exists {
		a.counters[counterName] = 0
	}

	return 0
}

func (a *CounterActivity) getCounter(counterName string) int {
	a.Lock()
	defer a.Unlock()

	return a.counters[counterName]
}