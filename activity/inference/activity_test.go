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
	tc.SetInput("model", "/Users/avanders/working/working_python/accelerometer/models/TB/1531761580")
	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")

	var features = make(map[string]interface{})
	features["x"] = -1.961143
	features["y"] = 2.371200
	features["z"] = 9.746658
	features["x1"] = -1.859383
	features["y1"] = 2.881411
	features["z1"] = 9.595146
	features["x2"] = -1.892647
	features["y2"] = 3.019059
	features["z2"] = 9.418503
	features["x3"] = -2.028000
	features["y3"] = 2.666912
	features["z3"] = 8.900028
	features["x4"] = -2.278714
	features["y4"] = 4.023686
	features["z4"] = 8.547684
	features["x5"] = -2.377852
	features["y5"] = 4.164970
	features["z5"] = 8.326502
	features["x6"] = -2.422589
	features["y6"] = 4.343912
	features["z6"] = 8.241617
	features["x7"] = -2.496000
	features["y7"] = 4.807030
	features["z7"] = 7.697485
	features["x8"] = -2.586618
	features["y8"] = 4.956441
	features["z8"] = 7.407708
	features["x9"] = -3.515735
	features["y9"] = 4.490735
	features["z9"] = 7.653178
	features["x10"] = -4.226486
	features["y10"] = 4.313400
	features["z10"] = 8.797283
	features["activity"] = 1

	tc.SetInput("features", features)

	done, _ := act.Eval(tc)
	if done == false {
		assert.Fail(t, "Invalid framework specified")
	}

	//check result attr
	fmt.Println(tc.GetOutput("result"))
}
