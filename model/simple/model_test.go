package simple

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/model"
)

func TestRegistered(t *testing.T) {
	act := model.Get("tibco-simple")

	if act == nil {
		t.Error("Model Not Registered")
		t.Fail()
		return
	}
}
