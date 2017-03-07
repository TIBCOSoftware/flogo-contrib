package instance

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// ExecOptions are optional Patch & Interceptor to be used during instance execution
type ExecOptions struct {
	Patch       *support.Patch
	Interceptor *support.Interceptor
}

// IDGenerator generates IDs for flow instances
type IDGenerator interface {

	//NewFlowInstanceID generate a new instance ID
	NewFlowInstanceID() string
}

// ApplyExecOptions applies any execution options to the flow instance
func ApplyExecOptions(instance *Instance, execOptions *ExecOptions) {

	if execOptions != nil {

		if execOptions.Patch != nil {
			logger.Infof("Instance [%s] has patch", instance.ID())
			instance.Patch = execOptions.Patch
			instance.Patch.Init()
		}

		if execOptions.Interceptor != nil {
			logger.Infof("Instance [%s] has interceptor", instance.ID)
			instance.Interceptor = execOptions.Interceptor
			instance.Interceptor.Init()
		}
	}
}
