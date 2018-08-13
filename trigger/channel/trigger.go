package channel

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-channel")

// ChannelTrigger CHANNEL trigger struct
type ChannelTrigger struct {
	metadata *trigger.Metadata
	//runner   action.Runner
	config *trigger.Config
	//handlers []*handler.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &ChannelFactory{metadata: md}
}

// ChannelFactory CHANNEL Trigger factory
type ChannelFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *ChannelFactory) New(config *trigger.Config) trigger.Trigger {
	return &ChannelTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *ChannelTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *ChannelTrigger) Initialize(ctx trigger.InitContext) error {

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		// setup handlers
		channel := strings.ToLower(handler.GetStringSetting("channel"))
		log.Debugf("Registering handler for channel [%s]", channel)
	}

	return nil
}

func (t *ChannelTrigger) Start() error {
	//ignore
	return nil
}

// Stop implements util.Managed.Stop
func (t *ChannelTrigger) Stop() error {
	//ignore
	return nil
}
