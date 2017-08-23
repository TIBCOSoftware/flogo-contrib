package app

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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

func TestAdd(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "ADD")
	tc.SetInput(ivType, "string")
	tc.SetInput(ivValue, "test")

	act.Eval(tc)

	value, found := tc.GetOutput(ovValue).(string)

	assert.True(t, found, "not found")
	if found {
		assert.Equal(t, "test", value, "not equal")
	}
}

func TestGet(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//add attribute
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "GET")

	act.Eval(tc)

	value, found := tc.GetOutput(ovValue).(string)

	assert.True(t, found, "not found")
	if found {
		assert.Equal(t, "test", value, "not equal")
	}
}

func TestUpdate(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "UPDATE")
	tc.SetInput(ivValue, "test3")

	act.Eval(tc)

	value, found := tc.GetOutput(ovValue).(string)

	assert.True(t, found, "not found")
	if found {
		assert.Equal(t, "test3", value, "not equal")
	}
}
