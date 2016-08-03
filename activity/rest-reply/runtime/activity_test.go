package rest

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo-lib/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-restreply")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestSimpleReply(t *testing.T) {

	act := activity.Get("tibco-restreply")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("code", 200)
	//tc.SetInput("data", "")

	//eval
	act.Eval(tc)
}
