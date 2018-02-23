package mapper

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

func TestSimpleMapper(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	ah := newActivityHost()
	tc := test.NewTestActivityContextWithAction(getActivityMetadata(), ah)

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

	//assert.Nil(t, ah.ReplyErr)
	o1, exists1 := ah.HostData.GetAttr("Output1")
	assert.True(t, exists1, "Output1 not set")
	if exists1 {
		assert.Equal(t, "1", o1.Value())
	}
	o2, exists2 := ah.HostData.GetAttr("Output2")
	assert.True(t, exists2, "Output2 not set")
	if exists2 {
		assert.Equal(t, 2, o2.Value())
	}
}

func newActivityHost() *test.TestActivityHost {
	input := []*data.Attribute{data.NewZeroAttribute("Input1", data.STRING)}
	output := []*data.Attribute{data.NewZeroAttribute("Output1", data.STRING), data.NewZeroAttribute("Output2", data.INTEGER)}

	ac := &test.TestActivityHost{
		HostId:     "1",
		HostRef:    "github.com/TIBCOSoftware/flogo-contrib/action/flow",
		IoMetadata: &data.IOMetadata{Input: input, Output: output},
		HostData:   data.NewFixedScope(output),
	}

	return ac
}
