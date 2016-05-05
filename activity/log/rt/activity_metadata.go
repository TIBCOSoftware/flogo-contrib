package log

var jsonMetadata = `{
  "name": "log",
  "version": "0.0.1",
  "description": "log activity",
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
