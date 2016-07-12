package sendWSMessage

import (
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo-lib/test"
	"testing"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("sendWSMessage")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}
