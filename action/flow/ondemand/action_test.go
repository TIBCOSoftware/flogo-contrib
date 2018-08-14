package ondemand

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/stretchr/testify/assert"

	_ "github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
)

const testEventJson =`
{
  "payload": {
    "in1":"in1_value",
    "in2":"in2_value"
  },
  "flogo" : {
      "inputMappings": [
        { "type": "assign", "value": "$.payload.in1", "mapTo": "customerId" },
        { "type": "assign", "value": "$.payload.in2", "mapTo": "orderId" }
      ],
      "flow": {
        "metadata" : {
          "input":[
            { "name":"customerId", "type":"string" },
            { "name":"orderId", "type":"string" }
          ],
          "output":[
            { "name":"value", "type":"string" }
          ]
        },
        "tasks": [
          {
            "id": "LogResult",
            "name": "Log Results",
            "activity" : {
              "ref":"test-log",
              "input" : {
                "message" : "REST results"
              },
              "mappings" : {
                "input": [
                  { "type": "assign", "value": "$flow.customerId", "mapTo": "message" }
                ]
              }
            }
          }
        ]
      }
  }
}`

type Event struct {
	Payload interface{} `json:"payload"`
	Flogo   json.RawMessage `json:"flogo"`
}

//TestInitNoFlavorError
func TestFlowAction_Run(t *testing.T) {

	var evt Event

	// Unmarshall evt
	if err := json.Unmarshal([]byte(testEventJson), &evt); err != nil {
		assert.Nil(t,err)
		return
	}

	cfg := &action.Config{}

	ff := ActionFactory{}
	ff.Init()
	fa, err := ff.New(cfg)
	assert.Nil(t, err)

	flowAction, ok := fa.(action.AsyncAction)
	assert.True(t, ok)

	inputs := make(map[string]*data.Attribute, 2)
	attr, _ := data.NewAttribute("payload", data.TypeObject, evt.Payload)
	inputs[attr.Name()] = attr
	attr, _ = data.NewAttribute("flowPackage", data.TypeAny, evt.Flogo)
	inputs[attr.Name()] = attr

	r := runner.NewDirect()
	_, err = r.Execute(context.Background(), flowAction, inputs)
	assert.Nil(t, err)
}
