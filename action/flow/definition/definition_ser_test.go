package definition

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const defJSON = `
{
    "type": 1,
    "name": "Demo Flow",
    "model": "simple",
    "attributes": [
      { "name": "petInfo", "type": "string", "value": "" }
    ],
    "rootTask": {
      "id": 1,
      "type": 1,
      "activityType": "",
      "name": "root",
      "tasks": [
        {
          "id": 2,
          "type": 1,
          "activityType": "log",
          "name": "Log Start",
          "attributes": [
            { "type": "string", "name": "message", "value": "Find Pet Flow Started!"},
            { "type": "boolean", "name": "flowInfo", "value": "true"}
          ]
        },
        {
          "id": 3,
          "type": 1,
          "activityType": "rest",
          "name": "Pet Query",
          "attributes": [
            { "type": "string", "name": "uri", "value": "http://petstore.swagger.io/v2/pet/{petId}" },
            { "type": "string", "name": "method", "value": "GET" },
            { "type": "string", "name": "petId", "value": "" },
            { "type": "string", "name": "result", "value": "" }
          ],
          "inputMappings": [
            { "type": 1, "value": "petId", "mapTo": "petId" }
          ],
          "ouputMappings": [
            { "type": 1, "value": "result", "mapTo": "petInfo" }
          ]
        },
        {
          "id": 4,
          "type": 1,
          "activityType": "log",
          "name": "Log Results",
          "attributes": [
            { "type": "string", "name": "message", "value": "REST results" },
            { "type": "boolean", "name": "flowInfo", "value": "true" }
          ],
          "inputMappings": [
            { "type": 1, "value": "petInfo", "result": "message" }
          ]
        }
      ],
      "links": [
        { "id": 1, "type": 1,  "name": "", "to": 3,  "from": 2 },
        { "id": 2, "type": 1, "name": "", "to": 4, "from": 3 }
      ]
    }
  }
`

func TestRestartWithFlowData(t *testing.T) {

	defRep := &DefinitionRep{}

	json.Unmarshal([]byte(defJSON), defRep)

	def, _ := NewDefinition(defRep)

	fmt.Printf("Definition: %v", def)
}

type MyDummyJson struct {
	Value interface{}
}

func TestConvertInterfaceToString(t *testing.T) {
	intJson := `{"Value": 1}`
	dummyJ := &MyDummyJson{}

	err := json.Unmarshal([]byte(intJson), &dummyJ)

	assert.Nil(t, err)

	result := convertInterfaceToString(dummyJ.Value)

	assert.Equal(t, "1", result)

	strJson := `{"Value": "stringId"}`
	dummyJ = &MyDummyJson{}

	err = json.Unmarshal([]byte(strJson), &dummyJ)

	assert.Nil(t, err)

	result = convertInterfaceToString(dummyJ.Value)

	assert.Equal(t, "stringId", result)

}
