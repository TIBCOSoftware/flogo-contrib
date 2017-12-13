package awssns

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"testing"
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

func TestStringConfig(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("AWS_ACCESS_KEY_ID", "")
	tc.SetInput("AWS_SECRET_ACCESS_KEY", "")
	tc.SetInput("AWS_DEFAULT_REGION", "ap-southeast-2")
	tc.SetInput("SMS_TYPE", "Promotional")
	tc.SetInput("SMS_FROM", "Sender")
	tc.SetInput("SMS_TO", "+XXXXXXXXXXXX")
	tc.SetInput("SMS_MESSAGE", "Hello world !")

	success, err := act.Eval(tc)

	if err != nil {
		t.Error("Error while sending SMS")
		t.Fail()
		return
	}
	if success {
		val := tc.GetOutput("MESSAGE_ID")
		fmt.Printf("Message ID : %v\n", val)
	} else {
		t.Error("Error while sending SMS")
		t.Fail()
		return
	}

}
