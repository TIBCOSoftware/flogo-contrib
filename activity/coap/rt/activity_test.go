package coap

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-coap")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}
