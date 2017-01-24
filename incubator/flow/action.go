package flow

import (
	"context"
	"errors"
	"fmt"
	//"encoding/json"

	//flow_types "github.com/TIBCOSoftware/flogo-contrib/incubator/flow/types"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/flow/service/flowprovider"
	"github.com/TIBCOSoftware/flogo-lib/flow/service/staterecorder"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/incubator/flow"
)

var log = logging.MustGetLogger("flow")

type FlowAction struct {
	stateRecorder flowinst.StateRecorder
	flowProvider  flowdef.Provider
	idGenerator   *util.Generator
	actionOptions *flowinst.ActionOptions
}

type FlowFactory struct{}

func init() {
	action.RegisterFactory(FLOW_REF, &FlowFactory{})
}

func (fa *FlowFactory) New(id string) action.Action2 {
	return &FlowAction{}
}

func (fa *FlowAction) Init(config types.ActionConfig, serviceManager *util.ServiceManager) {
	log.Infof("In Flow Init")

	embeddedJSONFlows := make(map[string]string)
	embeddedJSONFlows["embedded://"+config.Id] = string(config.Data[:])
	embeddedFlowMgr := support.NewEmbeddedFlowManager(false, embeddedJSONFlows)

	// TODO extract this to contribution
	fpConfig := &util.ServiceConfig{Name: "flowProvider", Enabled: true}
	flowProvider := flowprovider.NewRemoteFlowProvider(fpConfig, embeddedFlowMgr)
	serviceManager.RegisterService(flowProvider)
	fa.flowProvider = flowProvider

	// TODO extract this to contribution
	srSettings := make(map[string]string, 2)
	srSettings["host"] = ""
	srSettings["port"] = ""
	srConfig := &util.ServiceConfig{Name: "stateRecorder", Enabled: false, Settings: srSettings}
	sr := staterecorder.NewRemoteStateRecorder(srConfig)
	serviceManager.RegisterService(sr)
	fa.stateRecorder = sr

	// TODO extract this to contribution
	//	etConfig := &util.ServiceConfig{Name: "engineTester", Enabled: false}
	//	engineTester := tester.NewRestEngineTester(etConfig)
	//	serviceManager.RegisterService(engineTester)

	options := &flowinst.ActionOptions{Record: sr.Enabled()}

	fa.idGenerator, _ = util.NewGenerator()

	if options.MaxStepCount < 1 {
		options.MaxStepCount = int(^uint16(0))
	}

	fa.actionOptions = options
}

// RunOptions the options when running a FlowAction
type RunOptions struct {
	Op           int
	ReturnID     bool
	InitialState *flowinst.Instance
	ExecOptions  *flowinst.ExecOptions
}

// Run implements action.Action.Run
func (fa *FlowAction) Run(context context.Context, uri string, options interface{}, handler action.ResultHandler) error {

	log.Infof("In Flow Run uri: '%s'", uri)
	//todo: catch panic
	//todo: consider switch to URI to dictate flow operation (ex. flow://blah/resume)

	op := flowinst.AoStart
	retID := false

	ro, ok := options.(*flowinst.RunOptions)

	if ok {
		op = ro.Op
		retID = ro.ReturnID
	}

	var instance *flowinst.Instance

	switch op {
	case flowinst.AoStart:
		flow := fa.flowProvider.GetFlow("embedded://" + uri)

		if flow == nil {
			err := fmt.Errorf("Flow [%s] not found", uri)
			return err
		}

		instanceID := fa.idGenerator.NextAsString()
		log.Debug("Creating Instance: ", instanceID)

		instance = flowinst.NewFlowInstance(instanceID, uri, flow)
	case flowinst.AoResume:
		if ok {
			instance = ro.InitialState
			log.Debug("Resuming Instance: ", instance.ID())
		} else {
			return errors.New("Unable to resume instance, resume options not provided")
		}
	case flowinst.AoRestart:
		if ok {
			instance = ro.InitialState
			instanceID := fa.idGenerator.NextAsString()
			instance.Restart(instanceID, fa.flowProvider)

			log.Debug("Restarting Instance: ", instanceID)
		} else {
			return errors.New("Unable to restart instance, restart options not provided")
		}
	}

	if ok && ro.ExecOptions != nil {
		log.Debugf("Applying Exec Options to instance: %s\n", instance.ID())
		flowinst.ApplyExecOptions(instance, ro.ExecOptions)
	}

	triggerAttrs, ok := trigger.FromContext(context)

	if log.IsEnabledFor(logging.DEBUG) && ok {
		if len(triggerAttrs) > 0 {
			log.Debug("Run Attributes:")
			for _, attr := range triggerAttrs {
				log.Debugf(" Attr:%s, Type:%s, Value:%v", attr.Name, attr.Type.String(), attr.Value)
			}
		}
	}

	if op == flowinst.AoStart {
		instance.Start(triggerAttrs)
	} else {
		instance.UpdateAttrs(triggerAttrs)
	}

	log.Debugf("Executing instance: %s\n", instance.ID())

	stepCount := 0
	hasWork := true

	instance.SetReplyHandler(&SimpleReplyHandler{resultHandler: handler})

	go func() {

		defer handler.Done()

		if !instance.Flow.ExplicitReply() {
			handler.HandleResult(200, &IDResponse{ID: instance.ID()}, nil)
		}

		for hasWork && instance.Status() < flowinst.StatusCompleted && stepCount < fa.actionOptions.MaxStepCount {
			stepCount++
			log.Debugf("Step: %d\n", stepCount)
			hasWork = instance.DoStep()

			if fa.actionOptions.Record {
				fa.stateRecorder.RecordSnapshot(instance)
				fa.stateRecorder.RecordStep(instance)
			}
		}

		if retID {
			handler.HandleResult(200, &IDResponse{ID: instance.ID()}, nil)
		}

		log.Debugf("Done Executing A.instance [%s] - Status: %d\n", instance.ID(), instance.Status())

		if instance.Status() == flowinst.StatusCompleted {
			log.Infof("Flow [%s] Completed", instance.ID())
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

// IDResponse is a response object consists of an ID
type IDResponse struct {
	ID string `json:"id"`
}
