package flow

import (
	"context"
	"encoding/json"

	flow_types "github.com/TIBCOSoftware/flogo-contrib/incubator/flow/types"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/op/go-logging"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/incubator/flow"
)

var log = logging.MustGetLogger("flow")

type FlowAction struct {
}

type FlowFactory struct{}

func init() {
	action.RegisterFactory(FLOW_REF, &FlowFactory{})
}

func (fa *FlowFactory) New(id string) action.Action2 {
	return &FlowAction{}
}

func (fa *FlowAction) Init(config types.ActionConfig) {
	log.Infof("In Flow Init")
	// Parse to flow_types.Config
	var flowConfig flow_types.FlowConfig
	err := json.Unmarshal(config.Data, &flowConfig)
	if err != nil {
		panic(err.Error())
	}

}

// Run implements action.Action.Run
func (fa *FlowAction) Run(context context.Context, uri string, options interface{}, handler action.ResultHandler) error {
	log.Infof("In Flow Run")
	return nil
}
