package inference

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-contrib/activity/ml-inference/framework/tf"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var _ tf.TensorflowModel

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

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("model", "/Users/mellis/Documents/IoT/models/tn_demo/Archive.zip")
	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("features", "z-axis-q75:4.140586,corr-x-z:0.1381063882214782,x-axis-mean:1.7554575428900194,z-axis-sd:4.6888631696380765,z-axis-skew:-0.3619011587545954,y-axis-sd:7.959084724314854,y-axis-q75:16.467001,corr-z-y:0.3467060369518231,x-axis-sd:6.450293741961166,x-axis-skew:0.09756801680727022,y-axis-mean:9.389463650669393,y-axis-skew:-0.49036224958471764,z-axis-mean:1.1226106985139188,x-axis-q25:-3.1463003,x-axis-q75:6.3198414,y-axis-q25:3.0645783,z-axis-q25:-1.9477097,corr-x-y:0.08100326860866637")

	done, _ := act.Eval(tc)
	if done == false {
		assert.Fail(t, "Invalid framework specified")
	}

	//check result attr
	fmt.Println(tc.GetOutput("result"))
}
