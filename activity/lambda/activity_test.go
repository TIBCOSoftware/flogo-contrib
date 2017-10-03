package lambda

import (
	"testing"

	"io/ioutil"

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

func TestLambdaInvokeWithSecurity(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("arn", "arn:aws:lambda:us-east-1:658833855929:function:myFlogoApp")
	tc.SetInput("region", "us-east-1")
	tc.SetInput("accessKey", "AKIAJ2HZLVTLTDSIPJCA")
	tc.SetInput("secretKey", "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")
	tc.SetInput("payload", "hello")

	//eval
	_, err := act.Eval(tc)

	if err == nil {
		t.Fail()
	}
}
