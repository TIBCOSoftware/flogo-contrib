package lambda

import (
	"context"
	"encoding/json"
	syslog "log"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"flag"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

// LambdaTrigger AWS Lambda trigger struct
type LambdaTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &LambdaFactory{metadata: md}
}

// LambdaFactory AWS Lambda Trigger factory
type LambdaFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *LambdaFactory) New(config *trigger.Config) trigger.Trigger {
	return &LambdaTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *LambdaTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *LambdaTrigger) Init(runner action.Runner) {
	t.runner = runner
}

func (t *LambdaTrigger) Start() error {

	syslog.Printf("Starting AWS Lambda Trigger\n")

	// Parse the flags
	flag.Parse()

	// Looking up the arguments
	evtArg := flag.Lookup("evt")
	var evt interface{}
	// Unmarshall evt
	if err := json.Unmarshal([]byte(evtArg.Value.String()), &evt); err != nil {
		return err
	}

	syslog.Printf("Received evt: '%+v'\n", evt)

	ctxArg := flag.Lookup("ctx")
	var ctx *runtime.Context
	// Unmarshall ctx
	if err := json.Unmarshal([]byte(ctxArg.Value.String()), &ctx); err != nil {
		return err
	}

	syslog.Printf("Received ctx: '%+v'\n", ctx)

	actionId := t.config.Handlers[0].ActionId
	syslog.Printf("Hi there inside trigger calling actionid: '%s'!!\n", actionId)

	action := action.Get(actionId)
	syslog.Printf("Found action' %+x'\n", action)

	context := trigger.NewContext(context.Background(), make([]*data.Attribute,0))
	code, data, err := t.runner.Run(context, action, actionId, nil)

	syslog.Printf("After run error code: '%d', data: '%+v'\n", code, data)

	if err != nil {
		syslog.Printf("After run error' %s'\n", err)
	}

	return err
}

// Stop implements util.Managed.Stop
func (t *LambdaTrigger) Stop() error {
	return nil
}