package flow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowinst"
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
	flowProvider  Provider
	idGenerator   *util.Generator
	actionOptions *flowinst.ActionOptions
}

var flowAction *FlowAction

type FlowFactory struct{}

func init() {
	action.RegisterFactory(FLOW_REF, &FlowFactory{})
	flowAction = NewFlowAction()
}

func (fa *FlowFactory) New(id string) action.Action2 {
	return flowAction
}

// NewFlowAction creates a new FlowAction
func NewFlowAction() *FlowAction {

	fa := &FlowAction{}
	fa.flowProvider = NewRemoteFlowProvider()

	// TODO add state recorder
	//	srSettings := make(map[string]string, 2)
	//	srSettings["host"] = ""
	//	srSettings["port"] = ""
	//	srConfig := &util.ServiceConfig{Name: "stateRecorder", Enabled: false, Settings: srSettings}
	//	sr := staterecorder.NewRemoteStateRecorder(srConfig)
	//	serviceManager.RegisterService(sr)
	//	fa.stateRecorder = sr

	// TODO add engine tester
	//	etConfig := &util.ServiceConfig{Name: "engineTester", Enabled: false}
	//	engineTester := tester.NewRestEngineTester(etConfig)
	//	serviceManager.RegisterService(engineTester)

	options := &flowinst.ActionOptions{Record: false}

	fa.idGenerator, _ = util.NewGenerator()

	if options.MaxStepCount < 1 {
		options.MaxStepCount = int(^uint16(0))
	}

	fa.actionOptions = options
	return fa
}

func (fa *FlowAction) Init(config types.ActionConfig) {
	log.Debugf("Initializing flow '%s'", config.Id)

	var flavor Flavor
	err := json.Unmarshal(config.Data, &flavor)
	if err != nil {
		errorMsg := fmt.Sprintf("Error while loading flow '%s' error '%s'", config.Id, err.Error())
		log.Errorf(errorMsg)
		panic(errorMsg)
	}

	if len(flavor.Flow) > 0 {
		// It is an uncompressed and embedded flow
		err := fa.flowProvider.AddUncompressedFlow(config.Id, flavor.Flow)
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading uncompressed flow '%s' error '%s'", config.Id, err.Error())
			log.Errorf(errorMsg)
			panic(errorMsg)
		}
		return
	}

	if len(flavor.FlowCompressed) > 0 {
		// It is a compressed and embedded flow
		err := fa.flowProvider.AddCompressedFlow(config.Id, string(flavor.FlowCompressed[:]))
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading compressed flow '%s' error '%s'", config.Id, err.Error())
			log.Errorf(errorMsg)
			panic(errorMsg)
		}
		return
	}

	if len(flavor.FlowURI) > 0 {
		// It is a URI flow
		err := fa.flowProvider.AddFlowURI(config.Id, string(flavor.FlowURI[:]))
		if err != nil {
			errorMsg := fmt.Sprintf("Error while loading flow URI '%s' error '%s'", config.Id, err.Error())
			log.Errorf(errorMsg)
			panic(errorMsg)
		}
		return
	}

	errorMsg := fmt.Sprintf("No flow found in action data for id '%s'", config.Id)
	log.Errorf(errorMsg)
	panic(errorMsg)

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
		flow, err := fa.flowProvider.GetFlow(uri)
		if err != nil {
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
