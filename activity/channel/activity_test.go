package channel

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/engine/channels"
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

	channels.Add("test:5")
	ch := channels.Get("test")

	//setup attrs
	tc.SetSetting(sChannel, "test")
	tc.SetInput(ivValue, 2)

	done, err := act.Eval(tc)

	if !done {
		t.Error("activity should be done")
		return
	}

	if err != nil {
		t.Error("activity has an error: ", err)
		return
	}

	expected := 2
	found := <-ch

	if found != expected {
		t.Errorf("Expected %s, found %s", expected, found)
	}

	channels.Close()
}