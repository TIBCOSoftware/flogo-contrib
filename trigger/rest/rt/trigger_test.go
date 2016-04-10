package rest

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
)

func TestRegistered(t *testing.T) {
	act := trigger.Get("rest")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}
