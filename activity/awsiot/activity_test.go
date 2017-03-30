package awsiot

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
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
	act := activity.Get("github.com/TIBCOSoftware/flogo-contrib/activity/awsiot")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}
