package log

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-log")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	act := activity.Get("tibco-log")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("message", "test message")
	tc.SetInput("flowInfo", true)

	act.Eval(tc)
}

func TestAddToFlow(t *testing.T) {

	act := activity.Get("tibco-log")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("message", "test message")
	tc.SetInput("flowInfo", true)
	tc.SetInput("addToFlow", true)

	act.Eval(tc)

	msg := tc.GetOutput("message")

	fmt.Println("Message: ", msg)

	if msg == nil {
		t.Fail()
	}
}
