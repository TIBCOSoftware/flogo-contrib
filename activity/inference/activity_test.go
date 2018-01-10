package inference

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework/tf"
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

	var features = make(map[string]interface{})
	features["z-axis-q75"] = 4.140586
	features["corr-x-z"] = 0.1381063882214782
	features["x-axis-mean"] = 1.7554575428900194
	features["z-axis-sd"] = 4.6888631696380765
	features["z-axis-skew"] = -0.3619011587545954
	features["y-axis-sd"] = -7.959084724314854
	features["y-axis-q75"] = 16.467001
	features["corr-z-y"] = 0.3467060369518231
	features["x-axis-sd"] = 6.450293741961166
	features["x-axis-skew"] = 0.09756801680727022
	features["y-axis-mean"] = 9.389463650669393
	features["y-axis-skew"] = -0.49036224958471764
	features["z-axis-mean"] = 1.1226106985139188
	features["x-axis-q25"] = -3.1463003
	features["x-axis-q75"] = 6.3198414
	features["y-axis-q25"] = 3.0645783
	features["z-axis-q25"] = -1.9477097
	features["corr-x-y"] = 0.08100326860866637

	tc.SetInput("features", features)

	done, _ := act.Eval(tc)
	if done == false {
		assert.Fail(t, "Invalid framework specified")
	}

	//check result attr
	fmt.Println(tc.GetOutput("result"))
}
