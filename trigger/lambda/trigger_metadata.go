package lambda

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = `{
  "name": "tibco-lambda",
  "type": "flogo:trigger",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/lambda",
  "version": "0.0.1",
  "title": "Start Flow as a function in Lambda",
  "description": "Simple Lambda Trigger",
  "homepage": "https://github.com/TIBCOSoftware/flogo-contrib/tree/master/trigger/lambda",
  "settings": [
  ],
  "outputs": [
    {
      "name": "logStreamName",
      "type": "string"
    },
    {
      "name": "logGroupName",
      "type": "string"
    },
    {
      "name": "awsRequestId",
      "type": "string"
    },
    {
      "name": "memoryLimitInMB",
      "type": "string"
    }
  ]
}
`

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
