package actreturn

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"encoding/json"
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
	o1,exists1 := ac.ReplyDataAttr["Output1"]
	assert.True(t, exists1, "Output1 not set")
	if exists1 {
		assert.Equal(t, "1", o1.Value)
	}
	o2,exists2 := ac.ReplyDataAttr["Output2"]
	assert.True(t, exists2, "Output2 not set")
	if exists2 {
		assert.Equal(t, 2.0, o2.Value)
	}
}

func newActionContext() *test.TestActionCtx {
	input := []*data.Attribute{{Name: "Input1", Type: data.STRING}}
	output := []*data.Attribute{{Name: "Output1", Type: data.STRING}, {Name: "Output2", Type: data.INTEGER}}

	ac := &test.TestActionCtx{
		ActionId:   "1",
		ActionRef:  "github.com/TIBCOSoftware/flogo-contrib/action/flow",
		ActionMd:   &action.ConfigMetadata{Input: input, Output: output},
		ActionData: data.NewSimpleScope(nil, nil),
	}

	return ac
}
