package coap

var jsonMetadata = `{
  "name": "tibco-coap",
  "version": "0.0.1",
  "title": "CoAP",
  "description": "Simple CoAP Activity",
  "inputs":[
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "method",
      "type": "string",
      "required": true
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "type",
      "type": "string"
    },
    {
      "name": "messageId",
      "type": "integer"
    },
    {
      "name": "options",
      "type": "params"
    },
    {
      "name": "payload",
      "type": "string"
    }

  ],
  "outputs": [
    {
      "name": "response",
      "type": "string"
    }
  ]
}
`
