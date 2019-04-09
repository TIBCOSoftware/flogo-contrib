package wits0

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
	"id": "wits0",
	"settings": {

	},
	"handlers": [{
		"action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:query"
            }
          },
		"settings": {
			"SerialPort": "/dev/ttyUSB0",
			"HeartBeatValue": "&&\n0111-9999\n!!",
			"PacketHeader": "&&",
			"PacketFooter": "!!",
			"LineSeparator":"\r\n",
			"HeartBeatInterval": 1
		}
	}]
}`

const testConfigRaw string = `{
	"id": "wits0",
	"settings": {

	},
	"handlers": [{
		"action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:query"
            }
          },
		"settings": {
			"SerialPort": "/dev/ttyUSB0",
			"HeartBeatValue": "&&\n0111-9999\n!!",
			"PacketHeader": "&&",
			"PacketFooter": "!!",
			"LineSeparator":"\r\n",			
			"OutputRaw": true
		}
	}]
}`

const testConfigBaseSerialPort string = `{
	"id": "wits0",
	"settings": {
	},
	"handlers": [{
		"action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:query"
            }
          },
		"settings": {
			"SerialPort": "/dev/dummy",
			"HeartBeatValue": "&&\n0111-9999\n!!",
			"PacketHeader": "&&",
			"PacketFooter": "!!",
			"LineSeparator":"\r\n",
			"HeartBeatInterval": 2
		}
	}]
}`

const testData string = `
&&
1984PASON/EDR
01085871.95
01105893.00
0112110.00
01130.00
011544.38
01170.00
01190.00
01200.00
01210.00
01220.00
01230.00
01240.00
0125-9999.00
0126718.42
012717.50
01281.06
01300.00
01374483.00
0139-8888.00
0140-9999.00
0141-9999.00
0142313557.96
01430.00
01444483.24
01450.00
01695896.69
017023.30
0171-4066.40
0172-9999.00
0173-9999.00
!!

&&
1984PASON/EDR
01085871.95
01105893.00
0112108.46
01130.00
011544.79
01170.00
01190.01
01200.00
01210.00
01220.00
01230.00
01240.00
0125-9999.00
0126718.03
012717.12
01281.06
01300.00
01374483.00
0139-8888.00
0140-9999.00
0141-9999.00
0142313557.96
01430.00
01444483.24
01450.00
01695896.69
017023.37
0171-4066.40
0172-9999.00
0173-9999.00
!!

&&
1984PASON/EDR
631070.53
631170.53
633970.53
634070.53
!!

`

type initContext struct {
	handlers []*trigger.Handler
}

func (ctx initContext) GetHandlers() []*trigger.Handler {
	return ctx.handlers
}

type TestRunner struct {
}

var Test action.Runner

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	//log.Infof("Ran Action (Run): %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	//log.Infof("Ran Action (RunAction): %v", act)
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	//log.Infof("Ran Action (Execute): %v", act)
	value := inputs["data"].Value().(string)
	log.Info(value)
	return nil, nil
}

type TestAction struct {
}

func (tr *TestAction) Metadata() *action.Metadata {
	//log.Infof("Metadata")
	return nil
}

func (tr *TestAction) IOMetadata() *data.IOMetadata {
	//log.Infof("IOMetadata")
	return nil
}

func TestParse(t *testing.T) {
	trg, config := createTrigger(t, testConfig)
	initializeTrigger(t, trg, config)
	serialPort := &serialPort{}
	trgWits0 := trg.(*wits0Trigger)
	serialPort.Init(trgWits0, trgWits0.handlers[0])
	replaceData := strings.ReplaceAll(testData, "\n", "\r\n")
	data := bytes.NewBufferString(replaceData)
	outputBuffer := serialPort.parseBuffer(data)
	log.Debug(outputBuffer)
}

func TestParseRaw(t *testing.T) {
	trg, config := createTrigger(t, testConfigRaw)
	initializeTrigger(t, trg, config)
	serialPort := &serialPort{}
	trgWits0 := trg.(*wits0Trigger)
	serialPort.Init(trgWits0, trgWits0.handlers[0])
	replaceData := strings.ReplaceAll(testData, "\n", "\r\n")
	data := bytes.NewBufferString(replaceData)
	outputBuffer := serialPort.parseBuffer(data)
	log.Debug(outputBuffer)
}
func TestConnect(t *testing.T) {
	trg, config := createTrigger(t, testConfig)
	initializeTrigger(t, trg, config)
	runTrigger(5, trg)
}

func TestConnectRaw(t *testing.T) {
	trg, config := createTrigger(t, testConfigRaw)
	initializeTrigger(t, trg, config)
	runTrigger(5, trg)
}

func TestBadSerialPort(t *testing.T) {
	trg, config := createTrigger(t, testConfigBaseSerialPort)
	initializeTrigger(t, trg, config)
	runTrigger(5, trg)
}

func runTrigger(timeout int, trg trigger.Trigger) {
	go func() {
		time.Sleep(time.Second * time.Duration(timeout))
		trg.Stop()
	}()

	trg.Start()
}

func createTrigger(t *testing.T, conf string) (trigger.Trigger, trigger.Config) {
	log.SetLogLevel(logger.DebugLevel)
	md := trigger.NewMetadata(getJSONMetadata())
	f := NewFactory(md)
	config := trigger.Config{}
	if f == nil {
		t.Fail()
		return nil, config
	}

	jsonErr := json.Unmarshal([]byte(conf), &config)
	if jsonErr != nil {
		log.Error(jsonErr)
		t.Fail()
		return nil, config
	}
	trg := f.New(&config)
	return trg, config
}

func initializeTrigger(t *testing.T, trg trigger.Trigger, config trigger.Config) trigger.Initializable {
	if trg == nil {
		t.Fail()
		return nil
	}
	newTrg, _ := trg.(trigger.Initializable)

	initCtx := &initContext{handlers: make([]*trigger.Handler, 0, len(config.Handlers))}
	runner := &TestRunner{}
	action := &TestAction{}
	//create handlers for that trigger and init
	for _, hConfig := range config.Handlers {
		log.Infof("hConfig: %v", hConfig)
		log.Infof("trg.Metadata().Output: %v", trg.Metadata().Output)
		log.Infof("trg.Metadata().Reply: %v", trg.Metadata().Reply)
		handler := trigger.NewHandler(hConfig, action, trg.Metadata().Output, trg.Metadata().Reply, runner)
		initCtx.handlers = append(initCtx.handlers, handler)
	}

	newTrg.Initialize(initCtx)
	return newTrg
}
