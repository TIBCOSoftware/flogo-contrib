package channel

import (
	"io/ioutil"
)

var jsonTestMetadata = getTestJsonMetadata()

func getTestJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "id": "flogo-channel",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/channel",
  "handlers": [
    {
      "settings": {
        "channel": "test"
      }
      "action" : {
		"id": "my_test_flow"
      }
    }
  ]
}
`


