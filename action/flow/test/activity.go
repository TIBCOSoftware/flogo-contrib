package test

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func init() {
	activity.Register(NewLogActivity())
	activity.Register(NewCounterActivity())
}

type LogActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewLogActivity() activity.Activity {
	metadata := &activity.Metadata{ID: "test-log"}
	input := map[string]*data.Attribute{
		"message": data.NewZeroAttribute("message", data.STRING),
	}
	metadata.Input = input
	return &LogActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *LogActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *LogActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	message, _ := context.GetInput("message").(string)

	fmt.Println("Message :", message)
	return true, nil
}


type CounterActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewCounterActivity() activity.Activity {
	metadata := &activity.Metadata{ID: "test-counter"}
	input := map[string]*data.Attribute{
		"counterName": data.NewZeroAttribute("counterName", data.STRING),
	}
	metadata.Input = input
	return &CounterActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *CounterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *CounterActivity) Eval(context activity.Context) (done bool, err error) {

	//mv := context.GetInput(ivMessage)
	counterName, _ := context.GetInput("counterName").(string)
	fmt.Println("counterName :", counterName)


	return true, nil
}
