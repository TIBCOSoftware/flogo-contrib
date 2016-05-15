package awsiot

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-awsiot")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}
