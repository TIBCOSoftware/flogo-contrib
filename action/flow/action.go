package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/extension"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/provider"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/tester"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/action/flow"
)

// ActionOptions are the options for the FlowAction
type ActionOptions struct {
	MaxStepCount int
	Record       bool
}

type FlowAction struct {
	idGenerator   *util.Generator
	actionOptions *ActionOptions
}

// Provides the different extension points to the Flow Action
type ExtensionProvider interface {
	GetFlowProvider() provider.Provider
	GetFlowModel() *model.FlowModel
	GetStateRecorder() instance.StateRecorder
	GetMapperFactory() definition.MapperFactory
	GetLinkExprManagerFactory() definition.LinkExprManagerFactory
	GetFlowTester() *tester.RestEngineTester
}

var actionMu sync.Mutex
var ep ExtensionProvider
var flowAction *FlowAction

func init() {
	action.RegisterFactory(FLOW_REF, &FlowFactory{})
}

func SetExtensionProvider(provider ExtensionProvider) {
	actionMu.Lock()
	defer actionMu.Unlock()

	ep = provider
}

type FlowFactory struct{}

func (ff *FlowFactory) New(config *action.Config) action.Action {

	actionMu.Lock()
	defer actionMu.Unlock()

	if flowAction == nil {
		options := &ActionOptions{Record: false}

		if ep == nil {
			testerEnabled := os.Getenv(tester.ENV_ENABLED)
			if strings.ToLower(testerEnabled) == "true" {
				ep = tester.NewExtensionProvider()

				sm := util.GetDefaultServiceManager()
				sm.RegisterService(ep.GetFlowTester())
				options.Record = true
			} else {
				ep = extension.New()
			}
		}

		definition.SetMapperFactory(ep.GetMapperFactory())
		definition.SetLinkExprManagerFactory(ep.GetLinkExprManagerFactory())

		if options.MaxStepCount < 1 {
			options.MaxStepCount = int(^uint16(0))
		}

		flowAction = &FlowAction{}

		flowAction.actionOptions = options
		flowAction.idGenerator, _ = util.NewGenerator()
	}

	//temporary hack to support dynamic process running by tester
	if config.Data == nil {
		return flowAction
	}

	var flavor Flavor
	err := json.Unmarshal(config.Data, &flavor)
	if err != nil {
		errorMsg := fmt.Sprintf("Error while loading flow '%s' error '%s'", config.Id, err.Error())
		logger.Errorf(errorMsg)
		panic(errorMsg)
	}

	if len(flavor.Flow) > 0 {
		// It is an uncompressed and embedded flow
		err := ep.GetFlowProvider().AddUncompressedFlow(config.Id, flavor.Flow)
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading uncompressed flow '%s' error '%s'", config.Id, err.Error())
			logger.Errorf(errorMsg)
			panic(errorMsg)
		}
		return flowAction
	}

	if len(flavor.FlowCompressed) > 0 {
		// It is a compressed and embedded flow
		err := ep.GetFlowProvider().AddCompressedFlow(config.Id, string(flavor.FlowCompressed[:]))
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading compressed flow '%s' error '%s'", config.Id, err.Error())
			logger.Errorf(errorMsg)
			panic(errorMsg)
		}
		return flowAction
	}

	if len(flavor.FlowURI) > 0 {
		// It is a URI flow
		err := ep.GetFlowProvider().AddFlowURI(config.Id, string(flavor.FlowURI[:]))
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading flow URI '%s' error '%s'", config.Id, err.Error())
			logger.Errorf(errorMsg)
			panic(errorMsg)
		}
		return flowAction
	}

	errorMsg := fmt.Sprintf("No flow found in action data for id '%s'", config.Id)
	logger.Errorf(errorMsg)
	panic(errorMsg)

	return flowAction
}

// Run implements action.Action.Run
func (fa *FlowAction) Run(context context.Context, uri string, options interface{}, handler action.ResultHandler) error {

	logger.Infof("In Flow Run uri: '%s'", uri)
	//todo: catch panic
	//todo: consider switch to URI to dictate flow operation (ex. flow://blah/resume)

	op := instance.OpStart
	retID := false

	ro, ok := options.(*instance.RunOptions)

	if ok {
		op = ro.Op
		retID = ro.ReturnID
	}

	var inst *instance.Instance

	switch op {
	case instance.OpStart:
		flow, err := ep.GetFlowProvider().GetFlow(uri)
		if err != nil {
			return err
		}

		instanceID := fa.idGenerator.NextAsString()
		logger.Debug("Creating Instance: ", instanceID)

		inst = instance.New(instanceID, uri, flow, ep.GetFlowModel())
	case instance.OpResume:
		if ok {
			inst = ro.InitialState
			logger.Debug("Resuming Instance: ", inst.ID())
		} else {
			return errors.New("Unable to resume instance, resume options not provided")
		}
	case instance.OpRestart:
		if ok {
			inst = ro.InitialState
			instanceID := fa.idGenerator.NextAsString()
			inst.Restart(instanceID, ep.GetFlowProvider())

			logger.Debug("Restarting Instance: ", instanceID)
		} else {
			return errors.New("Unable to restart instance, restart options not provided")
		}
	}

	if ok && ro.ExecOptions != nil {
		logger.Debugf("Applying Exec Options to instance: %s\n", inst.ID())
		instance.ApplyExecOptions(inst, ro.ExecOptions)
	}

	triggerAttrs, ok := trigger.FromContext(context)

	if ok {
		if len(triggerAttrs) > 0 {
			logger.Debug("Run Attributes:")
			for _, attr := range triggerAttrs {
				logger.Debugf(" Attr:%s, Type:%s, Value:%v", attr.Name, attr.Type.String(), attr.Value)
			}
		}
	}

	if op == instance.OpStart {
		inst.Start(triggerAttrs)
	} else {
		inst.UpdateAttrs(triggerAttrs)
	}

	logger.Debugf("Executing instance: %s\n", inst.ID())

	stepCount := 0
	hasWork := true

	inst.SetReplyHandler(&SimpleReplyHandler{resultHandler: handler})

	go func() {

		defer handler.Done()

		if !inst.Flow.ExplicitReply() {
			handler.HandleResult(200, &instance.IDResponse{ID: inst.ID()}, nil)
		}

		for hasWork && inst.Status() < instance.StatusCompleted && stepCount < fa.actionOptions.MaxStepCount {
			stepCount++
			logger.Debugf("Step: %d\n", stepCount)
			hasWork = inst.DoStep()

			if fa.actionOptions.Record {
				ep.GetStateRecorder().RecordSnapshot(inst)
				ep.GetStateRecorder().RecordStep(inst)
			}
		}

		if retID {
			handler.HandleResult(200, &instance.IDResponse{ID: inst.ID()}, nil)
		}

		logger.Debugf("Done Executing A.instance [%s] - Status: %d\n", inst.ID(), inst.Status())

		if inst.Status() == instance.StatusCompleted {
			logger.Infof("Flow [%s] Completed", inst.ID())
		}
	}()

	return nil
}

// SimpleReplyHandler is a simple ReplyHandler that is pass-thru to the action ResultHandler
type SimpleReplyHandler struct {
	resultHandler action.ResultHandler
}

// Reply implements ReplyHandler.Reply
func (rh *SimpleReplyHandler) Reply(replyCode int, replyData interface{}, err error) {

	rh.resultHandler.HandleResult(replyCode, replyData, err)
}
