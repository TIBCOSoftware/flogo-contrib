package instance

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"sync"
	"github.com/TIBCOSoftware/flogo-lib/app"
	"time"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"fmt"
	"errors"
	"runtime/debug"
)

type Status string

const (
	Created   = "Created"
	Completed = "Completed"
	Cancelled = "Cancelled"
	Failed    = "Failed"
	Scheduled = "Scheduled"
	Skipped   = "Skipped"
	Started   = "Started"
	Waiting   = "Waiting"
	Unknown   = "Created"
)

type EventListenerFunc func(*EventContext) error

var eventListeners = make(map[string]EventListenerFunc)

// Buffered channel
var eventQueue = make(chan interface{}, 100)
var publisherRoutineStarted = false
var shutdown = make(chan bool)

var lock = &sync.RWMutex{}

// Registers listener for flow events
func RegisterEventListener(name string, fel EventListenerFunc) error {
	lock.Lock()
	defer lock.Unlock()
	_, exists := eventListeners[name]
	if exists {
		errMsg := fmt.Sprintf("Event listener with name - '%s' already registered", name)
		logger.Error(errMsg)
		return errors.New(errMsg)
	}
	eventListeners[name] = fel
	startPublisherRoutine()
	return nil
}

// Unregisters flow event listener
func UnRegisterEventListener(name string) {
	lock.Lock()
	defer lock.Unlock()
	_, exists := eventListeners[name]
	if exists {
		delete(eventListeners, name)
	}
	stopPublisherRoutine()
}

func startPublisherRoutine() {
	if publisherRoutineStarted == true {
		return
	}

	if len(eventListeners) > 0 {
		// start publisher routine
		go publishEvents()
		publisherRoutineStarted = true
	}
}

func stopPublisherRoutine() {
	if publisherRoutineStarted == false {
		return
	}

	if len(eventListeners) == 0 {
		// No more listeners. Stop go routine
		shutdown <- true
		publisherRoutineStarted = false
	}
}

//  EventContext is a wrapper over specific event context
type EventContext struct {
	// Event can be FlowEventContext or TaskEventContext
	event interface{}
}

// Returns wrapped event
func (ec *EventContext) GetEvent() interface{} {
	return ec.event
}

// FlowEventContext provides access to flow instance execution details
type FlowEventContext struct {
	time         time.Time
	status       Status
	flowInstance *Instance
}

// Returns flow name
func (fe *FlowEventContext) Name() string {
	return fe.flowInstance.Name()
}

// Returns flow ID
func (fe *FlowEventContext) ID() string {
	return fe.flowInstance.ID()
}

// In case of subflow, returns parent flow name
func (fe *FlowEventContext) ParentName() string {
	if fe.flowInstance.master != nil {
		return fe.flowInstance.master.Name()
	}
	return ""
}

// In case of subflow, returns parent flow ID
func (fe *FlowEventContext) ParentID() string {
	if fe.flowInstance.master != nil {
		return fe.flowInstance.master.ID()
	}
	return ""
}

// Returns event time
func (fe *FlowEventContext) Time() time.Time {
	return fe.time
}

// Returns application name
func (fe *FlowEventContext) AppName() string {
	return app.GetName()
}

// Returns application version
func (fe *FlowEventContext) AppVersion() string {
	return app.GetVersion()
}

// Returns current flow status
func (fe *FlowEventContext) Status() Status {
	return fe.status
}

//TODO: Should we read once?
// Returns flow runtime attributes
func (fe *FlowEventContext) GetWorkingData() map[string]interface{} {
	attrs := make(map[string]interface{})
	if fe.flowInstance.attrs != nil && len(fe.flowInstance.attrs) > 0 {
		for k, v := range fe.flowInstance.attrs {
			attrs[k] = v.Value()
		}
	}
	return attrs
}

// Returns output data for completed flow instance
func (fe *FlowEventContext) Output() map[string]interface{} {
	attrs := make(map[string]interface{})
	if fe.Status() == Completed && fe.flowInstance.returnData != nil && len(fe.flowInstance.returnData) > 0 {
		for k, v := range fe.flowInstance.returnData {
			attrs[k] = v.Value()
		}
	}
	return attrs
}

