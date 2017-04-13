package tester

import (
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/model/simple"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	ENV_ENABLED         = "TESTER_ENABLED"
	ENV_SETTING_PORT    = "TESTER_PORT"
	ENV_SETTING_SR_HOST = "TESTER_SR_SERVER"
)

//Provider is the extension provider for the flow action
type TesterProvider struct {
}

func NewExtensionProvider() *TesterProvider {
	return &TesterProvider{}
}

func (fp *TesterProvider) GetFlowProvider() definition.Provider {
	return definition.NewRemoteFlowProvider()
}

func (fp *TesterProvider) GetFlowModel() *model.FlowModel {
	return simple.New()
}

func (fp *TesterProvider) GetStateRecorder() instance.StateRecorder {

	config := &util.ServiceConfig{Enabled: true}

	server := os.Getenv(ENV_SETTING_SR_HOST)

	if server != "" {
		parts := strings.Split(server, ":")

		host := parts[0]
		port := "9090"

		if len(parts) > 1 {
			port = parts[1]
		}

		settings := map[string]string{
			"host": host,
			"port": port,
		}
		config.Settings = settings
	} else {
		config.Enabled = false
	}

	return instance.NewRemoteStateRecorder(config)
}

func (fp *TesterProvider) GetMapperFactory() flowdef.MapperFactory {
	return nil
}

func (fp *TesterProvider) GetLinkExprManagerFactory() flowdef.LinkExprManagerFactory {
	return nil
}

func (fp *TesterProvider) GetFlowTester() *RestEngineTester {

	config := &util.ServiceConfig{Enabled: true}

	settings := map[string]string{
		"port": os.Getenv(ENV_SETTING_PORT),
	}
	config.Settings = settings
	return NewRestEngineTester(config)
}
