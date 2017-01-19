package flow

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/op/go-logging"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/incubator/flow"
)

var log = logging.MustGetLogger("flow")

type FlowAction struct {
	Id string
}

type FlowFactory struct{}

func init() {
	action.GetRegistry().RegisterFactory(FLOW_REF, &FlowFactory{})
}

func (fa *FlowFactory) New(id string) action.Action2 {
	return &FlowAction{Id: id}
}

func (fa *FlowAction) Init(config types.ActionConfig) {
	log.Infof("In Flow Init '%s'", fa.Id)
}

// Run implements action.Action.Run
func (fa *FlowAction) Run(context context.Context, uri string, options interface{}, handler action.ResultHandler) error {
	log.Infof("In Flow Run '%s'", fa.Id)
	return nil
}
