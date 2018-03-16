package cli

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
  "id": "flogo-cli",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/cli",
  "handlers": [
    {
      "actionId": "",
      "settings": {
        "command": "run",
        "default": "true"
      }
    },
    {
      "actionId": "version_flow",
      "settings": {
        "command": "version"
      }
    }
  ]
}
`

/*
//TODO fix this test
func TestInitOk(t *testing.T) {
	// New  factory
	f := &CliTriggerFactory{}
	tgr := f.New("flogo-cli")

	runner := &TestRunner{}

	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr.Init(config, runner)
}
*/

/*
//TODO fix this test
func TestHandlerOk(t *testing.T) {

	// New  factory
	f := &CliTriggerFactory{}
	tgr := f.New("flogo-cli")

	runner := &TestRunner{}

	config := trigger.Config{}
	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	uri := "http://127.0.0.1:8091/device/12345/reset"

	req, err := http.NewRequest("POST", uri, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		t.Fail()
	}
}
*/
