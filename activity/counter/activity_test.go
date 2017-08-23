package counter

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
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

func TestIncrement(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivCounterName, "messages")
	tc.SetInput(ivIncrement, true)

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(int)

	if value != 1 {
		t.Fail()
	}
}

func TestGet(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := getActivityMetadata()

	counters := map[string]int{
		"messages": 5,
	}

	act := &CounterActivity{metadata: md, counters: counters}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivCounterName, "messages")

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(int)

	if value != 5 {
		t.Fail()
	}
}

func TestReset(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := getActivityMetadata()
	counters := map[string]int{
		"messages": 3,
	}

	act := &CounterActivity{metadata: md, counters: counters}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivCounterName, "messages")
	tc.SetInput(ivReset, true)

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(int)

	if value != 0 {
		t.Fail()
	}
}
