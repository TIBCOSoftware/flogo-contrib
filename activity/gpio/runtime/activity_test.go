package salesforce

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-gpio")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestReadState(t *testing.T) {

	act := activity.Get("tibco-gpio")

	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "Read State")
	tc.SetInput("pin number", 10)
	//eval
	_, err := act.Eval(tc)
	if err != nil {
		log.Errorf("Error occured: %+v", err)
	}
	val := tc.GetOutput("result")
	log.Debugf("Resut %s", val)

}
