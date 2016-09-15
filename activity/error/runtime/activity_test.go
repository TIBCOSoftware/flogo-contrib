package rest

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-error")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestSimpleError(t *testing.T) {

	act := activity.Get("tibco-error")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("message", "test error")

	//eval
	_, err := act.Eval(tc)

	if err == nil {
		t.Fail()
	}
}
