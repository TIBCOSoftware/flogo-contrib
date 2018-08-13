package filter

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetSetting(sType, "non-zero")
	tc.SetInput(ivValue, 2)

	done, err := act.Eval(tc)

	if err != nil {
		t.Error("Error evaluating activity")
		t.Fail()
		return
	}

	if !done {
		t.Error("activity should be done")
		t.Fail()
		return
	}

	report := tc.GetOutput(ovValue).(bool)
	result := tc.GetOutput(ovFiltered)

	if result != 2 {
		t.Errorf("Result is %d instead of 2", result)
	}

	if !report {
		t.Error("value should be reported")
	}

	tc.SetSetting(sProceedOnlyOnEmit, false)
	tc.SetInput(ivValue, 0)

	done, err = act.Eval(tc)

	if err != nil {
		t.Error("Error evaluating activity")
		t.Fail()
		return
	}

	if !done {
		t.Error("activity should be done")
		t.Fail()
		return
	}

	report = tc.GetOutput(ovValue).(bool)
	result = tc.GetOutput(ovFiltered)

	if result != 0 {
		t.Errorf("Result is %d instead of 0", result)
	}

	if report {
		t.Error("value should not be reported")
	}

	tc.SetSetting(sProceedOnlyOnEmit, true)
	tc.SetInput(ivValue, 0)

	done, err = act.Eval(tc)

	if err != nil {
		t.Error("Error evaluating activity")
		t.Fail()
		return
	}

	if done {
		t.Error("activity should not be done")
		t.Fail()
		return
	}

	report = tc.GetOutput(ovValue).(bool)
	result = tc.GetOutput(ovFiltered)

	if result != 0 {
		t.Errorf("Result is %d instead of 0", result)
	}

	if report {
		t.Error("value should not be reported")
	}
}
