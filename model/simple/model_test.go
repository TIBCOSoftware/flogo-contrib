package simple

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/model"
)

func TestRegistered(t *testing.T) {
	act := model.Get("simple")

	if act == nil {
		t.Error("Model Not Registered")
		t.Fail()
		return
	}
}