// Returns error for failed flow instance
func (fe *FlowEventContext) Error() error {
	return fe.flowInstance.returnError
}

// TaskEventContext provides access to task instance execution details
type TaskEventContext struct {
	time         time.Time
	status       Status
	taskInstance *TaskInst
}

// Returns flow name
func (te *TaskEventContext) FlowName() string {
	return te.taskInstance.flowInst.Name()
}

// Returns flow ID
func (te *TaskEventContext) FlowID() string {
	return te.taskInstance.flowInst.ID()
}

// Returns task name
func (te *TaskEventContext) Name() string {
	return te.taskInstance.task.Name()
}

// Returns task type
func (te *TaskEventContext) Type() string {
	return te.taskInstance.task.TypeID()
}

// Returns task status
func (te *TaskEventContext) Status() Status {
	return te.status
}

// Returns application name
func (te *TaskEventContext) AppName() string {
	return app.GetName()
}

// Returns application version
func (te *TaskEventContext) AppVersion() string {
	return app.GetVersion()
}

// Returns event time
func (te *TaskEventContext) Time() time.Time {
	return te.time
}

// Returns working data of current instance. e.g. key and value of current iteration for iterator task.
func (te *TaskEventContext) GetWorkingData() map[string]interface{} {
	attrs := make(map[string]interface{})
	if te.taskInstance.HasWorkingData() {
		for name, value := range te.taskInstance.workingData {
			attrs[name] = value.Value()
		}
	}
	return attrs
}

// Returns activity input data
func (te *TaskEventContext) ActivityInput() map[string]interface{} {
	attrs := make(map[string]interface{})
	if te.taskInstance.task.ActivityConfig().GetInputAttrs() != nil && te.taskInstance.inScope != nil {
		for name := range te.taskInstance.task.ActivityConfig().GetInputAttrs() {
			inVal, _ := te.taskInstance.inScope.GetAttr(name)
			attrs[name] = inVal
		}
	}
	return attrs
}

// Returns output data for completed activity
func (te *TaskEventContext) ActivityOutput() map[string]interface{} {
	attrs := make(map[string]interface{})
	if te.Status() == Completed && te.taskInstance.task.ActivityConfig().GetOutputAttrs() != nil && te.taskInstance.outScope != nil {
		for name := range te.taskInstance.task.ActivityConfig().GetOutputAttrs() {
			outVal, _ := te.taskInstance.outScope.GetAttr(name)
			attrs[name] = outVal
		}
	}
	return attrs
}

// Returns error for failed task
func (te *TaskEventContext) Error() error {
	return te.taskInstance.returnError
}

func convertFlowStatus(code model.FlowStatus) Status {
	switch code {
	case model.FlowStatusNotStarted:
		return Created
	case model.FlowStatusActive:
		return Started
	case model.FlowStatusCancelled:
		return Cancelled
	case model.FlowStatusCompleted:
		return Completed
	case model.FlowStatusFailed:
		return Failed
	}
	return Unknown
}

func convertTaskStatus(code model.TaskStatus) Status {
	switch code {
	case model.TaskStatusNotStarted:
		return Created
	case model.TaskStatusEntered:
		return Scheduled
	case model.TaskStatusSkipped:
		return Skipped
	case model.TaskStatusReady:
		return Started
	case model.TaskStatusFailed:
		return Failed
	case model.TaskStatusDone:
		return Completed
	case model.TaskStatusWaiting:
		return Waiting
	}
	return Unknown
}

func publishEvents() {
	for {
		select {
		case event := <-eventQueue:
			lock.RLock()
			evtContext := &EventContext{event: event}
			publishEvent(evtContext)
			lock.RUnlock()
		case <-shutdown:
			return
		}
	}
}

func publishEvent(fe *EventContext) {

	for name, fel := range eventListeners {
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("Registered event handler - '%s' failed to process event due to error - '%v' ", name, r)
					logger.Errorf("StackTrace: %s", debug.Stack())
				}
			}()
			err := fel(fe)
			if err != nil {
				logger.Errorf("Registered event handler - '%s' failed to process event due to error - '%s' ", name, err.Error())
			}
		}()
	}
}
