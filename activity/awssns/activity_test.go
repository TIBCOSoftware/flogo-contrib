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
	tc.SetInput("accessKey", "")
	tc.SetInput("secretKey", "")
	tc.SetInput("region", "ap-southeast-2")
	tc.SetInput("smsType", "Promotional")
	tc.SetInput("from", "Sender")
	tc.SetInput("to", "+XXXXXXXXXXXX")
	tc.SetInput("message", "Hello world !")

	success, err := act.Eval(tc)

	if err != nil {
		t.Error("Error while sending SMS")
		t.Fail()
		return
	}
	if success {
		val := tc.GetOutput("messageId")
		fmt.Printf("Message ID : %v\n", val)
	} else {
		t.Error("Error while sending SMS")
		t.Fail()
		return
	}

}
