package extension

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/model/simple"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
)

//Provider is the extension provider for the flow action
type Provider struct {
}

func New() *Provider {
	return &Provider{}
}

func (fp *Provider) GetFlowProvider() definition.Provider {
	return definition.NewRemoteFlowProvider()
}

func (fp *Provider) GetFlowModel() *model.FlowModel {
	return simple.New()
}

func (fp *Provider) GetStateRecorder() flowinst.StateRecorder {
	return nil
}
