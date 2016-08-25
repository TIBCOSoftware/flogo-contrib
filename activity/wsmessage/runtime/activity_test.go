package sendWSMessage

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("sendWSMessage")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}
