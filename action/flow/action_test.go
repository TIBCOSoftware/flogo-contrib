package flow

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
)

//TestInitNoFlavorError
func TestInitNoFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "No flow found in action data") {
				t.Fail()
			}
		}
	}()

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: []byte(`{}`)}
	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitUnCompressedFlowFlavorError
func TestInitUnCompressedFlowFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading uncompressed flow") {
				t.Fail()
			}
		}
	}()

	mockFlowData := []byte(`{"flow":{}}`)

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitCompressedFlowFlavorError
func TestInitCompressedFlowFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading compressed flow") {
				t.Fail()
			}
		}
	}()
	mockFlowData := []byte(`{"flowCompressed":""}`)

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitURIFlowFlavorError
func TestInitURIFlowFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading flow URI") {
				t.Fail()
			}
		}
	}()
	mockFlowData := []byte(`{"flowURI":""}`)

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

var testFlowActionCfg = `{
  "id": "flow",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
  "metadata": {
    "input": [],
    "output": []
  },
  "data":{
  "flow": {
    "model": "tibco-simple",
    "type": 1,
    "attributes": [],
    "rootTask": {
      "id": 1,
      "type": 1,
      "activityType": "",
      "ref": "",
      "name": "root",
      "tasks": [
        {
          "id": "log_2",
          "name": "Log Message",
          "description": "Simple Log Activity",
          "type": 1,
          "activityType": "tibco-log",
          "activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
          "attributes": [
            {
              "name": "message",
              "value": "First log",
              "required": false,
              "type": "string"
            },
            {
              "name": "flowInfo",
              "value": "false",
              "required": false,
              "type": "boolean"
            },
            {
              "name": "addToFlow",
              "value": "true",
              "required": false,
              "type": "boolean"
            }
          ]
        },
        {
          "id": "log_3",
          "name": "Log Message (2)",
          "description": "Simple Log Activity",
          "type": 1,
          "activityType": "tibco-log",
          "activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
          "attributes": [
            {
              "name": "message",
              "value": "Second log",
              "required": false,
              "type": "string"
            },
            {
              "name": "flowInfo",
              "value": "false",
              "required": false,
              "type": "boolean"
            },
            {
              "name": "addToFlow",
              "value": "true",
              "required": false,
              "type": "boolean"
            }
          ]
        }
      ],
      "links": [
        {
          "id": 1,
          "from": "log_2",
          "to": "log_3",
          "type": 0
        }
      ]
    }
  }
  }
}
`
var testRestartInitialState = `{
  "initialState": {
    "id": "90c3f713bf2b87e4e9a584892039a76b",
    "state": 0,
    "status": 100,
    "attrs": [],
    "flowUri": "flow",
    "workQueue": [
      {
        "id": 2,
        "execType": 10,
        "taskID": "log_2",
        "code": 0
      }
    ],
    "rootTaskEnv": {
      "id": 1,
      "taskId": "1",
      "taskDatas": [
        {
          "state": 20,
          "done": false,
          "attrs": [],
          "taskId": "log_2"
        }
      ],
      "linkDatas": []
    },
    "actionUri": "flow"
  },
  "interceptor": {
    "tasks": [
      {
        "id": "log_2",
        "inputs": [
          {
            "name": "message",
            "type": "string",
            "value": "test rerun 1"
          },
          {
            "name": "flowInfo",
            "type": "boolean",
            "value": "false"
          },
          {
            "name": "addToFlow",
            "type": "boolean",
            "value": "true"
          }
        ]
      }
    ]
  }
}
`

func TestFlowAction_Run_Restart(t *testing.T) {

	cfg := &action.Config{}
	err := json.Unmarshal([]byte(testFlowActionCfg), cfg)

	if err != nil {
		t.Error(err)
		return
	}

	ff := FlowFactory{}
	flowAction := ff.New(cfg)

	req := &RestartRequest{}
	err = json.Unmarshal([]byte(testRestartInitialState), req)

	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	if req.Data != nil {

		attrs := make([]*data.Attribute, len(req.Data))

		for k, v := range req.Data {
			attrs = append(attrs, data.NewAttribute(k, data.ANY, v))
		}

		ctx = trigger.NewContext(context.Background(), attrs)
	}

	execOptions := &instance.ExecOptions{Interceptor: req.Interceptor, Patch: req.Patch}
	ro := &instance.RunOptions{Op: instance.OpRestart, ReturnID: true, FlowURI: req.InitialState.FlowURI, InitialState: req.InitialState, ExecOptions: execOptions}

	r := runner.NewDirect()
	r.Run(ctx, flowAction, req.InitialState.FlowURI, ro)
}

type RestartRequest struct {
	InitialState *instance.Instance     `json:"initialState"`
	Data         map[string]interface{} `json:"data"`
	Interceptor  *support.Interceptor   `json:"interceptor"`
	Patch        *support.Patch         `json:"patch"`
}
