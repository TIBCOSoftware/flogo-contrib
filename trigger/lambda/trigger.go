package lambda

import (
	"context"
	"encoding/json"
	"flag"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-lambda")

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

	log.Debug("Starting AWS Lambda Trigger")

	// Parse the flags
	flag.Parse()

	// Looking up the arguments
	evtArg := flag.Lookup("evt")
	var evt interface{}
	// Unmarshall evt
	if err := json.Unmarshal([]byte(evtArg.Value.String()), &evt); err != nil {
		return err
	}

	log.Debugf("Received evt: '%+v'\n", evt)

	ctxArg := flag.Lookup("ctx")
	var ctx *runtime.Context
	// Unmarshall ctx
	if err := json.Unmarshal([]byte(ctxArg.Value.String()), &ctx); err != nil {
		return err
	}

	log.Debugf("Received ctx: '%+v'\n", ctx)

	actionId := t.config.Handlers[0].ActionId
	log.Debugf("Calling actionid: '%s'\n", actionId)

	action := action.Get(actionId)

	data := map[string]interface{}{
		"logStreamName":   ctx.LogStreamName,
		"logGroupName":    ctx.LogGroupName,
		"awsRequestId":    ctx.AWSRequestID,
		"memoryLimitInMB": ctx.MemoryLimitInMB,
	}

	startAttrs, err := t.metadata.OutputsToAttrs(data, false)
	if err != nil {
		log.Errorf("After run error' %s'\n", err)
		return err
	}

	context := trigger.NewContext(context.Background(), startAttrs)
	code, result, err := t.runner.Run(context, action, actionId, nil)

	log.Debugf("After run code: '%d', result: '%+v'\n", code, result)

	if err != nil {
		log.Debugf("After run error' %s'\n", err)
		return err
	}

	return err
}

// Stop implements util.Managed.Stop
func (t *LambdaTrigger) Stop() error {
	return nil
}
