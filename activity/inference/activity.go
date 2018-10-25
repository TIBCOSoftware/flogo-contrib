package inference

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework/tf"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/model"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var _ tf.TensorflowModel

// log is the default package logger
var log = logger.GetLogger("activity-tibco-inference")

// variables needed to persist model between inferences
var tfmodelmap map[string]*model.Model
var modelRunMutex sync.Mutex

const (
	ivModel     = "model"
	ivInputName = "inputName"
	ivFeatures  = "features"
	ivFramework = "framework"

	ivSigDef = "sigDefName"
	ivTag    = "tag"

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
	features := context.GetInput(ivFeatures).([]interface{})
	fw := context.GetInput(ivFramework).(string)

	tfFramework := framework.Get(fw)
	if tfFramework == nil {
		log.Errorf("%s framework not registered", fw)

		return false, fmt.Errorf("%s framework not registered", fw)
	}
	log.Debug("Loaded Framework: " + tfFramework.FrameworkTyp())

	// Defining the flags to be used to load model
	flags := model.ModelFlags{
		Tag:    context.GetInput("tag").(string),
		SigDef: context.GetInput("sigDefName").(string),
	}

	// if modelmap does not exist then make it
	if tfmodelmap == nil {
		tfmodelmap = make(map[string]*model.Model)
		log.Info("Making map of models with keys of 'ModelKey'.")
	}

	// check if this instance of tf model already exists if not load it
	modelKey := context.ActivityHost().Name() + ":" + context.Name() + context.GetInput("model").(string)
	log.Info("ModelKey is:", modelKey)
	if tfmodelmap[modelKey] == nil {
		tfmodelmap[modelKey], err = model.Load(modelName, tfFramework, flags)
		if err != nil {
			return false, err
		}

		// We should check types of features and TF expectations here
		// We should also check shapes

	} else {
		log.Debug("Model already loaded - skipping loading")
	}

	log.Debug("Incoming features: ")
	log.Debug(features)

	// Grab the input feature set and parse out the feature labels and values
	inputSample := make(map[string]interface{})
	for _, feat := range features {
		featMap := feat.(map[string]interface{})
		inputSample[featMap["name"].(string)] = featMap["data"]
	}

	log.Debug("Parsing of features completed")

	modelRunMutex.Lock()
	tfmodelmap[modelKey].SetInputs(inputSample)
	output, err := tfmodelmap[modelKey].Run(tfFramework)
	modelRunMutex.Unlock()
	if err != nil {
		return false, err
	}

	log.Debug("Model execution completed with result:")
	log.Info(output)

	if strings.Contains(tfmodelmap[modelKey].Metadata.Method, "tensorflow/serving/classify") {
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
