package actreturn

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
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

func TestSimpleReturn(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	ac := newActionContext()
	tc := test.NewTestActivityContextWithAction(getActivityMetadata(), ac)

	//set mappings
	mappingsJson := `[
      { "type": 2, "value": "1", "mapTo": "Output1" },
      { "type": 2, "value": 2, "mapTo": "Output2" }
    ]`

	var mappings interface{}
	err := json.Unmarshal([]byte(mappingsJson), &mappings)
	if err != nil {
		panic("Unable to parse mappings: " + err.Error())
	}

	//setup attrs
	tc.SetInput("mappings", mappings)

	//eval
	act.Eval(tc)

	assert.Nil(t, ac.ReplyErr)
	o1, exists1 := ac.ReplyDataAttr["Output1"]
	assert.True(t, exists1, "Output1 not set")
	if exists1 {
		assert.Equal(t, "1", o1.Value())
	}
	o2, exists2 := ac.ReplyDataAttr["Output2"]
	assert.True(t, exists2, "Output2 not set")
	if exists2 {
		assert.Equal(t, 2.0, o2.Value())
	}
}

func newActionContext() *test.TestActivityHost {

	input := []*data.Attribute{data.NewZeroAttribute("Input1", data.TypeString)}
	output := []*data.Attribute{data.NewZeroAttribute("Output1", data.TypeString), data.NewZeroAttribute("Output2", data.TypeInteger)}

	ac := &test.TestActivityHost{
		HostId:     "1",
		HostRef:    "github.com/TIBCOSoftware/flogo-contrib/action/flow",
		IoMetadata: &data.IOMetadata{Input: input, Output: output},
		HostData:   data.NewSimpleScope(nil, nil),
	}

	return ac
}
