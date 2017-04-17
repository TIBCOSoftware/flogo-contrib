package extension

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/tester"
	"github.com/TIBCOSoftware/flogo-contrib/model/simple"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
)

//Provider is the extension provider for the flow action
type Provider struct {
	flowProvider definition.Provider
	flowModel    *model.FlowModel
}

func New() *Provider {
	return &Provider{}
}

func (fp *Provider) GetFlowProvider() definition.Provider {

	if fp.flowProvider == nil {
		fp.flowProvider = definition.NewRemoteFlowProvider()
	}

	return fp.flowProvider
}

func (fp *Provider) GetFlowModel() *model.FlowModel {

	if fp.flowModel == nil {
		fp.flowModel = simple.New()
	}

	return fp.flowModel
}

func (fp *Provider) GetStateRecorder() instance.StateRecorder {
	return nil
}

func (fp *Provider) GetMapperFactory() flowdef.MapperFactory {
	return nil
}

func (fp *Provider) GetLinkExprManagerFactory() flowdef.LinkExprManagerFactory {
	return nil
}

//todo make FlowTester an interface
func (fp *Provider) GetFlowTester() *tester.RestEngineTester {
	return nil
}