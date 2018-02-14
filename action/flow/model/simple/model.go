package simple

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("flowmodel-simple")

const (
	MODEL_NAME = "flogo-simple"
)

func init() {
	model.RegisterDefault(New())
}

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&behaviors.Flow{})
	m.RegisterDefaultTaskBehavior(&behaviors.Task{})
	m.RegisterTaskBehavior(2, "iterator", &behaviors.IteratorTask{})
	return m
}
