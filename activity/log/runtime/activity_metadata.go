package log

var jsonMetadata = `{
  "name": "tibco-log",
  "version": "0.0.1",
  "description": "log Message",
  "inputs":[
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "flowInfo",
      "type": "boolean"
    },
    {
      "name": "addToFlow",
      "type": "boolean"
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    }
  ]
}`
