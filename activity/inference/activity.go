package inference

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework/tf"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/model"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var _ tf.TensorflowModel

// log is the default package logger
var log = logger.GetLogger("activity-tibco-inference")

const (
	ivModel     = "model"
	ivInputName = "inputName"
	ivFeatures  = "features"
	ivFramework = "framework"

	ovResult = "result"
)

// InferenceActivity is an Activity that is used to invoke a a ML Model using flogo-ml framework
type InferenceActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new InferenceActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &InferenceActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *InferenceActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Runs an ML model
func (a *InferenceActivity) Eval(context activity.Context) (done bool, err error) {

	modelName := context.GetInput(ivModel).(string)
	inputName := context.GetInput(ivInputName).(string)
	features := context.GetInput(ivFeatures)
	fw := context.GetInput(ivFramework).(string)

	tfFramework := framework.Get(fw)
	if tfFramework == nil {
		log.Errorf("%s framework not registered", fw)

		return false, fmt.Errorf("%s framework not registered", fw)
	}
	log.Debug("Loaded Framework: " + tfFramework.FrameworkTyp())

	model, _ := model.Load(modelName, tfFramework)

	// Grab the input feature set and parse out the feature labels and values
	inputSample := make(map[string]map[string]interface{})
	inputSample[inputName] = make(map[string]interface{})

	log.Debug("Incoming features: ")
	log.Debug(features)

	featureMap, ok := features.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("Cannot parse features, should be map[string]interface{}")
	}

	inputSample[inputName] = featureMap
	log.Debug("Parsing of features completed")

	model.SetInputs(inputSample)
	output, err := model.Run(tfFramework)

	if err != nil {
		return false, err
	}

	log.Debug("Model execution completed with result:")
	log.Info(output)

	if strings.Contains(model.Metadata.Method, "tensorflow/serving/classify") {
		var out = make(map[string]interface{})

		classes := output["classes"].([][]string)[0]
		scores := output["scores"].([][]float32)[0]

		for i := 0; i < len(classes); i++ {
			out[classes[i]] = scores[i]
		}
		context.SetOutput(ovResult, out)
	} else {
		context.SetOutput(ovResult, output)
	}

	return true, nil
}
