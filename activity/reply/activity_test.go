package reply

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"io/ioutil"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string{
	jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
	if err != nil{
		panic("No Json Metadata found for activity.json path")
	}
	return string(jsonMetadataBytes)
}

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-restreply")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestSimpleReply(t *testing.T) {

	act := activity.Get("tibco-restreply")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("code", 200)
	//tc.SetInput("data", "")

	//eval
	act.Eval(tc)
}
