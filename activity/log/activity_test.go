package log

import (
	"fmt"
	"testing"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

var jsonMetadataBytes, _ = ioutil.ReadFile("activity.json")
var jsonMetadata = string(jsonMetadataBytes)

func TestRegistered(t *testing.T) {
	act := activity.Get("github.com/TIBCOSoftware/flogo-contrib/activity/log")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	act := activity.Get("github.com/TIBCOSoftware/flogo-contrib/activity/log")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("message", "test message")
	tc.SetInput("flowInfo", true)

	act.Eval(tc)
}

func TestAddToFlow(t *testing.T) {

	act := activity.Get("github.com/TIBCOSoftware/flogo-contrib/activity/log")
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
