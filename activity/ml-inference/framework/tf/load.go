package tf

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/proto"
	models "github.com/TIBCOSoftware/flogo-contrib/activity/ml-inference/model"
	tfpb "github.com/TIBCOSoftware/flogo-contrib/activity/ml-inference/tensorflow/tensorflow/core/protobuf"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Load implements the backend framework specifics for loading a saved model
func (a *TensorflowModel) Load(modelPath string, modelFile string, model *models.Model) (err error) {
	var meta models.Metadata

	// Parse the protobuffer
	parseProtoBuf(modelFile, &meta)
	model.Metadata = &meta

	bundle, err := tf.LoadSavedModel(modelPath, []string{"serve"}, nil)
	if err != nil {
		return err
	}
	model.Instance = bundle

	return nil
}

func parseProtoBuf(file string, model *models.Metadata) error {
	savedModelPb, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var savedModel tfpb.SavedModel
	loadErr := proto.Unmarshal(savedModelPb, &savedModel)
	if loadErr != nil {
		fmt.Println(loadErr)
	}
	metaGraphs := savedModel.GetMetaGraphs()

	// Grab the default graph def
	sigDef := metaGraphs[0].SignatureDef["serving_default"]

	// Collect inputs
	inputs := getValues(sigDef.GetInputs())
	outputs := getValues(sigDef.GetOutputs())

	// Determine the feature keys
	if model.Inputs.Features == nil {
		model.Inputs.Features = make(map[string]models.Feature)
	}
	var featureIndx = 0
	for _, node := range metaGraphs[0].GetGraphDef().GetNode() {
		// Grab dense features
		if strings.Contains(node.GetName(), "ParseExample/ParseExample/dense_keys") {
			k := string(node.GetAttr()["value"].GetTensor().GetStringVal()[0][:])

			// Now determine the shape and type of each feature
			for _, n2 := range metaGraphs[0].GetGraphDef().GetNode() {
				if n2.GetName() == "ParseExample/ParseExample" {
					var featureTyp string
					var featureShape []int64

					for attr, val := range n2.GetAttr() {

						if attr == "Tdense" {
							featureTyp = val.GetList().GetType()[featureIndx].String()
						} else if attr == "_output_shapes" {
							for i := 0; i < len(val.GetList().GetShape()[featureIndx].Dim); i++ {
								featureShape = append(featureShape, val.GetList().GetShape()[featureIndx].Dim[i].GetSize())
							}
						}
					}
					feat := models.Feature{
						Shape: featureShape,
						Type:  featureTyp,
					}

					model.Inputs.Features[k] = feat
				}
			}

			featureIndx++
		}
	}

	model.Inputs.Params = inputs
	model.Outputs = outputs

	return nil
}

// Used to extract input and output ops and data from the singdef in the pb
func getValues(sigDef map[string]*tfpb.TensorInfo) map[string]models.OperationParam {

	params := make(map[string]models.OperationParam)
	var i = 0
	for key, ins := range sigDef {
		var p models.OperationParam
		p.Name = strings.Split(ins.GetName(), ":")[0]
		p.Type = ins.GetDtype().String()

		// grab the shape
		for _, dim := range ins.GetTensorShape().GetDim() {
			p.Shape = append(p.Shape, dim.GetSize())
		}

		params[key] = p
		i++
	}

	return params
}
