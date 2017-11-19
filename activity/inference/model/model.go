package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/utils"
)

// Framework interface used to implement specific ml framework implementations
type Framework interface {
	Load(modelPath string, modelFile string, model *Model) (err error)
	Run(model *Model) (out map[string]interface{}, err error)
	FrameworkTyp() string
}

// Model represents a ML model. The metadata defines inputs, outputs, shapes, etc.
// The Instance object is a pointer to the actual framework obj. This will be used by the framework impl to execute models
type Model struct {
	Metadata *Metadata
	Instance interface{}
	Inputs   map[string]map[string]interface{}
}

func Load(modelArchive string, framework Framework) (*Model, error) {
	f, _ := os.Open(modelArchive)
	defer f.Close()
	var outDir string
	if fi, err := f.Stat(); err == nil && fi.IsDir() {
		outDir = modelArchive
	} else {
		tmpDir := os.TempDir()
		outDir = filepath.Join(tmpDir, utils.PseudoUuid())
		if err := utils.Unzip(modelArchive, outDir); err != nil {
			return nil, fmt.Errorf("Failed to extract model archive: %v", err)
		}
	}

	var model Model
	framework.Load(outDir, filepath.Join(outDir, "saved_model.pb"), &model)

	return &model, nil
}

func (m *Model) AppendInput(inputName string, feature string, in interface{}) {
	if m.Inputs == nil {
		m.Inputs = make(map[string]map[string]interface{})
	}

	m.Inputs[inputName][feature] = in
}

func (m *Model) SetInputs(in map[string]map[string]interface{}) {
	m.Inputs = in
}

func (m *Model) RemoveInput(featureName string) {

}

func (m *Model) Run(framework Framework) (map[string]interface{}, error) {
	// check if inputs are available. maybe validate inputs against metadata?

	// run the model
	out, _ := framework.Run(m)

	return out, nil
}
