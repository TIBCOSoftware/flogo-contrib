package wits0

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-wits0")

// wits0TriggerFactory My wits0Trigger factory
type wits0TriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory create a new wits0Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &wits0TriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (trigger *wits0TriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &wits0Trigger{metadata: trigger.metadata, config: config}
}

// wits0Trigger is the wits0Trigger implementation
type wits0Trigger struct {
	metadata  *trigger.Metadata
	runner    action.Runner
	config    *trigger.Config
	handlers  []*trigger.Handler
	stopCheck chan bool
}

func (trigger *wits0Trigger) Initialize(ctx trigger.InitContext) error {
	log.Debug("Initialize")
	trigger.handlers = ctx.GetHandlers()
	return nil
}

func (trigger *wits0Trigger) Metadata() *trigger.Metadata {
	return trigger.metadata
}

// Start implements trigger.wits0Trigger.Start
func (trigger *wits0Trigger) Start() error {
	log.Debug("Start")
	trigger.stopCheck = make(chan bool)
	handlers := trigger.handlers
	for _, handler := range handlers {
		serial := &serialPort{}
		serial.Init(trigger, handler)
		serial.createSerialConnection()
	}
	return nil
}

// Stop implements trigger.wits0Trigger.Start
func (trigger *wits0Trigger) Stop() error {
	// stop the trigger
	close(trigger.stopCheck)
	return nil
}
